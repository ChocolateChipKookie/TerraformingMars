package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func resolvePlayerName(name string) string {
	if !strings.Contains(name, "/") {
		return name
	}

	parts := strings.Split(name, "/")
	first := parts[0]

	if first == "Jan" {
		return "Jan Mastrović"
	} else if first == "Mia" {
		return "Mia Čaušević"
	}

	return name
}

// Legacy game structures (old format without detailed milestones/awards)
type LegacyGameScore struct {
	Player      string `json:"player"`
	Corporation string `json:"corporation"`
	TR          int    `json:"tr"`
	Awards      int    `json:"awards"`
	Milestones  int    `json:"milestones"`
	Greenery    int    `json:"greenery"`
	Cities      int    `json:"cities"`
	Cards       int    `json:"cards"`
	Lead        int    `json:"lead"`
	Total       int    `json:"total"`
	Rank        int    `json:"rank"`
	Note        string `json:"note"`
	Turmoil     *int   `json:"turmoil,omitempty"`
}

type LegacyGame struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Note       string            `json:"note"`
	Date       string            `json:"date"`
	Mode       []string          `json:"mode"`
	ModeCode   string            `json:"mode_code"`
	Map        string            `json:"map"`
	Winner     string            `json:"winner"`
	WinCorp    string            `json:"win_corp"`
	WinScore   int               `json:"win_score"`
	Players    int               `json:"players"`
	Generation int               `json:"generation"`
	Image      string            `json:"image"`
	Scores     []LegacyGameScore `json:"scores"`
}

// Modern game structures (new format with detailed milestones/awards)
type ModernGameScore struct {
	Player      string `json:"player"`
	Corporation string `json:"corporation"`
	TR          int    `json:"tr"`
	Awards      int    `json:"awards"`
	Milestones  int    `json:"milestones"`
	Greenery    int    `json:"greenery"`
	Cities      int    `json:"cities"`
	Cards       int    `json:"cards"`
	Lead        int    `json:"lead"`
	Total       int    `json:"total"`
	Rank        int    `json:"rank"`
	Note        string `json:"note"`
	Turmoil     *int   `json:"turmoil,omitempty"`
}

type AwardWinners struct {
	First  []string `json:"first"`
	Second []string `json:"second"`
}

type ModernGameAwards struct {
	Winners map[string]AwardWinners   `json:"winners"`
	Points  map[string]map[string]int `json:"points"`
}

type Colony struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ModernGame struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Note       string            `json:"note"`
	Date       string            `json:"date"`
	Mode       []string          `json:"mode"`
	ModeCode   string            `json:"mode_code"`
	Map        string            `json:"map"`
	Winner     string            `json:"winner"`
	WinCorp    string            `json:"win_corp"`
	WinScore   int               `json:"win_score"`
	Players    int               `json:"players"`
	Generation int               `json:"generation"`
	Image      string            `json:"image"`
	Scores     []ModernGameScore `json:"scores"`
	Milestones map[string]string `json:"milestones"`
	Awards     ModernGameAwards  `json:"awards"`
	Colonies   []Colony          `json:"colonies,omitempty"`
}

// API request structures
type CreatePlayerRequest struct {
	Name          string  `json:"name"`
	Password      *string `json:"password,omitempty"`
	Role          string  `json:"role"`
	ActorName     string  `json:"actor_name"`
	ActorPassword string  `json:"actor_password"`
}

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type CreateGameResponse struct {
	ID int `json:"id"`
}

type APIClient struct {
	baseURL       string
	actorName     string
	actorPassword string
	client        *http.Client
}

func NewAPIClient(baseURL, actorName, actorPassword string) *APIClient {
	return &APIClient{
		baseURL:       baseURL,
		actorName:     actorName,
		actorPassword: actorPassword,
		client:        &http.Client{},
	}
}

