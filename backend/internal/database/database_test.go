package database

import (
	"bytes"
	"strings"
	"testing"
	"terraforming-mars-backend/internal/models"
)

func setupTestDB(t *testing.T) (*Repository, *models.Player) {
	db, err := Init(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	
	repo := NewRepository(db)
	
	// Create a system admin for creating other players
	systemPassword := "systemadmin123"
	systemAdmin, err := repo.CreateSystemAdmin("system", &systemPassword)
	if err != nil {
		t.Fatalf("Failed to create system admin: %v", err)
	}
	
	return repo, systemAdmin
}

func TestCreateAndGetPlayer(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Test creating a player without password (role: player)
	player, err := repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	
	if player.Name != "Alice" {
		t.Errorf("Expected player name 'Alice', got '%s'", player.Name)
	}
	
	if player.ID == 0 {
		t.Error("Expected player ID to be non-zero")
	}
	
	if player.CreatedBy == nil || *player.CreatedBy != systemAdmin.ID {
		t.Errorf("Expected created_by %d, got %v", systemAdmin.ID, player.CreatedBy)
	}
	
	if player.PasswordHash != nil {
		t.Error("Expected password hash to be nil")
	}
	
	if player.Role != models.RolePlayer {
		t.Errorf("Expected role 'player', got '%s'", player.Role)
	}
	
	// Test retrieving player by name
	retrieved, err := repo.GetPlayerByName("Alice")
	if err != nil {
		t.Fatalf("Failed to retrieve player: %v", err)
	}
	
	if retrieved.ID != player.ID {
		t.Errorf("Expected player ID %d, got %d", player.ID, retrieved.ID)
	}
	
	if retrieved.Name != player.Name {
		t.Errorf("Expected player name '%s', got '%s'", player.Name, retrieved.Name)
	}
	
	if retrieved.Role != models.RolePlayer {
		t.Errorf("Expected role 'player', got '%s'", retrieved.Role)
	}
}

func TestGetAllPlayers(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create multiple players
	names := []string{"Alice", "Bob", "Charlie"}
	for _, name := range names {
		_, err := repo.CreatePlayer(name, nil, models.RolePlayer, *systemAdmin)
		if err != nil {
			t.Fatalf("Failed to create player %s: %v", name, err)
		}
	}
	
	// Get all players
	players, err := repo.GetAllPlayers()
	if err != nil {
		t.Fatalf("Failed to get all players: %v", err)
	}
	
	// Verify all test names are present (plus the system admin)
	foundNames := make(map[string]bool)
	for _, p := range players {
		foundNames[p.Name] = true
	}
	
	// Should find all test names plus the system admin
	expectedNames := append(names, "system")
	if len(players) != len(expectedNames) {
		t.Errorf("Expected %d players (including system), got %d", len(expectedNames), len(players))
	}
	
	for _, name := range names {
		if !foundNames[name] {
			t.Errorf("Player '%s' not found in results", name)
		}
	}
	
	// Verify system admin exists
	if !foundNames["system"] {
		t.Error("System admin not found in results")
	}
}

func TestCreateDuplicatePlayer(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create first player
	_, err := repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create first player: %v", err)
	}
	
	// Try to create duplicate
	_, err = repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err == nil {
		t.Error("Expected error when creating duplicate player, got nil")
	}
}

