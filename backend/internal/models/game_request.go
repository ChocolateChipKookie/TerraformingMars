package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	gamedata "terraforming-mars-backend/shared"
)

// Parse don't validate (see logic in in parse function)
type ParsedGameRequest struct {
	Normal *CreateGameRequest
	Legacy *CreateLegacyGameRequest
}

// Utility for fetching images from the contained request
func (r *ParsedGameRequest) GetImages() []ImageRequest {
	if r.Normal != nil {
		return r.Normal.Images
	}
	if r.Legacy != nil {
		return r.Legacy.Images
	}
	return nil
}

// Normal request and its children
type CreateGameRequest struct {
	Name        string             `json:"name"`
	Date        string             `json:"date"`
	Map         string             `json:"map"`
	Generations int                `json:"generations"`
	Expansions  Expansions         `json:"expansions"`
	Note        *string            `json:"note,omitempty"`
	Players     []PlayerRequest    `json:"players"`
	Milestones  []MilestoneRequest `json:"milestones"`
	Awards      []AwardRequest     `json:"awards"`
	Images      []ImageRequest     `json:"images"`
}

type PlayerRequest struct {
	Name               string `json:"name"`
	Corporation        string `json:"corporation"`
	TerraformingRating int    `json:"terraforming_rating"`
	Cities             int    `json:"cities"`
	Greeneries         int    `json:"greeneries"`
	Cards              int    `json:"cards"`
	TurmoilPoints      int    `json:"turmoil_points"`
}

type MilestoneRequest struct {
	Name                  string `json:"name"`
	WinnerGamePlayerIndex *int   `json:"winner_game_player_index"`
}

type AwardRequest struct {
	Name       string             `json:"name"`
	Placements []PlacementRequest `json:"placements"`
}

type PlacementRequest struct {
	PlayerIndex int       `json:"player_index"`
	Placement   Placement `json:"placement"`
}

// Legacy request and its children
type CreateLegacyGameRequest struct {
	Name    string                `json:"name"`
	Date    string                `json:"date"`
	Note    *string               `json:"note,omitempty"`
	Players []LegacyPlayerRequest `json:"players"`
	Images  []ImageRequest        `json:"images"`
}

type LegacyPlayerRequest struct {
	Name               string `json:"name"`
	TerraformingRating int    `json:"terraforming_rating"`
	Cities             int    `json:"cities"`
	Greeneries         int    `json:"greeneries"`
	Cards              int    `json:"cards"`
	TurmoilPoints      int    `json:"turmoil_points"`
	MilestonePoints    int    `json:"milestone_points"`
	AwardPoints        int    `json:"award_points"`
}

// Common
type ImageRequest struct {
	// For new images (empty if existing image)
	ImageData []byte `json:"image_data,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`

	// For existing images (reference by ID)
	ID *int `json:"id,omitempty"`
}

