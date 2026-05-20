package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"terraforming-mars-backend/internal/auth"
	"terraforming-mars-backend/internal/models"
	"terraforming-mars-backend/internal/rating"
	gamedata "terraforming-mars-backend/shared"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreatePlayer creates a new player with role-based validation
func (r *Repository) CreatePlayer(name string, password *string, role models.PlayerRole, actor models.Player) (*models.Player, error) {
	if err := auth.CanCreatePlayers(actor, role); err != nil {
		return nil, err
	}

	return r.createPlayer(name, password, role, &actor.ID)
}

// CreateSystemAdmin creates the initial system admin (no actor required)
func (r *Repository) CreateSystemAdmin(name string, password *string) (*models.Player, error) {
	return r.createPlayer(name, password, models.RoleAdmin, nil)
}

func (r *Repository) createPlayer(name string, password *string, role models.PlayerRole, createdBy *int) (*models.Player, error) {
	if !role.IsValid() {
		return nil, fmt.Errorf("Invalid role '%s'", role)
	}

	if auth.RequiresPassword(role) {
		if password == nil {
			return nil, fmt.Errorf("role '%s' requires a password", role)
		}
	} else {
		if password != nil {
			return nil, fmt.Errorf("role '%s' cannot have a password", role)
		}
	}

	var passwordHash *string = nil
	if password != nil {
		hash, err := auth.HashPassword(*password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		passwordHash = &hash
	}

	result, err := r.db.Exec("INSERT INTO player (name, password_hash, role, created_by) VALUES (?, ?, ?, ?)", name, passwordHash, role, createdBy)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetPlayerByID(int(id))
}

// GetPlayerByID retrieves a player by ID
func (r *Repository) GetPlayerByID(id int) (*models.Player, error) {
	var player models.Player
	err := r.db.QueryRow("SELECT id, name, password_hash, role, created_by, created_at, updated_at FROM player WHERE id = ?", id).Scan(
		&player.ID, &player.Name, &player.PasswordHash, &player.Role, &player.CreatedBy, &player.CreatedAt, &player.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// GetPlayerExtendedInfo retrieves player with additional statistics
func (r *Repository) GetPlayerExtendedInfo(playerID int) (*models.PlayerExtendedInfo, error) {
	// First get the player info
	player, err := r.GetPlayerByID(playerID)
	if err != nil {
		return nil, err
	}

	info := &models.PlayerExtendedInfo{
		Player:           *player,
		TotalGamesPlayed: 0,
		TotalGamesWon:    0,
	}

	// Count total games played (using latest revision of each game)
	err = r.db.QueryRow(`
		SELECT COUNT(DISTINCT gp.game_id)
		FROM game_player gp
		JOIN game g ON gp.game_id = g.id
		WHERE gp.player_id = ?
		AND g.id IN (
			SELECT id FROM game g2
			WHERE g2.game_id = g.game_id
			ORDER BY g2.revision DESC
			LIMIT 1
		)
	`, playerID).Scan(&info.TotalGamesPlayed)
	if err != nil {
		return nil, err
	}

	// Count total games won (where player has highest total_points in latest revision)
	err = r.db.QueryRow(`
		SELECT COUNT(DISTINCT gp.game_id)
		FROM game_player gp
		JOIN game g ON gp.game_id = g.id
		WHERE gp.player_id = ?
		AND g.id IN (
			SELECT id FROM game g2
			WHERE g2.game_id = g.game_id
			ORDER BY g2.revision DESC
			LIMIT 1
		)
		AND gp.total_points = (
			SELECT MAX(gp2.total_points)
			FROM game_player gp2
			WHERE gp2.game_id = gp.game_id
		)
	`, playerID).Scan(&info.TotalGamesWon)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// GetPlayerByName retrieves a player by name
func (r *Repository) GetPlayerByName(name string) (*models.Player, error) {
	var player models.Player
	err := r.db.QueryRow("SELECT id, name, password_hash, role, created_by, created_at, updated_at FROM player WHERE name = ?", name).Scan(
		&player.ID, &player.Name, &player.PasswordHash, &player.Role, &player.CreatedBy, &player.CreatedAt, &player.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// GetAllPlayers retrieves all players
func (r *Repository) GetAllPlayers() ([]models.Player, error) {
	rows, err := r.db.Query("SELECT id, name, password_hash, role, created_by, created_at, updated_at FROM player ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		err := rows.Scan(&p.ID, &p.Name, &p.PasswordHash, &p.Role, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

// AuthenticatePlayer checks if a player exists and the password is correct
func (r *Repository) AuthenticatePlayer(name string, password string) (*models.Player, error) {
	player, err := r.GetPlayerByName(name)
	if err != nil {
		return nil, err
	}

	if player.PasswordHash == nil {
		return nil, fmt.Errorf("player '%s' has no password set", name)
	}

	if player.Role == models.RolePlayer {
		return nil, fmt.Errorf("players should not have a password set, but cannot be authenticated anyways, player: %s", name)
	}

	if !auth.CheckPassword(password, *player.PasswordHash) {
		return nil, fmt.Errorf("invalid password for player '%s'", name)
	}

	return player, nil
}

// UpdatePlayer updates a player's information with role validation
func (r *Repository) UpdatePlayer(playerID int, name string, password *string, role models.PlayerRole, actor models.Player) (*models.Player, error) {
	currentPlayer, err := r.GetPlayerByID(playerID)
	if err != nil {
		return nil, err
	}

	if err := auth.CanUpdatePlayer(actor, *currentPlayer); err != nil {
		return nil, err
	}

	if err := auth.ValidateRoleTransition(actor, currentPlayer.Role, role); err != nil {
		return nil, err
	}

	if auth.RequiresPassword(role) {
		passwordHash := currentPlayer.PasswordHash
		if password != nil {
			hash, err := auth.HashPassword(*password)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			passwordHash = &hash
		}

		if passwordHash == nil {
			return nil, fmt.Errorf("role '%s' requires a password", role)
		}

		_, err = r.db.Exec("UPDATE player SET name = ?, password_hash = ?, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			name, passwordHash, role, playerID)
		if err != nil {
			return nil, err
		}
	} else {
		if password != nil {
			return nil, fmt.Errorf("role 'player' cannot have a password")
		}

		_, err := r.db.Exec("UPDATE player SET name = ?, password_hash = NULL, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			name, role, playerID)
		if err != nil {
			return nil, err
		}
	}

	// Return the updated player in any case
	return r.GetPlayerByID(playerID)
}

// createNormalGameData creates all game-related data (players, milestones, awards) for a normal game revision
func (r *Repository) createNormalGameData(tx *sql.Tx, gameID int64, req *models.CreateGameRequest) error {
	var gamePlayers []models.GamePlayer
	for _, playerReq := range req.Players {
		// Get player by name
		var playerID int
		err := tx.QueryRow("SELECT id FROM player WHERE name = ?", playerReq.Name).Scan(&playerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("player '%s' not found. Please create the player first", playerReq.Name)
			}
			return fmt.Errorf("error finding player '%s': %v", playerReq.Name, err)
		}

		// Create game_player entry with scores from the player request
		// Milestone, award, and total points will be calculated later
		gamePlayerResult, err := tx.Exec(`
			INSERT INTO game_player (
				game_id, player_id, corporation,
				terraforming_rating, cities, greeneries, cards, turmoil_points,
				milestone_points, award_points, total_points
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, 0, 0)
		`, gameID, playerID, playerReq.Corporation,
			playerReq.TerraformingRating, playerReq.Cities, playerReq.Greeneries,
			playerReq.Cards, playerReq.TurmoilPoints)

		if err != nil {
			return err
		}

		gamePlayerID, err := gamePlayerResult.LastInsertId()
		if err != nil {
			return err
		}

		gamePlayers = append(gamePlayers, models.GamePlayer{
			ID:                 int(gamePlayerID),
			GameID:             int(gameID),
			PlayerID:           playerID,
			Corporation:        &playerReq.Corporation,
			TerraformingRating: playerReq.TerraformingRating,
			Cities:             playerReq.Cities,
			Greeneries:         playerReq.Greeneries,
			Cards:              playerReq.Cards,
			TurmoilPoints:      playerReq.TurmoilPoints,
			MilestonePoints:    0,
			AwardPoints:        0,
			TotalPoints:        0,
		})
	}

	// Create milestones and track them for point calculation
	var milestones []models.Milestone
	for _, milestoneReq := range req.Milestones {
		// Validate winner_game_player_index if provided
		var winnerGamePlayerID *int
		if milestoneReq.WinnerGamePlayerIndex != nil {
			// WinnerGamePlayerIndex in the request is the index of the player in the req.Players array
			if *milestoneReq.WinnerGamePlayerIndex < 0 || *milestoneReq.WinnerGamePlayerIndex >= len(gamePlayers) {
				return fmt.Errorf("invalid winner_game_player_index %d for milestone '%s': must be between 0 and %d",
					*milestoneReq.WinnerGamePlayerIndex, milestoneReq.Name, len(gamePlayers)-1)
			}
			winnerGamePlayerID = &gamePlayers[*milestoneReq.WinnerGamePlayerIndex].ID
		}

		result, err := tx.Exec(`
			INSERT INTO milestone (game_id, name, winner_game_player_id)
			VALUES (?, ?, ?)
		`, gameID, milestoneReq.Name, winnerGamePlayerID)
		if err != nil {
			return fmt.Errorf("error creating milestone '%s': %v", milestoneReq.Name, err)
		}

		milestoneID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		milestones = append(milestones, models.Milestone{
			ID:                 int(milestoneID),
			GameID:             int(gameID),
			Name:               milestoneReq.Name,
			WinnerGamePlayerID: winnerGamePlayerID,
		})
	}

	// Create awards and their placements, tracking for point calculation
	var placements []models.AwardPlacement
	for _, awardReq := range req.Awards {
		// Create the award
		result, err := tx.Exec(`
			INSERT INTO award (game_id, name)
			VALUES (?, ?)
		`, gameID, awardReq.Name)
		if err != nil {
			return fmt.Errorf("error creating award '%s': %v", awardReq.Name, err)
		}

		awardID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Create award placements
		for _, placementReq := range awardReq.Placements {
			// Validate player index
			if placementReq.PlayerIndex < 0 || placementReq.PlayerIndex >= len(gamePlayers) {
				return fmt.Errorf("invalid player_index %d for award '%s': must be between 0 and %d",
					placementReq.PlayerIndex, awardReq.Name, len(gamePlayers)-1)
			}

			// Validate placement value
			if !placementReq.Placement.IsValid() {
				return fmt.Errorf("invalid placement %d for award '%s': must be %d (gold) or %d (silver)",
					placementReq.Placement, awardReq.Name, models.PlacementFirst, models.PlacementSecond)
			}

			// Create the placement
			placementResult, err := tx.Exec(`
				INSERT INTO award_placement (award_id, game_player_id, placement)
				VALUES (?, ?, ?)
			`, awardID, gamePlayers[placementReq.PlayerIndex].ID, placementReq.Placement)
			if err != nil {
				return fmt.Errorf("error creating award placement: %v", err)
			}

			placementID, err := placementResult.LastInsertId()
			if err != nil {
				return err
			}

			placements = append(placements, models.AwardPlacement{
				ID:           int(placementID),
				AwardID:      int(awardID),
				GamePlayerID: gamePlayers[placementReq.PlayerIndex].ID,
				Placement:    placementReq.Placement,
			})
		}
	}

	// Calculate and update total points for each game player (normal games only)
	for i := range gamePlayers {
		gp := &gamePlayers[i]

		// Calculate points from milestones and awards
		milestonePoints := 0
		for _, milestone := range milestones {
			if milestone.WinnerGamePlayerID != nil && *milestone.WinnerGamePlayerID == gp.ID {
				milestonePoints += gamedata.Data.Constants.MilestonePoints
			}
		}

		awardPoints := 0
		for _, placement := range placements {
			if placement.GamePlayerID == gp.ID {
				switch placement.Placement {
				case models.PlacementFirst:
					awardPoints += gamedata.Data.Constants.AwardPointsGold
				case models.PlacementSecond:
					awardPoints += gamedata.Data.Constants.AwardPointsSilver
				}
			}
		}

		// Calculate total points
		totalPoints := gp.TerraformingRating + gp.Cities + gp.Greeneries +
			gp.Cards + gp.TurmoilPoints + milestonePoints + awardPoints

		// Update the game_player record with calculated points
		_, err := tx.Exec(`
			UPDATE game_player
			SET milestone_points = ?, award_points = ?, total_points = ?
			WHERE id = ?
		`, milestonePoints, awardPoints, totalPoints, gp.ID)
		if err != nil {
			return fmt.Errorf("error updating points for game_player %d: %v", gp.ID, err)
		}

		// Update the struct for consistency
		gp.MilestonePoints = milestonePoints
		gp.AwardPoints = awardPoints
		gp.TotalPoints = totalPoints
	}

	return nil
}

// CreateGame creates a new game with all related data
func (r *Repository) CreateGame(req *models.ParsedGameRequest, actor models.Player) (*models.GameWithDetails, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Check if the actor can create games
	if err := auth.CanCreateGames(actor); err != nil {
		return nil, err
	}

	// Determine which type of game to create
	if req.Normal != nil {
		return r.createNormalGame(req.Normal, actor)
	}
	if req.Legacy != nil {
		return r.createLegacyGame(req.Legacy, actor)
	}

	return nil, fmt.Errorf("invalid request: neither normal nor legacy game data provided")
}

// createNormalGame handles creating a non-legacy game
func (r *Repository) createNormalGame(req *models.CreateGameRequest, actor models.Player) (*models.GameWithDetails, error) {
	if len(req.Players) == 0 {
		return nil, fmt.Errorf("game must have at least one player")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Marshal expansions to JSON (always present for normal games)
	expansionsJSON, err := json.Marshal(req.Expansions)
	if err != nil {
		return nil, err
	}

	// Get next game ID using max(game_id) + 1
	var gameID int
	err = tx.QueryRow("SELECT COALESCE(MAX(game_id), 0) + 1 FROM game").Scan(&gameID)
	if err != nil {
		return nil, err
	}

	// Create game with first revision (normal game, legacy_mode = FALSE)
	result, err := tx.Exec(`
		INSERT INTO game (game_id, revision, name, date, map, generations, expansions, note, legacy_mode, created_by)
		VALUES (?, 1, ?, ?, ?, ?, ?, ?, FALSE, ?)
	`, gameID, req.Name, req.Date, req.Map, req.Generations, string(expansionsJSON), req.Note, actor.ID)
	if err != nil {
		return nil, err
	}

	internalID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := r.createNormalGameData(tx, internalID, req); err != nil {
		return nil, err
	}

	for i, imageReq := range req.Images {
		_, err = tx.Exec(`
			INSERT INTO game_image (game_id, image_data, mime_type, display_order)
			VALUES (?, ?, ?, ?)
		`, internalID, imageReq.ImageData, imageReq.MimeType, i)
		if err != nil {
			return nil, fmt.Errorf("error creating image %d: %v", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetGameByID(gameID)
}

// createLegacyGame handles creating a legacy game
func (r *Repository) createLegacyGame(req *models.CreateLegacyGameRequest, actor models.Player) (*models.GameWithDetails, error) {
	if len(req.Players) == 0 {
		return nil, fmt.Errorf("game must have at least one player")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get next game ID using max(game_id) + 1
	var gameID int
	err = tx.QueryRow("SELECT COALESCE(MAX(game_id), 0) + 1 FROM game").Scan(&gameID)
	if err != nil {
		return nil, err
	}

	// Create game with first revision (legacy game, legacy_mode = TRUE)
	result, err := tx.Exec(`
		INSERT INTO game (game_id, revision, name, date, map, generations, expansions, note, legacy_mode, created_by)
		VALUES (?, 1, ?, ?, NULL, NULL, NULL, ?, TRUE, ?)
	`, gameID, req.Name, req.Date, req.Note, actor.ID)
	if err != nil {
		return nil, err
	}

	internalID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create game players
	for _, player := range req.Players {
		// Get player ID
		var playerID int
		err = tx.QueryRow("SELECT id FROM player WHERE name = ?", player.Name).Scan(&playerID)
		if err != nil {
			return nil, fmt.Errorf("player %s not found: %v", player.Name, err)
		}

		totalPoints := player.TerraformingRating + player.Cities + player.Greeneries +
			player.Cards + player.TurmoilPoints + player.MilestonePoints + player.AwardPoints

		_, err = tx.Exec(`
			INSERT INTO game_player (
				game_id, player_id, corporation,
				terraforming_rating, cities, greeneries, cards, turmoil_points,
				milestone_points, award_points, total_points
			) VALUES (?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?)
		`, internalID, playerID,
			player.TerraformingRating, player.Cities, player.Greeneries,
			player.Cards, player.TurmoilPoints, player.MilestonePoints, player.AwardPoints,
			totalPoints)
		if err != nil {
			return nil, fmt.Errorf("error creating game_player for %s: %v", player.Name, err)
		}
	}

	// Create game images
	for i, imageReq := range req.Images {
		_, err = tx.Exec(`
			INSERT INTO game_image (game_id, image_data, mime_type, display_order)
			VALUES (?, ?, ?, ?)
		`, internalID, imageReq.ImageData, imageReq.MimeType, i)
		if err != nil {
			return nil, fmt.Errorf("error creating image %d: %v", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetGameByID(gameID)
}

// UpdateGame creates a new revision of an existing game
func (r *Repository) UpdateGame(gameID int, req *models.ParsedGameRequest, actor models.Player) (*models.GameWithDetails, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	var createdBy int
	err := r.db.QueryRow("SELECT created_by FROM game WHERE game_id = ? LIMIT 1", gameID).Scan(&createdBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game with ID %d not found", gameID)
		}
		return nil, err
	}

	if err := auth.CanModifyGame(actor, createdBy); err != nil {
		return nil, err
	}

	if req.Normal != nil {
		return r.updateNormalGame(gameID, req.Normal, createdBy)
	}
	if req.Legacy != nil {
		return r.updateLegacyGame(gameID, req.Legacy, createdBy)
	}

	return nil, fmt.Errorf("invalid request: neither normal nor legacy game data provided")
}

// updateNormalGame handles updating a non-legacy game
func (r *Repository) updateNormalGame(gameID int, req *models.CreateGameRequest, createdBy int) (*models.GameWithDetails, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get the next revision number
	var maxRevision int
	err = tx.QueryRow("SELECT MAX(revision) FROM game WHERE game_id = ?", gameID).Scan(&maxRevision)
	if err != nil {
		return nil, err
	}

	// Marshal expansions to JSON (always present for normal games)
	expansionsJSON, err := json.Marshal(req.Expansions)
	if err != nil {
		return nil, err
	}

	// Create new revision
	result, err := tx.Exec(`
		INSERT INTO game (game_id, revision, name, date, map, generations, expansions, note, legacy_mode, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, FALSE, ?)
	`, gameID, maxRevision+1, req.Name, req.Date, req.Map, req.Generations, string(expansionsJSON), req.Note, createdBy)
	if err != nil {
		return nil, err
	}

	internalID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := r.createNormalGameData(tx, internalID, req); err != nil {
		return nil, err
	}

	// Create game images
	images := req.Images
	for i, imageReq := range images {
		if imageReq.ID != nil {
			// For "references" just copy the image, it is easier to handle, and we will never actually have too many images
			_, err = tx.Exec(`
				INSERT INTO game_image (game_id, image_data, mime_type, display_order)
				SELECT ?, image_data, mime_type, ?
				FROM game_image
				WHERE id = ?
			`, internalID, i, *imageReq.ID)
			if err != nil {
				return nil, fmt.Errorf("error copying existing image %d: %v", *imageReq.ID, err)
			}
		} else {
			_, err = tx.Exec(`
				INSERT INTO game_image (game_id, image_data, mime_type, display_order)
				VALUES (?, ?, ?, ?)
			`, internalID, imageReq.ImageData, imageReq.MimeType, i)
			if err != nil {
				return nil, fmt.Errorf("error creating image %d: %v", i, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetGameByID(gameID)
}

// updateLegacyGame handles updating a legacy game
func (r *Repository) updateLegacyGame(gameID int, req *models.CreateLegacyGameRequest, createdBy int) (*models.GameWithDetails, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get the next revision number
	var maxRevision int
	err = tx.QueryRow("SELECT MAX(revision) FROM game WHERE game_id = ?", gameID).Scan(&maxRevision)
	if err != nil {
		return nil, err
	}

	// Create new revision
	result, err := tx.Exec(`
		INSERT INTO game (game_id, revision, name, date, map, generations, expansions, note, legacy_mode, created_by)
		VALUES (?, ?, ?, ?, NULL, NULL, NULL, ?, TRUE, ?)
	`, gameID, maxRevision+1, req.Name, req.Date, req.Note, createdBy)
	if err != nil {
		return nil, err
	}

	internalID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create game players for legacy game
	for _, player := range req.Players {
		var playerID int
		err = tx.QueryRow("SELECT id FROM player WHERE name = ?", player.Name).Scan(&playerID)
		if err != nil {
			return nil, fmt.Errorf("player %s not found: %v", player.Name, err)
		}

		totalPoints := player.TerraformingRating + player.Cities + player.Greeneries +
			player.Cards + player.TurmoilPoints + player.MilestonePoints + player.AwardPoints

		_, err = tx.Exec(`
			INSERT INTO game_player (
				game_id, player_id, corporation,
				terraforming_rating, cities, greeneries, cards, turmoil_points,
				milestone_points, award_points, total_points
			) VALUES (?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?)
		`, internalID, playerID,
			player.TerraformingRating, player.Cities, player.Greeneries,
			player.Cards, player.TurmoilPoints, player.MilestonePoints, player.AwardPoints,
			totalPoints)
		if err != nil {
			return nil, fmt.Errorf("error creating game_player for %s: %v", player.Name, err)
		}
	}

	for i, imageReq := range req.Images {
		if imageReq.ID != nil {
			// For "references" just copy the image, it is easier to handle, and we will never actually have too many images
			_, err = tx.Exec(`
				INSERT INTO game_image (game_id, image_data, mime_type, display_order)
				SELECT ?, image_data, mime_type, ?
				FROM game_image
				WHERE id = ?
			`, internalID, i, *imageReq.ID)
			if err != nil {
				return nil, fmt.Errorf("error copying existing image %d: %v", *imageReq.ID, err)
			}
		} else {
			_, err = tx.Exec(`
				INSERT INTO game_image (game_id, image_data, mime_type, display_order)
				VALUES (?, ?, ?, ?)
			`, internalID, imageReq.ImageData, imageReq.MimeType, i)
			if err != nil {
				return nil, fmt.Errorf("error creating image %d: %v", i, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetGameByID(gameID)
}

// GetAllGames retrieves all latest game revisions
func (r *Repository) GetAllGames() ([]models.Game, error) {
	rows, err := r.db.Query(`
		SELECT id, game_id, revision, name, date, map, generations, expansions, note, created_by, created_at
		FROM game
		GROUP BY game_id
		HAVING revision = MAX(revision)
		ORDER BY date DESC, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var g models.Game
		var expansionsJSON sql.NullString
		err := rows.Scan(
			&g.ID, &g.GameID, &g.Revision, &g.Name, &g.Date,
			&g.Map, &g.Generations, &expansionsJSON, &g.Note, &g.CreatedBy, &g.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if expansionsJSON.Valid {
			if err := json.Unmarshal([]byte(expansionsJSON.String), &g.Expansions); err != nil {
				return nil, fmt.Errorf("error parsing expansions JSON: %v", err)
			}
		} else {
			// For legacy games, expansions is nil
			g.Expansions = nil
		}

		games = append(games, g)
	}
	return games, nil
}

// GetGameByID retrieves a specific game revision by ID with all related data
func (r *Repository) GetGameByID(gameID int) (*models.GameWithDetails, error) {
	// Get the game
	var game models.Game
	var expansionsJSON sql.NullString
	err := r.db.QueryRow(`
		SELECT id, game_id, revision, name, date, map, generations, expansions, note, legacy_mode, created_by, created_at
		FROM game
		WHERE game_id = ?
		ORDER BY revision DESC
		LIMIT 1
	`, gameID).Scan(
		&game.ID, &game.GameID, &game.Revision, &game.Name, &game.Date,
		&game.Map, &game.Generations, &expansionsJSON, &game.Note, &game.LegacyMode, &game.CreatedBy, &game.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game with ID %d not found", gameID)
		}
		return nil, err
	}

	if expansionsJSON.Valid {
		if err := json.Unmarshal([]byte(expansionsJSON.String), &game.Expansions); err != nil {
			return nil, fmt.Errorf("error parsing expansions JSON: %v", err)
		}
	} else {
		// For legacy games, expansions is nil
		game.Expansions = nil
	}

	// Get game players and their associated player info (preserve original order)
	rows, err := r.db.Query(`
		SELECT
			gp.id, gp.game_id, gp.player_id, gp.corporation,
			gp.terraforming_rating, gp.cities, gp.greeneries, gp.cards,
			gp.turmoil_points, gp.milestone_points, gp.award_points, gp.total_points,
			p.id, p.name, p.password_hash, p.role, p.created_by, p.created_at, p.updated_at
		FROM game_player gp
		JOIN player p ON gp.player_id = p.id
		WHERE gp.game_id = ?
		ORDER BY gp.id
	`, game.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gamePlayers []models.GamePlayer
	var players []models.Player
	playersSeen := make(map[int]bool)

	for rows.Next() {
		var gp models.GamePlayer
		var p models.Player
		var corporationDB sql.NullString
		err := rows.Scan(
			&gp.ID, &gp.GameID, &gp.PlayerID, &corporationDB,
			&gp.TerraformingRating, &gp.Cities, &gp.Greeneries, &gp.Cards,
			&gp.TurmoilPoints, &gp.MilestonePoints, &gp.AwardPoints, &gp.TotalPoints,
			&p.ID, &p.Name, &p.PasswordHash, &p.Role, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if corporationDB.Valid {
			gp.Corporation = &corporationDB.String
		} else {
			gp.Corporation = nil
		}

		gamePlayers = append(gamePlayers, gp)

		if !playersSeen[p.ID] {
			players = append(players, p)
			playersSeen[p.ID] = true
		}
	}

	milestoneRows, err := r.db.Query(`
		SELECT id, game_id, name, winner_game_player_id
		FROM milestone 
		WHERE game_id = ?
	`, game.ID)
	if err != nil {
		return nil, err
	}
	defer milestoneRows.Close()

	var milestones []models.Milestone
	for milestoneRows.Next() {
		var m models.Milestone
		err := milestoneRows.Scan(&m.ID, &m.GameID, &m.Name, &m.WinnerGamePlayerID)
		if err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}

	awardRows, err := r.db.Query(`
		SELECT id, game_id, name
		FROM award 
		WHERE game_id = ?
	`, game.ID)
	if err != nil {
		return nil, err
	}
	defer awardRows.Close()

	var awards []models.Award
	for awardRows.Next() {
		var a models.Award
		err := awardRows.Scan(&a.ID, &a.GameID, &a.Name)
		if err != nil {
			return nil, err
		}
		awards = append(awards, a)
	}

	placementRows, err := r.db.Query(`
		SELECT ap.id, ap.award_id, ap.game_player_id, ap.placement
		FROM award_placement ap
		JOIN award a ON ap.award_id = a.id
		WHERE a.game_id = ?
	`, game.ID)
	if err != nil {
		return nil, err
	}
	defer placementRows.Close()

	var placements []models.AwardPlacement
	for placementRows.Next() {
		var p models.AwardPlacement
		err := placementRows.Scan(&p.ID, &p.AwardID, &p.GamePlayerID, &p.Placement)
		if err != nil {
			return nil, err
		}
		placements = append(placements, p)
	}

	// Get image metadata
	imageRows, err := r.db.Query(`
		SELECT id, display_order, mime_type, uploaded_at
		FROM game_image
		WHERE game_id = ?
		ORDER BY display_order
	`, game.ID)
	if err != nil {
		return nil, err
	}
	defer imageRows.Close()

	var images []models.GameImageMeta
	for imageRows.Next() {
		var img models.GameImageMeta
		err := imageRows.Scan(&img.ID, &img.DisplayOrder, &img.MimeType, &img.UploadedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return &models.GameWithDetails{
		Game:        game,
		GamePlayers: gamePlayers,
		Players:     players,
		Milestones:  milestones,
		Awards:      awards,
		Placements:  placements,
		Images:      images,
	}, nil
}

// GetGameImageData retrieves the actual image data for a specific image
func (r *Repository) GetGameImageData(imageID int) ([]byte, string, error) {
	var imageData []byte
	var mimeType string
	err := r.db.QueryRow(`
		SELECT image_data, mime_type
		FROM game_image
		WHERE id = ?
	`, imageID).Scan(&imageData, &mimeType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", fmt.Errorf("image with ID %d not found", imageID)
		}
		return nil, "", err
	}
	return imageData, mimeType, nil
}

// GetImageGameID retrieves the user-facing game_id for a specific image
func (r *Repository) GetImageGameID(imageID int) (int, error) {
	var userFacingGameID int
	err := r.db.QueryRow(`
		SELECT game.game_id
		FROM game_image
		JOIN game ON game_image.game_id = game.id
		WHERE game_image.id = ?
	`, imageID).Scan(&userFacingGameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("image with ID %d not found", imageID)
		}
		return 0, err
	}
	return userFacingGameID, nil
}

// GetGamesForRating returns every game's participants chronologically,
// using only the latest revision per game_id. Drives the rating snapshot.
func (r *Repository) GetGamesForRating() ([]rating.GameForRating, error) {
	rows, err := r.db.Query(`
		SELECT g.game_id, g.date, gp.player_id, gp.total_points
		FROM game g
		JOIN game_player gp ON gp.game_id = g.id
		WHERE (g.game_id, g.revision) IN (
			SELECT game_id, MAX(revision) FROM game GROUP BY game_id
		)
		ORDER BY g.date ASC, g.game_id ASC, gp.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []rating.GameForRating
	for rows.Next() {
		var gameID, playerID, totalPoints int
		var date string
		if err := rows.Scan(&gameID, &date, &playerID, &totalPoints); err != nil {
			return nil, err
		}
		// Group consecutive rows belonging to the same game (the query orders by game_id).
		if len(games) == 0 || games[len(games)-1].GameID != gameID {
			games = append(games, rating.GameForRating{GameID: gameID, Date: date})
		}
		g := &games[len(games)-1]
		g.Participants = append(g.Participants, rating.Participant{
			PlayerID:    playerID,
			TotalPoints: totalPoints,
		})
	}
	return games, rows.Err()
}

// DeleteGame deletes all revisions of a game
func (r *Repository) DeleteGame(gameID int, actor models.Player) error {
	var createdBy int
	err := r.db.QueryRow("SELECT created_by FROM game WHERE game_id = ? LIMIT 1", gameID).Scan(&createdBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("game with ID %d not found", gameID)
		}
		return err
	}

	if err := auth.CanModifyGame(actor, createdBy); err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM game WHERE game_id = ?", gameID)
	if err != nil {
		return fmt.Errorf("failed to delete game: %v", err)
	}

	return nil
}