func TestCreateGame(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create players first
	players := []string{"Alice", "Bob", "Charlie"}
	for _, name := range players {
		_, err := repo.CreatePlayer(name, nil, models.RolePlayer, *systemAdmin)
		if err != nil {
			t.Fatalf("Failed to create player %s: %v", name, err)
		}
	}
	
	// Create a game with cleaner structure
	noteText := "Test note"
	req := models.CreateGameRequest{
		Name:        "Test Game",
		Date:        "2024-01-15",
		Map:         "Hellas",
		Generations: 12,
		Note:        &noteText,
		Expansions: models.Expansions{
			"base":     true,
			"prelude":  true,
			"colonies": false,
			"venus":    true,
			"turmoil":  false,
		},
		CreatedBy: systemAdmin.ID,
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 20, Cities: 5, Greeneries: 8, Cards: 15, TurmoilPoints: 0},
			{Name: "Bob", Corporation: "Helion", TerraformingRating: 22, Cities: 6, Greeneries: 4, Cards: 12, TurmoilPoints: 5},
			{Name: "Charlie", Corporation: "UNMI", TerraformingRating: 25, Cities: 3, Greeneries: 5, Cards: 18, TurmoilPoints: 10},
		},
		Milestones: []models.MilestoneRequest{
			{Name: "Terraformer", WinnerGamePlayerIndex: intPtr(0)},
			{Name: "Mayor", WinnerGamePlayerIndex: intPtr(1)},
			{Name: "Gardener", WinnerGamePlayerIndex: nil},
			{Name: "Builder", WinnerGamePlayerIndex: intPtr(2)},
			{Name: "Planner", WinnerGamePlayerIndex: intPtr(0)},
		},
		Awards: []models.AwardRequest{
			{
				Name: "Landlord",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 0, Placement: models.PlacementFirst}, // Alice gets first
					{PlayerIndex: 1, Placement: models.PlacementSecond}, // Bob gets second
				},
			},
			{
				Name: "Banker",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 2, Placement: models.PlacementFirst}, // Charlie gets first
					{PlayerIndex: 0, Placement: models.PlacementSecond}, // Alice gets second
				},
			},
			{
				Name: "Scientist",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 1, Placement: models.PlacementFirst}, // Bob gets first
				},
			},
		},
	}
	
	game, err := repo.CreateGame(req, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	
	// Verify game core details
	if game.Game.Name != "Test Game" {
		t.Errorf("Expected game name 'Test Game', got '%s'", game.Game.Name)
	}
	
	if game.Game.Map != "Hellas" {
		t.Errorf("Expected map 'Hellas', got '%s'", game.Game.Map)
	}
	
	if game.Game.Generations != 12 {
		t.Errorf("Expected 12 generations, got %d", game.Game.Generations)
	}
	
	if game.Game.Revision != 1 {
		t.Errorf("Expected revision 1, got %d", game.Game.Revision)
	}
	
	if game.Game.Note == nil || *game.Game.Note != "Test note" {
		t.Errorf("Expected note 'Test note', got %v", game.Game.Note)
	}
	
	// No need to check IsLatest anymore - we use MAX(revision) instead
	
	// Verify expansions
	if !game.Game.Expansions["base"] {
		t.Error("Expected base expansion to be true")
	}
	if !game.Game.Expansions["prelude"] {
		t.Error("Expected prelude expansion to be true")
	}
	if game.Game.Expansions["colonies"] {
		t.Error("Expected colonies expansion to be false")
	}
	if !game.Game.Expansions["venus"] {
		t.Error("Expected venus expansion to be true")
	}
	if game.Game.Expansions["turmoil"] {
		t.Error("Expected turmoil expansion to be false")
	}
	
	// Verify players count
	if len(game.GamePlayers) != 3 {
		t.Fatalf("Expected 3 game players, got %d", len(game.GamePlayers))
	}
	
	// Create a map for easier player lookup
	playerMap := make(map[string]models.GamePlayer)
	for _, gp := range game.GamePlayers {
		playerMap[gp.Corporation] = gp
	}
	
	// Verify Alice's data
	alice := playerMap["Ecoline"]
	if alice.TerraformingRating != 20 {
		t.Errorf("Alice: Expected TR 20, got %d", alice.TerraformingRating)
	}
	if alice.Cities != 5 {
		t.Errorf("Alice: Expected 5 cities, got %d", alice.Cities)
	}
	if alice.Greeneries != 8 {
		t.Errorf("Alice: Expected 8 greeneries, got %d", alice.Greeneries)
	}
	if alice.Cards != 15 {
		t.Errorf("Alice: Expected 15 cards, got %d", alice.Cards)
	}
	if alice.TurmoilPoints != 0 {
		t.Errorf("Alice: Expected 0 turmoil points, got %d", alice.TurmoilPoints)
	}
	if alice.MilestonePoints != 10 { // 2 milestones * 5 points
		t.Errorf("Alice: Expected 10 milestone points, got %d", alice.MilestonePoints)
	}
	if alice.AwardPoints != 7 { // 1 gold (5) + 1 silver (2)
		t.Errorf("Alice: Expected 7 award points, got %d", alice.AwardPoints)
	}
	expectedAliceTotal := 20 + 5 + 8 + 15 + 0 + 10 + 7 // = 65
	if alice.TotalPoints != expectedAliceTotal {
		t.Errorf("Alice: Expected %d total points, got %d", expectedAliceTotal, alice.TotalPoints)
	}
	
	// Verify Bob's data
	bob := playerMap["Helion"]
	if bob.TerraformingRating != 22 {
		t.Errorf("Bob: Expected TR 22, got %d", bob.TerraformingRating)
	}
	if bob.Cities != 6 {
		t.Errorf("Bob: Expected 6 cities, got %d", bob.Cities)
	}
	if bob.Greeneries != 4 {
		t.Errorf("Bob: Expected 4 greeneries, got %d", bob.Greeneries)
	}
	if bob.Cards != 12 {
		t.Errorf("Bob: Expected 12 cards, got %d", bob.Cards)
	}
	if bob.TurmoilPoints != 5 {
		t.Errorf("Bob: Expected 5 turmoil points, got %d", bob.TurmoilPoints)
	}
	if bob.MilestonePoints != 5 { // 1 milestone * 5 points
		t.Errorf("Bob: Expected 5 milestone points, got %d", bob.MilestonePoints)
	}
	if bob.AwardPoints != 7 { // 1 silver (2) + 1 gold (5)
		t.Errorf("Bob: Expected 7 award points, got %d", bob.AwardPoints)
	}
	expectedBobTotal := 22 + 6 + 4 + 12 + 5 + 5 + 7 // = 61
	if bob.TotalPoints != expectedBobTotal {
		t.Errorf("Bob: Expected %d total points, got %d", expectedBobTotal, bob.TotalPoints)
	}
	
	// Verify Charlie's data
	charlie := playerMap["UNMI"]
	if charlie.TerraformingRating != 25 {
		t.Errorf("Charlie: Expected TR 25, got %d", charlie.TerraformingRating)
	}
	if charlie.Cities != 3 {
		t.Errorf("Charlie: Expected 3 cities, got %d", charlie.Cities)
	}
	if charlie.Greeneries != 5 {
		t.Errorf("Charlie: Expected 5 greeneries, got %d", charlie.Greeneries)
	}
	if charlie.Cards != 18 {
		t.Errorf("Charlie: Expected 18 cards, got %d", charlie.Cards)
	}
	if charlie.TurmoilPoints != 10 {
		t.Errorf("Charlie: Expected 10 turmoil points, got %d", charlie.TurmoilPoints)
	}
	if charlie.MilestonePoints != 5 { // 1 milestone * 5 points
		t.Errorf("Charlie: Expected 5 milestone points, got %d", charlie.MilestonePoints)
	}
	if charlie.AwardPoints != 5 { // 1 gold (5)
		t.Errorf("Charlie: Expected 5 award points, got %d", charlie.AwardPoints)
	}
	expectedCharlieTotal := 25 + 3 + 5 + 18 + 10 + 5 + 5 // = 71
	if charlie.TotalPoints != expectedCharlieTotal {
		t.Errorf("Charlie: Expected %d total points, got %d", expectedCharlieTotal, charlie.TotalPoints)
	}
	
	// Verify milestones
	if len(game.Milestones) != 5 {
		t.Fatalf("Expected 5 milestones, got %d", len(game.Milestones))
	}
	
	// Check specific milestones
	milestoneMap := make(map[string]models.Milestone)
	for _, m := range game.Milestones {
		milestoneMap[m.Name] = m
	}
	
	if milestoneMap["Terraformer"].WinnerGamePlayerID == nil || *milestoneMap["Terraformer"].WinnerGamePlayerID != alice.ID {
		t.Error("Expected Terraformer milestone to be won by Alice")
	}
	if milestoneMap["Mayor"].WinnerGamePlayerID == nil || *milestoneMap["Mayor"].WinnerGamePlayerID != bob.ID {
		t.Error("Expected Mayor milestone to be won by Bob")
	}
	if milestoneMap["Gardener"].WinnerGamePlayerID != nil {
		t.Error("Expected Gardener milestone to have no winner")
	}
	if milestoneMap["Builder"].WinnerGamePlayerID == nil || *milestoneMap["Builder"].WinnerGamePlayerID != charlie.ID {
		t.Error("Expected Builder milestone to be won by Charlie")
	}
	if milestoneMap["Planner"].WinnerGamePlayerID == nil || *milestoneMap["Planner"].WinnerGamePlayerID != alice.ID {
		t.Error("Expected Planner milestone to be won by Alice")
	}
	
	// Verify awards
	if len(game.Awards) != 3 {
		t.Fatalf("Expected 3 awards, got %d", len(game.Awards))
	}
	
	// Check award placements
	if len(game.Placements) != 5 { // Total placements across all awards
		t.Errorf("Expected 5 total award placements, got %d", len(game.Placements))
	}
	
	// Verify placement details
	awardMap := make(map[string]int)
	for _, a := range game.Awards {
		awardMap[a.Name] = a.ID
	}
	
	placementsByAward := make(map[int][]models.AwardPlacement)
	for _, p := range game.Placements {
		placementsByAward[p.AwardID] = append(placementsByAward[p.AwardID], p)
	}
	
	// Check Landlord award placements
	landlordPlacements := placementsByAward[awardMap["Landlord"]]
	if len(landlordPlacements) != 2 {
		t.Errorf("Expected 2 placements for Landlord, got %d", len(landlordPlacements))
	}
	
	// Check Banker award placements
	bankerPlacements := placementsByAward[awardMap["Banker"]]
	if len(bankerPlacements) != 2 {
		t.Errorf("Expected 2 placements for Banker, got %d", len(bankerPlacements))
	}
	
	// Check Scientist award placements
	scientistPlacements := placementsByAward[awardMap["Scientist"]]
	if len(scientistPlacements) != 1 {
		t.Errorf("Expected 1 placement for Scientist, got %d", len(scientistPlacements))
	}
}

