package database

import (
	"fmt"
	"terraforming-mars-backend/internal/models"
)

// ValidateGameRequest does domain vaildation of the request
func (r *Repository) ValidateGameRequest(req *models.ParsedGameRequest) error {
	if req.Normal != nil {
		return r.validateNormalGameRequest(req.Normal)
	}
	if req.Legacy != nil {
		return r.validateLegacyGameRequest(req.Legacy)
	}
	return fmt.Errorf("invalid request: neither normal nor legacy game data provided")
}

func (r *Repository) validateNormalGameRequest(req *models.CreateGameRequest) error {
	playerNames := make([]string, len(req.Players))
	for i, p := range req.Players {
		playerNames[i] = p.Name
	}

	return r.validateCommon(playerNames)
}

func (r *Repository) validateLegacyGameRequest(req *models.CreateLegacyGameRequest) error {
	playerNames := make([]string, len(req.Players))
	for i, p := range req.Players {
		playerNames[i] = p.Name
	}

	return r.validateCommon(playerNames)
}

// validateCommon performs validation common to both normal and legacy games
func (r *Repository) validateCommon(playerNames []string) error {
	// Validate all player names exist
	for _, name := range playerNames {
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM player WHERE name = ?)", name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check player %s: %w", name, err)
		}
		if !exists {
			return fmt.Errorf("player %s does not exist", name)
		}
	}

	return nil
}

// Checks that the game exists and that referenced images belong to that game
func (r *Repository) ValidateGameUpdateRequest(gameID int, req *models.ParsedGameRequest) error {
	// Check if game exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM game WHERE game_id = ?)", gameID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check game existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("game with ID %d not found", gameID)
	}

	images := req.GetImages()

	// Check that referenced images belong to this game (any revision of it)
	for i, img := range images {
		if img.ID != nil {
			var belongsToGame bool
			err := r.db.QueryRow(
				`SELECT EXISTS(
					SELECT 1 FROM game_image gi
					JOIN game g1 ON gi.game_id = g1.id
					JOIN game g2 ON g1.game_id = g2.game_id
					WHERE gi.id = ? AND g2.game_id = ?
				)`,
				*img.ID, gameID,
			).Scan(&belongsToGame)
			if err != nil {
				return fmt.Errorf("failed to check image ownership: %w", err)
			}
			if !belongsToGame {
				return fmt.Errorf("image %d: image ID %d does not belong to game %d", i+1, *img.ID, gameID)
			}
		}
	}

	return r.ValidateGameRequest(req)
}

