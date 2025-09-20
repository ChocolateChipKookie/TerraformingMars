package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"terraforming-mars-backend/internal/auth"
	"terraforming-mars-backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreatePlayer creates a new player with role-based validation
func (r *Repository) CreatePlayer(name string, password *string, role models.PlayerRole, actor models.Player) (*models.Player, error) {
	// Check if the actor can create players
	if err := auth.CanCreatePlayers(actor, role); err != nil {
		return nil, err
	}

	return r.createPlayer(name, password, role, &actor.ID)
}

// CreateSystemAdmin creates the initial system admin (no actor required)
func (r *Repository) CreateSystemAdmin(name string, password *string) (*models.Player, error) {
	return r.createPlayer(name, password, models.RoleAdmin, nil)
}

// createPlayer is the internal implementation shared by CreatePlayer and CreateSystemAdmin
func (r *Repository) createPlayer(name string, password *string, role models.PlayerRole, createdBy *int) (*models.Player, error) {
	if err := auth.IsValidRole(role); err != nil {
		return nil, err
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

	// Check if player has no password set
	if player.PasswordHash == nil {
		return nil, fmt.Errorf("player '%s' has no password set", name)
	}

	// Players cannot be authenticated, they should not have any password anyways
	if player.Role == models.RolePlayer {
		return nil, fmt.Errorf("players should not have a password set, but cannot be authenticated anyways, player: %s", name)
	}

	// Check password
	if !auth.CheckPassword(password, *player.PasswordHash) {
		return nil, fmt.Errorf("invalid password for player '%s'", name)
	}

	return player, nil
}

// UpdatePlayer updates a player's information with role validation
func (r *Repository) UpdatePlayer(playerID int, name string, password *string, role *models.PlayerRole, actor models.Player) (*models.Player, error) {
	// Get current player info
	currentPlayer, err := r.GetPlayerByID(playerID)
	if err != nil {
		return nil, err
	}

	// Check if the actor can update this player
	if err := auth.CanUpdatePlayer(actor, *currentPlayer); err != nil {
		return nil, err
	}

	// Check if name is being changed and if actor can update names
	if name != currentPlayer.Name {
		if err := auth.CanUpdatePlayerName(actor, *currentPlayer); err != nil {
			return nil, err
		}
	}

	// Determine final role
	finalRole := currentPlayer.Role
	if role != nil {
		// Validate new role
		if err := auth.IsValidRole(*role); err != nil {
			return nil, err
		}
		finalRole = *role
		// Validate role transition
		if err := auth.ValidateRoleTransition(actor, currentPlayer.Role, finalRole); err != nil {
			return nil, err
		}
	}

	if auth.RequiresPassword(finalRole) {
		passwordHash := currentPlayer.PasswordHash
		if password != nil {
			hash, err := auth.HashPassword(*password)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			passwordHash = &hash
		}

		// Password does not already exist nor was it supplied
		if passwordHash == nil {
			return nil, fmt.Errorf("role '%s' requires a password", finalRole)
		}

		_, err = r.db.Exec("UPDATE player SET name = ?, password_hash = ?, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			name, passwordHash, finalRole, playerID)
		if err != nil {
			return nil, err
		}
	} else {
		if password != nil {
			return nil, fmt.Errorf("role 'player' cannot have a password")
		}
		// If changing to role that does not need password, clear existing password
		_, err := r.db.Exec("UPDATE player SET name = ?, password_hash = NULL, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			name, finalRole, playerID)
		if err != nil {
			return nil, err
		}
	}

	// Return the new player in any case
	return r.GetPlayerByID(playerID)
}

// createGameData creates all game-related data (players, milestones, awards) for a game revision
// This is used by both CreateGame and UpdateGame
func (r *Repository) createGameData(tx *sql.Tx, gameID int64, req models.CreateGameRequest) error {
	// Validate that we have at least one player
	if len(req.Players) == 0 {
		return fmt.Errorf("validation error: at least one player is required to create a game")
	}

	// Create game_players entries
	var gamePlayers []models.GamePlayer
	for _, playerReq := range req.Players {
		// Get player by name (must exist)
		var playerID int
		err := tx.QueryRow("SELECT id FROM player WHERE name = ?", playerReq.Name).Scan(&playerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("player '%s' not found. Please create the player first", playerReq.Name)
			}
			return fmt.Errorf("error finding player '%s': %v", playerReq.Name, err)
		}

		// Create game_player entry with scores from the player request
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
			Corporation:        playerReq.Corporation,
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

	// Create milestones
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

	// Create awards and their placements
	var awards []models.Award
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

		awards = append(awards, models.Award{
			ID:     int(awardID),
			GameID: int(gameID),
			Name:   awardReq.Name,
		})

		// Create award placements
		for _, placementReq := range awardReq.Placements {
			// Validate player index
			if placementReq.PlayerIndex < 0 || placementReq.PlayerIndex >= len(gamePlayers) {
				return fmt.Errorf("invalid player_index %d for award '%s': must be between 0 and %d",
					placementReq.PlayerIndex, awardReq.Name, len(gamePlayers)-1)
			}

			// Validate placement value
			if !placementReq.Placement.IsValid() {
				return fmt.Errorf("invalid placement %d for award '%s': must be 1 (first) or 2 (second)",
					placementReq.Placement, awardReq.Name)
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

	// Calculate and update points for each game player
	for i := range gamePlayers {
		gp := &gamePlayers[i]

		// Calculate milestone points (5 points per milestone won)
		milestonePoints := 0
		for _, milestone := range milestones {
			if milestone.WinnerGamePlayerID != nil && *milestone.WinnerGamePlayerID == gp.ID {
				milestonePoints += 5
			}
		}

		// Calculate award points (5 for first, 2 for second)
		awardPoints := 0
		for _, placement := range placements {
			if placement.GamePlayerID == gp.ID {
				switch placement.Placement {
				case models.PlacementFirst:
					awardPoints += 5
				case models.PlacementSecond:
					awardPoints += 2
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
func (r *Repository) CreateGame(req models.CreateGameRequest, actor models.Player) (*models.GameWithDetails, error) {
	// Check if the actor can create games
	if err := auth.CanCreateGames(actor); err != nil {
		return nil, err
	}

	// Ensure the actor is creating the game (set CreatedBy to actor's ID)
	req.CreatedBy = actor.ID

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Validate date format (but store as string)
	_, err = time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v (expected YYYY-MM-DD)", err)
	}

	// Marshal expansions to JSON
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

	// Create game with first revision
	result, err := tx.Exec(`
		INSERT INTO game (game_id, revision, name, date, map, generations, expansions, note, created_by)
		VALUES (?, 1, ?, ?, ?, ?, ?, ?, ?)
	`, gameID, req.Name, req.Date, req.Map, req.Generations, string(expansionsJSON), req.Note, req.CreatedBy)
	if err != nil {
		return nil, err
	}

	// Get the auto-increment primary key for foreign key references
	internalID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create all game-related data
	if err := r.createGameData(tx, internalID, req); err != nil {
		return nil, err
	}

	// Create game images (images should already be processed by API layer)
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

	// Return the created game with all its details
	return r.GetGameByID(gameID)
}

// UpdateGame creates a new revision of an existing game (for API use with permission checking)
func (r *Repository) UpdateGame(gameID int, req models.CreateGameRequest, actor models.Player) (*models.GameWithDetails, error) {
	// First, get the created_by of the game we're updating (from any revision)
	var createdBy int
	err := r.db.QueryRow("SELECT created_by FROM game WHERE game_id = ? LIMIT 1", gameID).Scan(&createdBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game with ID %d not found", gameID)
		}
		return nil, err
	}

	// Check if the actor can modify this game
	if err := auth.CanModifyGame(actor, createdBy); err != nil {
		return nil, err
	}

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

	// Validate date format (but store as string)
	_, err = time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	// Marshal expansions to JSON
	expansionsJSON, err := json.Marshal(req.Expansions)
	if err != nil {
		return nil, err
	}

	// Create new revision (keep the original creator)
	result, err := tx.Exec(`
		INSERT INTO game (game_id, revision, name, date, map, generations, expansions, note, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, gameID, maxRevision+1, req.Name, req.Date, req.Map, req.Generations, string(expansionsJSON), req.Note, createdBy)
	if err != nil {
		return nil, err
	}

	// Get the auto-increment primary key for foreign key references
	internalID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create all game-related data using the shared helper
	if err := r.createGameData(tx, internalID, req); err != nil {
		return nil, err
	}

	// Create game images (images should already be processed by API layer)
	for i, imageReq := range req.Images {
		if imageReq.ID != nil {
			// For existing images, just create a reference (copy the image to new revision)
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
			// For new images, insert the processed data
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
		var expansionsJSON string
		err := rows.Scan(
			&g.ID, &g.GameID, &g.Revision, &g.Name, &g.Date,
			&g.Map, &g.Generations, &expansionsJSON, &g.Note, &g.CreatedBy, &g.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse expansions JSON
		if err := json.Unmarshal([]byte(expansionsJSON), &g.Expansions); err != nil {
			return nil, fmt.Errorf("error parsing expansions JSON: %v", err)
		}

		games = append(games, g)
	}
	return games, nil
}

// GetGameByID retrieves a specific game revision by ID with all related data
func (r *Repository) GetGameByID(gameID int) (*models.GameWithDetails, error) {
	// Get the game
	var game models.Game
	var expansionsJSON string
	err := r.db.QueryRow(`
		SELECT id, game_id, revision, name, date, map, generations, expansions, note, created_by, created_at
		FROM game 
		WHERE game_id = ?
		ORDER BY revision DESC
		LIMIT 1
	`, gameID).Scan(
		&game.ID, &game.GameID, &game.Revision, &game.Name, &game.Date,
		&game.Map, &game.Generations, &expansionsJSON, &game.Note, &game.CreatedBy, &game.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game with ID %d not found", gameID)
		}
		return nil, err
	}

	// Parse expansions JSON
	if err := json.Unmarshal([]byte(expansionsJSON), &game.Expansions); err != nil {
		return nil, fmt.Errorf("error parsing expansions JSON: %v", err)
	}

	// Get game players and their associated player info (preserve original order)
	rows, err := r.db.Query(`
		SELECT 
			gp.id, gp.game_id, gp.player_id, gp.corporation,
			gp.terraforming_rating, gp.cities, gp.greeneries, gp.cards,
			gp.turmoil_points, gp.milestone_points, gp.award_points, gp.total_points,
			p.id, p.name
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
		err := rows.Scan(
			&gp.ID, &gp.GameID, &gp.PlayerID, &gp.Corporation,
			&gp.TerraformingRating, &gp.Cities, &gp.Greeneries, &gp.Cards,
			&gp.TurmoilPoints, &gp.MilestonePoints, &gp.AwardPoints, &gp.TotalPoints,
			&p.ID, &p.Name,
		)
		if err != nil {
			return nil, err
		}

		gamePlayers = append(gamePlayers, gp)

		// Add player to list if not already added
		if !playersSeen[p.ID] {
			players = append(players, p)
			playersSeen[p.ID] = true
		}
	}

	// Get milestones
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

	// Get awards
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

	// Get award placements
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

	// Get image metadata (without actual image data)
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

// GetImageGameID retrieves the user-facing game_id for a specific image (for ownership validation)
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