func TestGetGameByID(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create players
	_, err := repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	
	// Create a simple game
	req := models.CreateGameRequest{
		Name:        "Test Game",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 10,
		Expansions:  models.Expansions{"base": true},
		CreatedBy:   systemAdmin.ID,
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 20, Cities: 5, Greeneries: 8, Cards: 15, TurmoilPoints: 0},
		},
	}
	
	createdGame, err := repo.CreateGame(req, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	
	// Retrieve the game
	retrievedGame, err := repo.GetGameByID(createdGame.Game.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve game: %v", err)
	}
	
	if retrievedGame.Game.ID != createdGame.Game.ID {
		t.Errorf("Expected game ID %d, got %d", createdGame.Game.ID, retrievedGame.Game.ID)
	}
	
	if retrievedGame.Game.Name != createdGame.Game.Name {
		t.Errorf("Expected game name '%s', got '%s'", createdGame.Game.Name, retrievedGame.Game.Name)
	}
}

func TestCreateGameWithoutPlayers(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	req := models.CreateGameRequest{
		Name:        "Empty Game",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 10,
		Expansions:  models.Expansions{"base": true},
		CreatedBy:   systemAdmin.ID,
		Players:     []models.PlayerRequest{},
	}
	
	_, err := repo.CreateGame(req, *systemAdmin)
	if err == nil {
		t.Error("Expected error when creating game without players, got nil")
	}
}

