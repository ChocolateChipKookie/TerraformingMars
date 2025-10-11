package database

import (
	"fmt"
	"terraforming-mars-backend/internal/models"
)

// ValidateGameRequest performs domain validation on a parsed game request
// This includes checking that referenced entities exist in the database
func (r *Repository) ValidateGameRequest(req *models.ParsedGameRequest) error {
	if req.Normal != nil {
		return r.validateNormalGameRequest(req.Normal)
	}
	if req.Legacy != nil {
		return r.validateLegacyGameRequest(req.Legacy)
	}
	return fmt.Errorf("invalid request: neither normal nor legacy game data provided")
}

// validateNormalGameRequest validates a normal game request against the database
func (r *Repository) validateNormalGameRequest(req *models.CreateGameRequest) error {
	// Extract player names for validation
	playerNames := make([]string, len(req.Players))
	for i, p := range req.Players {
		playerNames[i] = p.Name
	}

	return r.validateCommon(playerNames, req.Images)
}

// validateLegacyGameRequest validates a legacy game request against the database
func (r *Repository) validateLegacyGameRequest(req *models.CreateLegacyGameRequest) error {
	// Extract player names for validation
	playerNames := make([]string, len(req.Players))
	for i, p := range req.Players {
		playerNames[i] = p.Name
	}

	return r.validateCommon(playerNames, req.Images)
}

// validateCommon performs validation common to both normal and legacy games
func (r *Repository) validateCommon(playerNames []string, images []models.ImageRequest) error {
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

	// For game creation, no image IDs should be present (already validated structurally)
	// For game updates, image validation is done in ValidateGameUpdateRequest

	return nil
}

// ValidateGameUpdateRequest validates a game update request
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

	// Get images from the request
	var images []models.ImageRequest
	if req.Normal != nil {
		images = req.Normal.Images
	} else if req.Legacy != nil {
		images = req.Legacy.Images
	}

	// Check that referenced images belong to this game (any revision of it)
	for i, img := range images {
		if img.ID != nil {
			var belongsToGame bool
			// Check if image belongs to any revision of this game by matching the game_id from the game table
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

	// Perform regular validation (player existence)
	return r.ValidateGameRequest(req)
}