func ParseCreateGameRequest(r io.Reader, isUpdate bool) (*CreateGameRequest, error) {
	var req CreateGameRequest

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Validate common data
	if req.Name == "" {
		return nil, errors.New("game name is required")
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}
	mapData, err := gamedata.ValidateMap(req.Map)
	if err != nil {
		return nil, err
	}
	if err := gamedata.ValidateGenerations(req.Generations); err != nil {
		return nil, err
	}
	if err := gamedata.ValidateExpansions(req.Expansions); err != nil {
		return nil, err
	}
	if err := gamedata.ValidatePlayerCount(len(req.Players), mapData); err != nil {
		return nil, err
	}

	// Check for duplicate player names and validate each player
	playerNames := make(map[string]bool)
	for i, player := range req.Players {
		if player.Name == "" {
			return nil, fmt.Errorf("player %d: name is required", i+1)
		}
		if playerNames[player.Name] {
			return nil, fmt.Errorf("duplicate player name: %s", player.Name)
		}
		playerNames[player.Name] = true

		if player.Corporation == "" {
			return nil, fmt.Errorf("player %s: corporation is required", player.Name)
		}

		// Validate corporation is available with selected expansions
		if err := gamedata.ValidateCorporation(player.Corporation, req.Expansions); err != nil {
			return nil, fmt.Errorf("player %s: %w", player.Name, err)
		}

		// Validate score ranges
		if player.TerraformingRating < 0 || player.TerraformingRating > 200 {
			return nil, fmt.Errorf("player %s: terraforming rating must be 0-200", player.Name)
		}
		if player.Cities < 0 || player.Cities > 100 {
			return nil, fmt.Errorf("player %s: cities points must be 0-100", player.Name)
		}
		if player.Greeneries < 0 || player.Greeneries > 100 {
			return nil, fmt.Errorf("player %s: greeneries points must be 0-100", player.Name)
		}
		if player.Cards < -50 || player.Cards > 100 {
			return nil, fmt.Errorf("player %s: cards points must be -50-100", player.Name)
		}
		if player.TurmoilPoints < 0 || player.TurmoilPoints > 100 {
			return nil, fmt.Errorf("player %s: turmoil points must be 0 to 100", player.Name)
		}
	}

	// Validate milestone and award counts
	milestoneNamesForCount := make([]string, len(req.Milestones))
	for i, milestone := range req.Milestones {
		milestoneNamesForCount[i] = milestone.Name
	}
	if err := gamedata.ValidateMilestoneCount(milestoneNamesForCount, req.Expansions); err != nil {
		return nil, err
	}

	awardNamesForCount := make([]string, len(req.Awards))
	for i, award := range req.Awards {
		awardNamesForCount[i] = award.Name
	}
	if err := gamedata.ValidateAwardCount(awardNamesForCount, req.Expansions); err != nil {
		return nil, err
	}

	// Validate milestones
	useCustomMilestones := req.Expansions["Milestones & Awards"]
	milestoneNames := make(map[string]bool)
	for _, milestone := range req.Milestones {
		if milestone.Name == "" {
			return nil, errors.New("milestone name is required")
		}
		if milestoneNames[milestone.Name] {
			return nil, fmt.Errorf("duplicate milestone: %s", milestone.Name)
		}
		milestoneNames[milestone.Name] = true

		if err := gamedata.ValidateMilestone(milestone.Name, req.Map, useCustomMilestones, req.Expansions); err != nil {
			return nil, err
		}

		if milestone.WinnerGamePlayerIndex != nil {
			if *milestone.WinnerGamePlayerIndex < 0 || *milestone.WinnerGamePlayerIndex >= len(req.Players) {
				return nil, fmt.Errorf("milestone %s: invalid winner player index %d",
					milestone.Name, *milestone.WinnerGamePlayerIndex)
			}
		}
	}

	// Validate awards
	useCustomAwards := req.Expansions["Milestones & Awards"]
	awardNames := make(map[string]bool)
	for _, award := range req.Awards {
		if award.Name == "" {
			return nil, errors.New("award name is required")
		}
		if awardNames[award.Name] {
			return nil, fmt.Errorf("duplicate award: %s", award.Name)
		}
		awardNames[award.Name] = true

		if err := gamedata.ValidateAward(award.Name, req.Map, useCustomAwards, req.Expansions); err != nil {
			return nil, err
		}

		playerPlacements := make(map[int]bool)
		for _, placement := range award.Placements {
			if placement.PlayerIndex < 0 || placement.PlayerIndex >= len(req.Players) {
				return nil, fmt.Errorf("award %s: invalid player index %d",
					award.Name, placement.PlayerIndex)
			}
			if playerPlacements[placement.PlayerIndex] {
				return nil, fmt.Errorf("award %s: player %d has multiple placements",
					award.Name, placement.PlayerIndex)
			}
			playerPlacements[placement.PlayerIndex] = true

			if !placement.Placement.IsValid() {
				return nil, fmt.Errorf("award %s: invalid placement value %d",
					award.Name, placement.Placement)
			}
		}
	}

	// Validate images
	if err := validateImages(req.Images, isUpdate); err != nil {
		return nil, err
	}

	return &req, nil
}