func TestCreateGameWithNonExistentPlayer(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	req := models.CreateGameRequest{
		Name:        "Test Game",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 10,
		Expansions:  models.Expansions{"base": true},
		CreatedBy:   systemAdmin.ID,
		Players: []models.PlayerRequest{
			{Name: "NonExistent", Corporation: "Ecoline", TerraformingRating: 20, Cities: 5, Greeneries: 8, Cards: 15, TurmoilPoints: 0},
		},
	}
	
	_, err := repo.CreateGame(req, *systemAdmin)
	if err == nil {
		t.Error("Expected error when creating game with non-existent player, got nil")
	}
}

func TestUpdateGame(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create players
	players := []string{"Alice", "Bob"}
	for _, name := range players {
		_, err := repo.CreatePlayer(name, nil, models.RolePlayer, *systemAdmin)
		if err != nil {
			t.Fatalf("Failed to create player %s: %v", name, err)
		}
	}
	
	// Create initial game
	initialReq := models.CreateGameRequest{
		Name:        "Original Game",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 10,
		Expansions:  models.Expansions{"base": true},
		CreatedBy:   systemAdmin.ID,
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 20, Cities: 5, Greeneries: 8, Cards: 15, TurmoilPoints: 0},
			{Name: "Bob", Corporation: "Helion", TerraformingRating: 18, Cities: 4, Greeneries: 6, Cards: 12, TurmoilPoints: 3},
		},
		Milestones: []models.MilestoneRequest{
			{Name: "Terraformer", WinnerGamePlayerIndex: intPtr(0)},
		},
		Awards: []models.AwardRequest{
			{
				Name: "Landlord",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 0, Placement: models.PlacementFirst},
				},
			},
		},
	}
	
	createdGame, err := repo.CreateGame(initialReq, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create initial game: %v", err)
	}
	
	// Update the game with corrected scores
	updateReq := models.CreateGameRequest{
		Name:        "Original Game (Corrected)",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 11, // Changed
		Expansions:  models.Expansions{"base": true, "prelude": true}, // Added prelude
		CreatedBy:   systemAdmin.ID,
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 22, Cities: 6, Greeneries: 9, Cards: 16, TurmoilPoints: 2}, // Updated scores
			{Name: "Bob", Corporation: "Helion", TerraformingRating: 20, Cities: 5, Greeneries: 7, Cards: 14, TurmoilPoints: 5}, // Updated scores
		},
		Milestones: []models.MilestoneRequest{
			{Name: "Terraformer", WinnerGamePlayerIndex: intPtr(0)},
			{Name: "Mayor", WinnerGamePlayerIndex: intPtr(1)}, // Added milestone
		},
		Awards: []models.AwardRequest{
			{
				Name: "Landlord",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 1, Placement: models.PlacementFirst}, // Changed winner
					{PlayerIndex: 0, Placement: models.PlacementSecond}, // Alice now second
				},
			},
			{
				Name: "Banker", // Added award
				Placements: []models.PlacementRequest{
					{PlayerIndex: 0, Placement: models.PlacementFirst},
				},
			},
		},
	}
	
	updatedGame, err := repo.UpdateGame(createdGame.Game.ID, updateReq, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to update game: %v", err)
	}
	
	// Verify the updated game has the same GameID across revisions
	if updatedGame.Game.GameID != createdGame.Game.GameID {
		t.Error("Expected updated game to have the same GameID across revisions")
	}
	
	// Verify revision number increased
	if updatedGame.Game.Revision != 2 {
		t.Errorf("Expected revision 2, got %d", updatedGame.Game.Revision)
	}
	
	// Updated game should have revision 2
	if updatedGame.Game.Revision != 2 {
		t.Errorf("Expected updated game to have revision 2, got %d", updatedGame.Game.Revision)
	}
	
	// Verify updated game details
	if updatedGame.Game.Name != "Original Game (Corrected)" {
		t.Errorf("Expected name 'Original Game (Corrected)', got '%s'", updatedGame.Game.Name)
	}
	
	if updatedGame.Game.Generations != 11 {
		t.Errorf("Expected 11 generations, got %d", updatedGame.Game.Generations)
	}
	
	if !updatedGame.Game.Expansions["prelude"] {
		t.Error("Expected prelude expansion to be true")
	}
	
	// Verify player scores were updated
	if len(updatedGame.GamePlayers) != 2 {
		t.Fatalf("Expected 2 game players, got %d", len(updatedGame.GamePlayers))
	}
	
	// Create a map for easier player lookup
	playerMap := make(map[string]models.GamePlayer)
	for _, gp := range updatedGame.GamePlayers {
		playerMap[gp.Corporation] = gp
	}
	
	// Check Alice's updated scores
	alice := playerMap["Ecoline"]
	if alice.TerraformingRating != 22 {
		t.Errorf("Alice: Expected TR 22, got %d", alice.TerraformingRating)
	}
	if alice.Cities != 6 {
		t.Errorf("Alice: Expected 6 cities, got %d", alice.Cities)
	}
	if alice.Greeneries != 9 {
		t.Errorf("Alice: Expected 9 greeneries, got %d", alice.Greeneries)
	}
	if alice.Cards != 16 {
		t.Errorf("Alice: Expected 16 cards, got %d", alice.Cards)
	}
	if alice.TurmoilPoints != 2 {
		t.Errorf("Alice: Expected 2 turmoil points, got %d", alice.TurmoilPoints)
	}
	
	// Check Bob's updated scores
	bob := playerMap["Helion"]
	if bob.TerraformingRating != 20 {
		t.Errorf("Bob: Expected TR 20, got %d", bob.TerraformingRating)
	}
	if bob.Cities != 5 {
		t.Errorf("Bob: Expected 5 cities, got %d", bob.Cities)
	}
	if bob.Greeneries != 7 {
		t.Errorf("Bob: Expected 7 greeneries, got %d", bob.Greeneries)
	}
	if bob.Cards != 14 {
		t.Errorf("Bob: Expected 14 cards, got %d", bob.Cards)
	}
	if bob.TurmoilPoints != 5 {
		t.Errorf("Bob: Expected 5 turmoil points, got %d", bob.TurmoilPoints)
	}
	
	// Verify milestones were updated
	if len(updatedGame.Milestones) != 2 {
		t.Errorf("Expected 2 milestones, got %d", len(updatedGame.Milestones))
	}
	
	// Verify awards were updated
	if len(updatedGame.Awards) != 2 {
		t.Errorf("Expected 2 awards, got %d", len(updatedGame.Awards))
	}
	
	// GetGameByID always returns the latest revision, so originalGame should be the updated one
	latestGame, err := repo.GetGameByID(createdGame.Game.GameID)
	if err != nil {
		t.Fatalf("Failed to retrieve latest game: %v", err)
	}
	
	// The latest game should be the updated revision
	if latestGame.Game.Revision != 2 {
		t.Errorf("Expected latest game to have revision 2, got %d", latestGame.Game.Revision)
	}
	
	// Verify the latest revision has the updated data
	if latestGame.Game.Generations != 11 {
		t.Errorf("Expected latest game to have 11 generations, got %d", latestGame.Game.Generations)
	}
}

