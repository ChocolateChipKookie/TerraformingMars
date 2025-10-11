package gamedata

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed game-data.json
var gameDataJSON []byte

// GameData represents the shared game configuration
type GameData struct {
	Constants     Constants           `json:"constants"`
	Maps          map[string]MapData  `json:"maps"`
	Expansions    []string            `json:"expansions"`
	Corporations  map[string][]string `json:"corporations"`
	AllMilestones []string            `json:"allMilestones"`
	AllAwards     []string            `json:"allAwards"`
}

// Constants holds game rule constants
type Constants struct {
	MaxMilestonesClaimed  int `json:"maxMilestonesClaimed"`
	MaxAwardsFunded       int `json:"maxAwardsFunded"`
	MilestonePoints       int `json:"milestonePoints"`
	AwardPointsGold       int `json:"awardPointsGold"`
	AwardPointsSilver     int `json:"awardPointsSilver"`
	AwardPlacementNone    int `json:"awardPlacementNone"`
	AwardPlacementGold    int `json:"awardPlacementGold"`
	AwardPlacementSilver  int `json:"awardPlacementSilver"`
	DefaultMilestoneSlots int `json:"defaultMilestoneSlots"`
	VenusMilestoneSlots   int `json:"venusMilestoneSlots"`
	DefaultAwardSlots     int `json:"defaultAwardSlots"`
	VenusAwardSlots       int `json:"venusAwardSlots"`
	MinPlayers            int `json:"minPlayers"`
	DefaultMaxPlayers     int `json:"defaultMaxPlayers"`
	MinGenerations        int `json:"minGenerations"`
	MaxGenerations        int `json:"maxGenerations"`
	DefaultGenerations    int `json:"defaultGenerations"`
	DefaultPlayerCount    int `json:"defaultPlayerCount"`
	DefaultTR             int `json:"defaultTR"`
}

// MapData represents configuration for a specific map
type MapData struct {
	MaxPlayers int      `json:"maxPlayers"`
	Milestones []string `json:"milestones"`
	Awards     []string `json:"awards"`
}

var Data *GameData

// init loads the game data at startup
func init() {
	Data = &GameData{}
	if err := json.Unmarshal(gameDataJSON, Data); err != nil {
		panic(fmt.Sprintf("failed to load game data: %v", err))
	}
}

// ValidateMap checks if a map exists and returns its data
func ValidateMap(mapName string) (*MapData, error) {
	mapData, exists := Data.Maps[mapName]
	if !exists {
		return nil, fmt.Errorf("invalid map: %s", mapName)
	}
	return &mapData, nil
}

// ValidateGenerations checks if generations count is valid
func ValidateGenerations(generations int) error {
	if generations < Data.Constants.MinGenerations || generations > Data.Constants.MaxGenerations {
		return fmt.Errorf("generations must be between %d and %d, got %d",
			Data.Constants.MinGenerations, Data.Constants.MaxGenerations, generations)
	}
	return nil
}

// ValidatePlayerCount checks if player count is valid for the given map
func ValidatePlayerCount(playerCount int, mapData *MapData) error {
	if playerCount < Data.Constants.MinPlayers {
		return fmt.Errorf("game must have at least %d player(s), got %d",
			Data.Constants.MinPlayers, playerCount)
	}
	if playerCount > mapData.MaxPlayers {
		return fmt.Errorf("maximum %d players allowed, got %d",
			mapData.MaxPlayers, playerCount)
	}
	return nil
}

// ValidateExpansions checks if all expansion names are valid and ensures Base Game is enabled
func ValidateExpansions(expansions map[string]bool) error {
	// Base Game must always be enabled
	if !expansions["Base Game"] {
		return fmt.Errorf("Base Game expansion must always be enabled")
	}

	for expansion := range expansions {
		// Check if expansion exists in our data
		validExpansion := false
		for _, e := range Data.Expansions {
			if e == expansion {
				validExpansion = true
				break
			}
		}
		if !validExpansion {
			return fmt.Errorf("invalid expansion: %s", expansion)
		}
	}
	return nil
}