func ParseCreateLegacyGameRequest(r io.Reader, isUpdate bool) (*CreateLegacyGameRequest, error) {
	var req CreateLegacyGameRequest

	// Decode JSON
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Validate required fields
	if req.Name == "" {
		return nil, errors.New("game name is required")
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}

	// Assume that all legacy games were played on the Tharsis map
	mapData, err := gamedata.ValidateMap("Tharsis")
	if err != nil {
		return nil, err
	}
	if err := gamedata.ValidatePlayerCount(len(req.Players), mapData); err != nil {
		return nil, err
	}

	// Check for duplicate player names and validate each player
	playerNames := make(map[string]bool)
	for i, player := range req.Players {
		if player.Name == "" {
			return nil, fmt.Errorf("player %d: name is required", i+1)
		}
		if playerNames[player.Name] {
			return nil, fmt.Errorf("duplicate player name: %s", player.Name)
		}
		playerNames[player.Name] = true

		// Validate score ranges
		if player.TerraformingRating < 0 || player.TerraformingRating > 100 {
			return nil, fmt.Errorf("player %s: terraforming rating must be 0 to 100", player.Name)
		}
		if player.Cities < 0 || player.Cities > 100 {
			return nil, fmt.Errorf("player %s: cities points must be 0-100", player.Name)
		}
		if player.Greeneries < 0 || player.Greeneries > 100 {
			return nil, fmt.Errorf("player %s: greeneries must be 0-100", player.Name)
		}
		if player.Cards < -50 || player.Cards > 100 {
			return nil, fmt.Errorf("player %s: cards must be -50 to 100", player.Name)
		}
		if player.TurmoilPoints < 0 || player.TurmoilPoints > 100 {
			return nil, fmt.Errorf("player %s: turmoil points must be 0 to 100", player.Name)
		}

		// Legacy mode specific: validate milestone and award points
		if player.MilestonePoints < 0 || player.MilestonePoints > 15 {
			return nil, fmt.Errorf("player %s: milestone points must be 0-15", player.Name)
		}
		if player.AwardPoints < 0 || player.AwardPoints > 15 {
			return nil, fmt.Errorf("player %s: award points must be 0-15", player.Name)
		}
	}

	// Validate images (reuse the same validation)
	if err := validateImages(req.Images, isUpdate); err != nil {
		return nil, err
	}

	return &req, nil
}

func validateImages(images []ImageRequest, isUpdate bool) error {
	for i, img := range images {
		hasNewImage := len(img.ImageData) > 0 || img.MimeType != ""
		hasExistingRef := img.ID != nil

		if !isUpdate && hasExistingRef {
			return fmt.Errorf("image %d: cannot reference existing images when creating a new game", i+1)
		}
		if hasNewImage && hasExistingRef {
			return fmt.Errorf("image %d: cannot specify both new image data and existing image ID", i+1)
		}
		if !hasNewImage && !hasExistingRef {
			return fmt.Errorf("image %d: must specify either new image data or existing image ID", i+1)
		}

		if hasNewImage {
			if len(img.ImageData) == 0 {
				return fmt.Errorf("image %d: image data is required for new images", i+1)
			}
			if img.MimeType == "" {
				return fmt.Errorf("image %d: mime type is required for new images", i+1)
			}

			validMimeTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/webp": true,
				"image/gif":  true,
			}
			if !validMimeTypes[img.MimeType] {
				return fmt.Errorf("image %d: invalid mime type %s (supported: jpeg, png, webp, gif)", i+1, img.MimeType)
			}
		}
	}
	return nil
}

// We have 2 possible requests coming in, legacy and normal requests
// This function parses the incoming request, discerns if it is legacy/normal, and parses it as such
// While parsing it also validates all the static data (eg. constrained values like enums/map names)
// As a result we have an union type that is either an legacy or normal request (or parse error)
// And the caller can be sure that the returned value is semantically correct, making later processing steps easier because we only have to validate at this one point
// Also we have the type of the request ingrained in the actual type
// Another validation step has to be done for the non-static data (eg. for reference to some IDs in the database)
func ParseGameRequest(r io.Reader, isUpdate bool) (*ParsedGameRequest, error) {
	// Read all bytes first so we can peek at legacy_mode
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	// Check if it's legacy mode
	var check struct {
		LegacyMode bool `json:"legacy_mode"`
	}
	if err := json.Unmarshal(data, &check); err != nil {
		return nil, fmt.Errorf("failed to check legacy_mode: %w", err)
	}

	// Create a new reader from the data
	reader := bytes.NewReader(data)

	// Parse as legacy or normal based on the flag
	if check.LegacyMode {
		validated, err := ParseCreateLegacyGameRequest(reader, isUpdate)
		if err != nil {
			return nil, err
		}
		return &ParsedGameRequest{Legacy: validated}, nil
	}

	// Parse as normal game
	validated, err := ParseCreateGameRequest(reader, isUpdate)
	if err != nil {
		return nil, err
	}
	return &ParsedGameRequest{Normal: validated}, nil
}