func (c *APIClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *APIClient) GetPlayers() ([]Player, error) {
	resp, err := c.doRequest("GET", "/api/players", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var players []Player
	if err := json.NewDecoder(resp.Body).Decode(&players); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return players, nil
}

func (c *APIClient) CreatePlayer(name string) (*Player, error) {
	req := CreatePlayerRequest{
		Name:          name,
		Role:          "player",
		ActorName:     c.actorName,
		ActorPassword: c.actorPassword,
	}

	resp, err := c.doRequest("POST", "/api/players", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var player Player
	if err := json.Unmarshal(body, &player); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &player, nil
}

func (c *APIClient) CreateGame(gameReq interface{}) (*CreateGameResponse, error) {
	resp, err := c.doRequest("POST", "/api/games", gameReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var gameResp CreateGameResponse
	if err := json.Unmarshal(body, &gameResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &gameResp, nil
}

type GameEntry struct {
	Legacy *LegacyGame
	Modern *ModernGame
}

// API game request structures
type APILegacyPlayerRequest struct {
	Name               string `json:"name"`
	TerraformingRating int    `json:"terraforming_rating"`
	Cities             int    `json:"cities"`
	Greeneries         int    `json:"greeneries"`
	Cards              int    `json:"cards"`
	TurmoilPoints      int    `json:"turmoil_points"`
	MilestonePoints    int    `json:"milestone_points"`
	AwardPoints        int    `json:"award_points"`
}

type APILegacyGameRequest struct {
	LegacyMode    bool                     `json:"legacy_mode"`
	Name          string                   `json:"name"`
	Date          string                   `json:"date"`
	Note          *string                  `json:"note,omitempty"`
	Players       []APILegacyPlayerRequest `json:"players"`
	Images        []interface{}            `json:"images"`
	ActorName     string                   `json:"actor_name"`
	ActorPassword string                   `json:"actor_password"`
}

type APIPlayerRequest struct {
	Name               string `json:"name"`
	Corporation        string `json:"corporation"`
	TerraformingRating int    `json:"terraforming_rating"`
	Cities             int    `json:"cities"`
	Greeneries         int    `json:"greeneries"`
	Cards              int    `json:"cards"`
	TurmoilPoints      int    `json:"turmoil_points"`
}

type APIMilestoneRequest struct {
	Name                  string `json:"name"`
	WinnerGamePlayerIndex *int   `json:"winner_game_player_index"`
}

type APIPlacementRequest struct {
	PlayerIndex int `json:"player_index"`
	Placement   int `json:"placement"`
}

type APIAwardRequest struct {
	Name       string                `json:"name"`
	Placements []APIPlacementRequest `json:"placements"`
}

type APIGameRequest struct {
	LegacyMode    bool                   `json:"legacy_mode"`
	Name          string                 `json:"name"`
	Date          string                 `json:"date"`
	Map           string                 `json:"map"`
	Generations   int                    `json:"generations"`
	Expansions    map[string]bool        `json:"expansions"`
	Note          *string                `json:"note,omitempty"`
	Players       []APIPlayerRequest     `json:"players"`
	Milestones    []APIMilestoneRequest  `json:"milestones"`
	Awards        []APIAwardRequest      `json:"awards"`
	Images        []interface{}          `json:"images"`
	ActorName     string                 `json:"actor_name"`
	ActorPassword string                 `json:"actor_password"`
}

func convertLegacyGame(game *LegacyGame, actorName, actorPassword string, gameIndex int) *APILegacyGameRequest {
	players := make([]APILegacyPlayerRequest, len(game.Scores))
	for i, score := range game.Scores {
		turmoilPoints := 0
		if score.Turmoil != nil {
			turmoilPoints = *score.Turmoil
		}

		players[i] = APILegacyPlayerRequest{
			Name:               score.Player,
			TerraformingRating: score.TR,
			Cities:             score.Cities,
			Greeneries:         score.Greenery,
			Cards:              score.Cards,
			TurmoilPoints:      turmoilPoints,
			MilestonePoints:    score.Milestones,
			AwardPoints:        score.Awards,
		}
	}

	var note *string
	if game.Note != "" {
		note = &game.Note
	}

	name := game.Name
	if name == "" {
		name = fmt.Sprintf("Game %d", gameIndex)
	}

	return &APILegacyGameRequest{
		LegacyMode:    true,
		Name:          name,
		Date:          game.Date,
		Note:          note,
		Players:       players,
		Images:        []interface{}{},
		ActorName:     actorName,
		ActorPassword: actorPassword,
	}
}

type MapData struct {
	Milestones []string
	Awards     []string
}

var mapDefaultsData = map[string]MapData{
	"Tharsis": {
		Milestones: []string{"Builder", "Gardener", "Mayor", "Planner", "Terraformer"},
		Awards:     []string{"Landlord", "Banker", "Scientist", "Thermalist", "Miner"},
	},
	"Hellas": {
		Milestones: []string{"Diversifier", "Tactician", "Polar Explorer", "Energizer", "Rim Settler"},
		Awards:     []string{"Cultivator", "Magnate", "Space Baron", "Excentric", "Contractor"},
	},
	"Elysium": {
		Milestones: []string{"Generalist", "Specialist", "Ecologist", "Tycoon", "Legend"},
		Awards:     []string{"Celebrity", "Industrialist", "Desert Settler", "Estate Dealer", "Benefactor"},
	},
}

var corporationNameMap = map[string]string{
	"Morning Star INC": "Morning Star Inc.",
	"UNMI":             "United Nations Mars Initiative",
	"Splice":           "Splice Tactical Genomics",
}

func normalizeCorporationName(name string) string {
	if mapped, ok := corporationNameMap[name]; ok {
		return mapped
	}
	return name
}

func modeToExpansions(mode []string) map[string]bool {
	expansions := make(map[string]bool)
	expansionMap := map[string]string{
		"Base":  "Base Game",
		"Venus": "Venus Next",
	}

	for _, m := range mode {
		if mapped, ok := expansionMap[m]; ok {
			expansions[mapped] = true
		} else {
			expansions[m] = true
		}
	}
	return expansions
}

func convertModernGame(game *ModernGame, actorName, actorPassword string, gameIndex int) *APIGameRequest {
	players := make([]APIPlayerRequest, len(game.Scores))
	playerNameToIndex := make(map[string]int)

	for i, score := range game.Scores {
		playerNameToIndex[score.Player] = i

		turmoilPoints := 0
		if score.Turmoil != nil {
			turmoilPoints = *score.Turmoil
		}

		players[i] = APIPlayerRequest{
			Name:               score.Player,
			Corporation:        normalizeCorporationName(score.Corporation),
			TerraformingRating: score.TR,
			Cities:             score.Cities,
			Greeneries:         score.Greenery,
			Cards:              score.Cards,
			TurmoilPoints:      turmoilPoints,
		}
	}

	expansions := modeToExpansions(game.Mode)
	hasVenus := expansions["Venus Next"]

	mapData, hasMapData := mapDefaultsData[game.Map]
	if !hasMapData {
		mapData = mapDefaultsData["Tharsis"]
	}

	allMilestones := make([]string, len(mapData.Milestones))
	copy(allMilestones, mapData.Milestones)
	if hasVenus {
		allMilestones = append(allMilestones, "Hoverlord")
	}

	milestones := make([]APIMilestoneRequest, len(allMilestones))
	for i, name := range allMilestones {
		var winnerIndex *int
		if winner, claimed := game.Milestones[name]; claimed {
			if idx, ok := playerNameToIndex[winner]; ok {
				winnerIndex = &idx
			}
		}
		milestones[i] = APIMilestoneRequest{
			Name:                  name,
			WinnerGamePlayerIndex: winnerIndex,
		}
	}

	allAwards := make([]string, len(mapData.Awards))
	copy(allAwards, mapData.Awards)
	if hasVenus {
		allAwards = append(allAwards, "Venuphile")
	}

	awards := make([]APIAwardRequest, len(allAwards))
	for i, awardName := range allAwards {
		placements := make([]APIPlacementRequest, 0)

		if winners, hasWinners := game.Awards.Winners[awardName]; hasWinners {
			for _, firstPlace := range winners.First {
				if idx, ok := playerNameToIndex[firstPlace]; ok {
					placements = append(placements, APIPlacementRequest{
						PlayerIndex: idx,
						Placement:   1,
					})
				}
			}

			for _, secondPlace := range winners.Second {
				if idx, ok := playerNameToIndex[secondPlace]; ok {
					placements = append(placements, APIPlacementRequest{
						PlayerIndex: idx,
						Placement:   2,
					})
				}
			}
		}

		awards[i] = APIAwardRequest{
			Name:       awardName,
			Placements: placements,
		}
	}

	var note *string
	if game.Note != "" {
		note = &game.Note
	}

	generations := game.Generation
	if generations == 0 {
		generations = 1
	}

	name := game.Name
	if name == "" {
		name = fmt.Sprintf("Game %d", gameIndex)
	}

	return &APIGameRequest{
		LegacyMode:    false,
		Name:          name,
		Date:          game.Date,
		Map:           game.Map,
		Generations:   generations,
		Expansions:    expansions,
		Note:          note,
		Players:       players,
		Milestones:    milestones,
		Awards:        awards,
		Images:        []interface{}{},
		ActorName:     actorName,
		ActorPassword: actorPassword,
	}
}

func loadGames(filepath string) ([]GameEntry, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var rawGames []json.RawMessage
	if err := json.Unmarshal(data, &rawGames); err != nil {
		return nil, fmt.Errorf("failed to parse JSON array: %w", err)
	}

	games := make([]GameEntry, 0, len(rawGames))
	for i, raw := range rawGames {
		var detector struct {
			Milestones interface{} `json:"milestones"`
		}
		if err := json.Unmarshal(raw, &detector); err != nil {
			return nil, fmt.Errorf("failed to detect game type at index %d: %w", i, err)
		}

		entry := GameEntry{}
		if detector.Milestones == nil {
			var legacy LegacyGame
			if err := json.Unmarshal(raw, &legacy); err != nil {
				return nil, fmt.Errorf("failed to parse legacy game at index %d: %w", i, err)
			}
			for j := range legacy.Scores {
				legacy.Scores[j].Player = resolvePlayerName(legacy.Scores[j].Player)
			}
			entry.Legacy = &legacy
		} else {
			var modern ModernGame
			if err := json.Unmarshal(raw, &modern); err != nil {
				return nil, fmt.Errorf("failed to parse modern game at index %d: %w", i, err)
			}
			for j := range modern.Scores {
				modern.Scores[j].Player = resolvePlayerName(modern.Scores[j].Player)
			}
			entry.Modern = &modern
		}
		games = append(games, entry)
	}

	return games, nil
}

func extractPlayerNames(games []GameEntry) []string {
	nameSet := make(map[string]bool)
	for _, game := range games {
		if game.Legacy != nil {
			for _, score := range game.Legacy.Scores {
				nameSet[score.Player] = true
			}
		} else if game.Modern != nil {
			for _, score := range game.Modern.Scores {
				nameSet[score.Player] = true
			}
		}
	}

	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	return names
}

func ensurePlayers(client *APIClient, playerNames []string) error {
	existingPlayers, err := client.GetPlayers()
	if err != nil {
		return fmt.Errorf("failed to get existing players: %w", err)
	}

	existingPlayerMap := make(map[string]bool)
	for _, player := range existingPlayers {
		existingPlayerMap[player.Name] = true
	}

	created := 0
	alreadyExisting := 0

	for _, name := range playerNames {
		if existingPlayerMap[name] {
			alreadyExisting++
			continue
		}

		_, err := client.CreatePlayer(name)
		if err != nil {
			return fmt.Errorf("failed to create player %s: %w", name, err)
		}
		created++
	}

	if created > 0 {
		fmt.Printf("Created %d players\n", created)
	}
	if alreadyExisting > 0 {
		fmt.Printf("%d players already existed\n", alreadyExisting)
	}
	fmt.Println("All players successfully ensured")

	return nil
}

func main() {
	actorName := os.Getenv("ACTOR_NAME")
	actorPassword := os.Getenv("ACTOR_PASSWORD")
	apiURL := os.Getenv("API_BASE_URL")
	dataDir := os.Getenv("DATA_DIR")

	if actorName == "" || actorPassword == "" || apiURL == "" || dataDir == "" {
		fmt.Println("Error: missing required environment variables")
		fmt.Println("Required: ACTOR_NAME, ACTOR_PASSWORD, API_BASE_URL, DATA_DIR")
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Terraforming Mars Data Migration")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("API Base URL: %s\n", apiURL)
	fmt.Printf("Actor: %s\n", actorName)
	fmt.Printf("Data Directory: %s\n", dataDir)
	fmt.Println(strings.Repeat("=", 60))

	client := NewAPIClient(apiURL, actorName, actorPassword)

	fmt.Println("\nLoading games from file...")
	games, err := loadGames(dataDir + "/games.json")
	if err != nil {
		fmt.Printf("Failed to load games: %v\n", err)
		os.Exit(1)
	}

	legacyCount := 0
	modernCount := 0
	for _, game := range games {
		if game.Legacy != nil {
			legacyCount++
		} else {
			modernCount++
		}
	}
	fmt.Printf("Loaded %d games (%d legacy, %d modern)\n", len(games), legacyCount, modernCount)

	playerNames := extractPlayerNames(games)
	fmt.Printf("Extracted %d unique player names\n", len(playerNames))

	if err := ensurePlayers(client, playerNames); err != nil {
		fmt.Printf("Failed to ensure players: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nInserting games...")
	createdGames := 0
	failedGames := 0

	for i, gameEntry := range games {
		var err error
		var gameName string

		if gameEntry.Legacy != nil {
			gameName = gameEntry.Legacy.Name
			if gameName == "" {
				gameName = fmt.Sprintf("Game %d", i+1)
			}
			req := convertLegacyGame(gameEntry.Legacy, actorName, actorPassword, i+1)
			_, err = client.CreateGame(req)
		} else if gameEntry.Modern != nil {
			gameName = gameEntry.Modern.Name
			if gameName == "" {
				gameName = fmt.Sprintf("Game %d", i+1)
			}
			req := convertModernGame(gameEntry.Modern, actorName, actorPassword, i+1)
			_, err = client.CreateGame(req)

			if err != nil && strings.Contains(err.Error(), "is not available with selected expansions") {
				fmt.Printf("[%d/%d] Retrying '%s' with Promo expansion enabled\n", i+1, len(games), gameName)
				req.Expansions["Promo"] = true
				_, err = client.CreateGame(req)
			}
		}

		if err != nil {
			fmt.Printf("[%d/%d] Failed to create game '%s': %v\n", i+1, len(games), gameName, err)
			failedGames++
		} else {
			fmt.Printf("[%d/%d] Inserted: %s\n", i+1, len(games), gameName)
			createdGames++
		}
	}

	fmt.Printf("\nGame insertion completed: %d created, %d failed\n", createdGames, failedGames)
}