func TestPlayerPasswordAuth(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create player with password (user role)
	password := "password123"
	player, err := repo.CreatePlayer("Alice", &password, models.RoleUser, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player with password: %v", err)
	}
	
	if player.PasswordHash == nil {
		t.Error("Expected password hash to be set")
	}
	
	// Test successful authentication
	authenticatedPlayer, err := repo.AuthenticatePlayer("Alice", "password123")
	if err != nil {
		t.Fatalf("Failed to authenticate player: %v", err)
	}
	
	if authenticatedPlayer.ID != player.ID {
		t.Errorf("Expected authenticated player ID %d, got %d", player.ID, authenticatedPlayer.ID)
	}
	
	// Test failed authentication with wrong password
	_, err = repo.AuthenticatePlayer("Alice", "wrongpassword")
	if err == nil {
		t.Error("Expected error when authenticating with wrong password, got nil")
	}
	
	// Test authentication with player that has no password
	_, err = repo.CreatePlayer("Bob", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player without password: %v", err)
	}
	
	_, err = repo.AuthenticatePlayer("Bob", "anypassword")
	if err == nil {
		t.Error("Expected error when authenticating player with no password, got nil")
	}
}

func TestUpdatePlayer(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create player
	player, err := repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	
	// Update player with password (change to user role)
	newPassword := "newpassword123"
	userRole := models.RoleUser
	updatedPlayer, err := repo.UpdatePlayer(player.ID, "Alice Updated", &newPassword, &userRole, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to update player: %v", err)
	}
	
	if updatedPlayer.Name != "Alice Updated" {
		t.Errorf("Expected name 'Alice Updated', got '%s'", updatedPlayer.Name)
	}
	
	if updatedPlayer.PasswordHash == nil {
		t.Error("Expected password hash to be set after update")
	}
	
	// Test authentication with new password
	_, err = repo.AuthenticatePlayer("Alice Updated", "newpassword123")
	if err != nil {
		t.Fatalf("Failed to authenticate with updated password: %v", err)
	}
}