// ValidateCorporation checks if a corporation is valid for the given expansions
func ValidateCorporation(corporation string, enabledExpansions map[string]bool) error {
	// Beginner is always valid
	if corporation == "Beginner" {
		return nil
	}

	// Check if corporation exists in any enabled expansion
	for expansion, corps := range Data.Corporations {
		if enabledExpansions[expansion] {
			for _, corp := range corps {
				if corp == corporation {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("corporation '%s' is not available with selected expansions", corporation)
}

// ValidateMilestone checks if a milestone is valid
func ValidateMilestone(milestone string, mapName string, useCustom bool, expansions map[string]bool) error {
	// Check for Venus Next specific milestone without Venus Next expansion
	if milestone == "Hoverlord" && !expansions["Venus Next"] {
		return fmt.Errorf("Hoverlord is only available with Venus Next")
	}

	if useCustom {
		for _, m := range Data.AllMilestones {
			if m == milestone {
				return nil
			}
		}
		return fmt.Errorf("invalid milestone: %s", milestone)
	}

	mapData, exists := Data.Maps[mapName]
	if !exists {
		return fmt.Errorf("invalid map: %s", mapName)
	}

	for _, m := range mapData.Milestones {
		if m == milestone {
			return nil
		}
	}

	// Special case: Hoverlord is allowed with Venus Next even if not on map
	if milestone == "Hoverlord" && expansions["Venus Next"] {
		return nil
	}

	return fmt.Errorf("milestone '%s' is not available for map %s", milestone, mapName)
}

// ValidateAward checks if an award is valid
func ValidateAward(award string, mapName string, useCustom bool, expansions map[string]bool) error {
	// Check for Venus Next specific award without Venus Next expansion
	if award == "Venuphile" && !expansions["Venus Next"] {
		return fmt.Errorf("Venuphile is only available with Venus Next")
	}

	if useCustom {
		for _, a := range Data.AllAwards {
			if a == award {
				return nil
			}
		}
		return fmt.Errorf("invalid award: %s", award)
	}

	mapData, exists := Data.Maps[mapName]
	if !exists {
		return fmt.Errorf("invalid map: %s", mapName)
	}

	for _, a := range mapData.Awards {
		if a == award {
			return nil
		}
	}

	// Special case: Venuphile is allowed with Venus Next even if not on map
	if award == "Venuphile" && expansions["Venus Next"] {
		return nil
	}

	return fmt.Errorf("award '%s' is not available for map %s", award, mapName)
}

// ValidateMilestoneCount checks if the milestone count is correct based on expansions
func ValidateMilestoneCount(milestones []string, expansions map[string]bool) error {
	expectedCount := Data.Constants.DefaultMilestoneSlots
	if expansions["Venus Next"] {
		expectedCount = Data.Constants.VenusMilestoneSlots
	}

	if len(milestones) != expectedCount {
		return fmt.Errorf("expected %d milestones, got %d", expectedCount, len(milestones))
	}

	// If Venus Next is enabled, Hoverlord is mandatory
	if expansions["Venus Next"] {
		hasHoverlord := false
		for _, milestone := range milestones {
			if milestone == "Hoverlord" {
				hasHoverlord = true
				break
			}
		}
		if !hasHoverlord {
			return fmt.Errorf("Hoverlord milestone is mandatory when Venus Next expansion is enabled")
		}
	}

	return nil
}

// ValidateAwardCount checks if the award count is correct based on expansions
func ValidateAwardCount(awards []string, expansions map[string]bool) error {
	expectedCount := Data.Constants.DefaultAwardSlots
	if expansions["Venus Next"] {
		expectedCount = Data.Constants.VenusAwardSlots
	}

	if len(awards) != expectedCount {
		return fmt.Errorf("expected %d awards, got %d", expectedCount, len(awards))
	}

	// If Venus Next is enabled, Venuphile is mandatory
	if expansions["Venus Next"] {
		hasVenuphile := false
		for _, award := range awards {
			if award == "Venuphile" {
				hasVenuphile = true
				break
			}
		}
		if !hasVenuphile {
			return fmt.Errorf("Venuphile award is mandatory when Venus Next expansion is enabled")
		}
	}

	return nil
}