func TestPasswordValidation(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Test password too short
	shortPassword := "short"
	_, err := repo.CreatePlayer("Alice", &shortPassword, models.RoleUser, *systemAdmin)
	if err == nil {
		t.Error("Expected error when creating player with short password, got nil")
	}
	
	// Test valid password
	validPassword := "validpassword123"
	_, err = repo.CreatePlayer("Bob", &validPassword, models.RoleUser, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player with valid password: %v", err)
	}
}

func TestRoleValidation(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Test creating admin with password (should work)
	adminPassword := "adminpassword123"
	admin, err := repo.CreatePlayer("Admin", &adminPassword, models.RoleAdmin, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}
	if admin.Role != models.RoleAdmin {
		t.Errorf("Expected role 'admin', got '%s'", admin.Role)
	}
	
	// Test creating user with password (should work)
	userPassword := "userpassword123"
	user, err := repo.CreatePlayer("User", &userPassword, models.RoleUser, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if user.Role != models.RoleUser {
		t.Errorf("Expected role 'user', got '%s'", user.Role)
	}
	
	// Test creating player without password (should work)
	player, err := repo.CreatePlayer("Player", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	if player.Role != models.RolePlayer {
		t.Errorf("Expected role 'player', got '%s'", player.Role)
	}
	
	// Test creating admin without password (should fail)
	_, err = repo.CreatePlayer("BadAdmin", nil, models.RoleAdmin, *systemAdmin)
	if err == nil {
		t.Error("Expected error when creating admin without password, got nil")
	}
	
	// Test creating user without password (should fail)
	_, err = repo.CreatePlayer("BadUser", nil, models.RoleUser, *systemAdmin)
	if err == nil {
		t.Error("Expected error when creating user without password, got nil")
	}
	
	// Test creating player with password (should fail)
	playerPassword := "shouldnotwork123"
	_, err = repo.CreatePlayer("BadPlayer", &playerPassword, models.RolePlayer, *systemAdmin)
	if err == nil {
		t.Error("Expected error when creating player with password, got nil")
	}
}

func TestRoleTransition(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create a player
	player, err := repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	
	// Upgrade player to user (should require password)
	userRole := models.RoleUser
	userPassword := "newuserpassword123"
	updatedPlayer, err := repo.UpdatePlayer(player.ID, "Alice", &userPassword, &userRole, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to upgrade player to user: %v", err)
	}
	if updatedPlayer.Role != models.RoleUser {
		t.Errorf("Expected role 'user', got '%s'", updatedPlayer.Role)
	}
	
	// Downgrade user to player (should clear password)
	playerRole := models.RolePlayer
	downgradedPlayer, err := repo.UpdatePlayer(updatedPlayer.ID, "Alice", nil, &playerRole, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to downgrade user to player: %v", err)
	}
	if downgradedPlayer.Role != models.RolePlayer {
		t.Errorf("Expected role 'player', got '%s'", downgradedPlayer.Role)
	}
	if downgradedPlayer.PasswordHash != nil {
		t.Error("Expected password hash to be cleared when downgrading to player")
	}
	
	// Try to upgrade to user without password (should fail)
	_, err = repo.UpdatePlayer(downgradedPlayer.ID, "Alice", nil, &userRole, *systemAdmin)
	if err == nil {
		t.Error("Expected error when upgrading to user without password, got nil")
	}
}

func TestGameModificationPermissions(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create a regular user who will create a game
	userPassword := "userpassword123"
	user, err := repo.CreatePlayer("RegularUser", &userPassword, models.RoleUser, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Create another user who should NOT be able to modify the game
	otherUser, err := repo.CreatePlayer("OtherUser", &userPassword, models.RoleUser, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create other user: %v", err)
	}
	
	// Create players for the game
	_, err = repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	
	// Create a game as the regular user
	req := models.CreateGameRequest{
		Name:        "Test Game",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 10,
		Expansions:  models.Expansions{"base": true},
		CreatedBy:   user.ID, // This will be overridden by the actor
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 20, Cities: 5, Greeneries: 8, Cards: 15, TurmoilPoints: 0},
		},
	}
	
	createdGame, err := repo.CreateGame(req, *user)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	
	// Verify the game was created by the user
	if createdGame.Game.CreatedBy != user.ID {
		t.Errorf("Expected game created by user %d, got %d", user.ID, createdGame.Game.CreatedBy)
	}
	
	// Test 1: User should be able to modify their own game
	updateReq := models.CreateGameRequest{
		Name:        "Updated Game",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 11,
		Expansions:  models.Expansions{"base": true},
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 22, Cities: 6, Greeneries: 9, Cards: 16, TurmoilPoints: 2},
		},
	}
	
	_, err = repo.UpdateGame(createdGame.Game.ID, updateReq, *user)
	if err != nil {
		t.Errorf("User should be able to update their own game, got error: %v", err)
	}
	
	// Test 2: Other user should NOT be able to modify the game
	_, err = repo.UpdateGame(createdGame.Game.ID, updateReq, *otherUser)
	if err == nil {
		t.Error("Other user should NOT be able to update someone else's game")
	}
	
	// Test 3: Admin should be able to modify any game
	_, err = repo.UpdateGame(createdGame.Game.ID, updateReq, *systemAdmin)
	if err != nil {
		t.Errorf("Admin should be able to update any game, got error: %v", err)
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

func TestGameImages(t *testing.T) {
	repo, systemAdmin := setupTestDB(t)
	
	// Create players
	_, err := repo.CreatePlayer("Alice", nil, models.RolePlayer, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	
	// Create sample image bytes
	image1 := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A}
	image2 := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	
	req := models.CreateGameRequest{
		Name:        "Game with Images",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 10,
		Expansions:  models.Expansions{"base": true},
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 20},
		},
		Images: []models.ImageRequest{
			{ImageData: image1, MimeType: "image/png"},
			{ImageData: image2, MimeType: "image/jpeg"},
		},
	}
	
	// Create game
	game, err := repo.CreateGame(req, *systemAdmin)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	
	// Verify image metadata
	if len(game.Images) != 2 {
		t.Fatalf("Expected 2 images, got %d", len(game.Images))
	}
	
	if game.Images[0].MimeType != "image/png" {
		t.Errorf("Expected first image mime type 'image/png', got '%s'", game.Images[0].MimeType)
	}
	
	if game.Images[1].MimeType != "image/jpeg" {
		t.Errorf("Expected second image mime type 'image/jpeg', got '%s'", game.Images[1].MimeType)
	}
	
	// Test fetching image data
	imageData, mimeType, err := repo.GetGameImageData(game.Images[0].ID)
	if err != nil {
		t.Fatalf("Failed to get image data: %v", err)
	}
	
	if !bytes.Equal(imageData, image1) {
		t.Error("Retrieved image data doesn't match")
	}
	
	if mimeType != "image/png" {
		t.Errorf("Expected mime type 'image/png', got '%s'", mimeType)
	}
	
	// Test max images limit
	tooManyImages := make([]models.ImageRequest, 6)
	for i := 0; i < 6; i++ {
		tooManyImages[i] = models.ImageRequest{ImageData: image1, MimeType: "image/png"}
	}
	
	req.Images = tooManyImages
	_, err = repo.CreateGame(req, *systemAdmin)
	if err == nil {
		t.Error("Expected error with 6 images")
	}
	if !strings.Contains(err.Error(), "too many images") {
		t.Errorf("Expected 'too many images' error, got: %v", err)
	}
}