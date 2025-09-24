package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"terraforming-mars-backend/internal/models"
)

func TestGetGames(t *testing.T) {
	ctx := setupGameFixture(t)

	// Test GET /games
	req, _ := http.NewRequest("GET", "/games", nil)
	rr := httptest.NewRecorder()
	ctx.Handler.getGames(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	var games []models.Game
	err := json.NewDecoder(rr.Body).Decode(&games)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(games) != 2 {
		t.Errorf("Expected 2 games, got %d", len(games))
	}

	// Verify game details
	gamesByName := make(map[string]models.Game)
	for _, g := range games {
		gamesByName[g.Name] = g
	}

	// Check Admin Game
	if adminGame, ok := gamesByName["Admin Game"]; !ok {
		t.Error("Admin Game not found")
	} else {
		if adminGame.Map == nil || *adminGame.Map != "Tharsis" {
			t.Errorf("Expected map 'Tharsis', got '%v'", adminGame.Map)
		}
		if adminGame.Generations == nil || *adminGame.Generations != 10 {
			t.Errorf("Expected 10 generations, got %v", adminGame.Generations)
		}
		if adminGame.CreatedBy != ctx.Admin.ID {
			t.Errorf("Expected created_by %d, got %d", ctx.Admin.ID, adminGame.CreatedBy)
		}
		// Verify it matches our fixture game
		if adminGame.GameID != ctx.Game1.Game.GameID {
			t.Errorf("Expected game ID %d, got %d", ctx.Game1.Game.GameID, adminGame.GameID)
		}
	}

	// Check User Game
	if userGame, ok := gamesByName["User Game With Multiple Users"]; !ok {
		t.Error("User Game With Multiple Users not found")
	} else {
		if userGame.Map == nil || *userGame.Map != "Hellas" {
			t.Errorf("Expected map 'Hellas', got '%v'", userGame.Map)
		}
		if userGame.Generations == nil || *userGame.Generations != 12 {
			t.Errorf("Expected 12 generations, got %v", userGame.Generations)
		}
		if userGame.CreatedBy != ctx.Bob.ID {
			t.Errorf("Expected created_by %d (Bob), got %d", ctx.Bob.ID, userGame.CreatedBy)
		}
		// Verify it matches our fixture game
		if userGame.GameID != ctx.Game2.Game.GameID {
			t.Errorf("Expected game ID %d, got %d", ctx.Game2.Game.GameID, userGame.GameID)
		}
	}

	// Test milestone and award points calculation by getting detailed games
	// Get Game1 details to verify milestone/award points
	req1, _ := http.NewRequest("GET", "/games/"+fmt.Sprintf("%d", ctx.Game1.Game.GameID), nil)
	req1 = mux.SetURLVars(req1, map[string]string{"id": fmt.Sprintf("%d", ctx.Game1.Game.GameID)})
	rr1 := httptest.NewRecorder()
	ctx.Handler.getGame(rr1, req1)

	var game1Details models.GameWithDetails
	json.NewDecoder(rr1.Body).Decode(&game1Details)

	// Verify milestone points (5 points each)
	// Alice won Terraformer, Charlie won Mayor, Bob won Gardener
	aliceMilestonePoints := 0
	bobMilestonePoints := 0
	charlieMilestonePoints := 0

	for _, player := range game1Details.GamePlayers {
		switch player.PlayerID {
		case ctx.Alice.ID:
			aliceMilestonePoints = player.MilestonePoints
		case ctx.Bob.ID:
			bobMilestonePoints = player.MilestonePoints
		case ctx.Charlie.ID:
			charlieMilestonePoints = player.MilestonePoints
		}
	}

	if aliceMilestonePoints != 5 {
		t.Errorf("Alice should have 5 milestone points (Terraformer), got %d", aliceMilestonePoints)
	}
	if bobMilestonePoints != 5 {
		t.Errorf("Bob should have 5 milestone points (Gardener), got %d", bobMilestonePoints)
	}
	if charlieMilestonePoints != 5 {
		t.Errorf("Charlie should have 5 milestone points (Mayor), got %d", charlieMilestonePoints)
	}

	// Verify award points (5 for 1st place, 2 for 2nd place)
	// Landlord: Charlie 1st (5), Alice 2nd (2)
	// Scientist: Charlie 1st (5), Bob 2nd (2)
	// Thermalist: Bob 1st (5)
	aliceAwardPoints := 0
	bobAwardPoints := 0
	charlieAwardPoints := 0

	for _, player := range game1Details.GamePlayers {
		switch player.PlayerID {
		case ctx.Alice.ID:
			aliceAwardPoints = player.AwardPoints
		case ctx.Bob.ID:
			bobAwardPoints = player.AwardPoints
		case ctx.Charlie.ID:
			charlieAwardPoints = player.AwardPoints
		}
	}

	if aliceAwardPoints != 2 {
		t.Errorf("Alice should have 2 award points (Landlord 2nd), got %d", aliceAwardPoints)
	}
	if bobAwardPoints != 7 { // Scientist 2nd (2) + Thermalist 1st (5)
		t.Errorf("Bob should have 7 award points (Scientist 2nd + Thermalist 1st), got %d", bobAwardPoints)
	}
	if charlieAwardPoints != 10 { // Landlord 1st (5) + Scientist 1st (5)
		t.Errorf("Charlie should have 10 award points (Landlord 1st + Scientist 1st), got %d", charlieAwardPoints)
	}

	// Verify total points calculation
	for _, player := range game1Details.GamePlayers {
		expectedTotal := player.TerraformingRating + player.Cities + player.Greeneries +
			player.Cards + player.TurmoilPoints + player.MilestonePoints + player.AwardPoints
		if player.TotalPoints != expectedTotal {
			t.Errorf("Player %d total points mismatch: expected %d, got %d",
				player.PlayerID, expectedTotal, player.TotalPoints)
		}
	}
}

func TestCreateGame(t *testing.T) {
	ctx := setupGameFixture(t)
	
	// Create a reusable game request
	baseRequest := models.CreateGameRequest{
		Name:        "Test New Game",
		Date:        "2024-01-25",
		Map:         models.Ptr("Tharsis"),
		Generations: models.Ptr(8),
		Expansions:  models.Ptr(models.Expansions{"Prelude": true}),
		Players: []models.PlayerRequest{
			{
				Name:               "Bob",
				Corporation:        "Helion",
				TerraformingRating: 20,
				Cities:             1,
				Greeneries:         2,
				Cards:              8,
				TurmoilPoints:      3,
			},
		},
		Milestones: []models.MilestoneRequest{
			{
				Name:                  "Builder",
				WinnerGamePlayerIndex: models.Ptr(0),
			},
		},
		Awards: []models.AwardRequest{
			{
				Name: "Contractor",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 0, Placement: models.PlacementFirst},
				},
			},
		},
	}
	
	// Test user cannot create game with wrong password
	t.Run("User cannot create game with wrong password", func(t *testing.T) {
		req := AuthenticatedGameRequest{
			CreateGameRequest: baseRequest,
			ActorName:         "Bob",
			ActorPassword:     "wrongpassword",
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createGame(rr, httpReq)
		
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test player cannot create game
	t.Run("Player cannot create game", func(t *testing.T) {
		req := AuthenticatedGameRequest{
			CreateGameRequest: baseRequest,
			ActorName:         "Alice",
			ActorPassword:     "", // Players don't have passwords
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createGame(rr, httpReq)
		
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test user can create game with correct password
	t.Run("User can create game with correct password", func(t *testing.T) {
		req := AuthenticatedGameRequest{
			CreateGameRequest: baseRequest,
			ActorName:         "Bob",
			ActorPassword:     ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createGame(rr, httpReq)
		
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
		}
		
		var game models.GameWithDetails
		err := json.NewDecoder(rr.Body).Decode(&game)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		
		if game.Game.Name != "Test New Game" {
			t.Errorf("Expected name 'Test New Game', got '%s'", game.Game.Name)
		}
		if game.Game.CreatedBy != ctx.Bob.ID {
			t.Errorf("Expected created_by %d (Bob), got %d", ctx.Bob.ID, game.Game.CreatedBy)
		}
		if game.Game.Map == nil || *game.Game.Map != "Tharsis" {
			t.Errorf("Expected map 'Tharsis', got '%v'", game.Game.Map)
		}
		if len(game.GamePlayers) != 1 {
			t.Errorf("Expected 1 player, got %d", len(game.GamePlayers))
		}
		if len(game.Milestones) != 1 {
			t.Errorf("Expected 1 milestone, got %d", len(game.Milestones))
		}
		if len(game.Awards) != 1 {
			t.Errorf("Expected 1 award, got %d", len(game.Awards))
		}
	})
}

func TestGameImages(t *testing.T) {
	ctx := setupGameFixture(t)
	
	// Create sample image data
	image1 := createTestImage("image/png")
	image2 := createTestImage("image/jpeg")
	
	// Create a game with images
	gameReq := models.CreateGameRequest{
		Name:        "Game with Images",
		Date:        "2024-01-16",
		Map:         models.Ptr("Tharsis"),
		Generations: models.Ptr(10),
		Expansions:  models.Ptr(models.Expansions{"base": true}),
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 20},
		},
		Images: []models.ImageRequest{
			{ImageData: image1, MimeType: "image/png"},
			{ImageData: image2, MimeType: "image/jpeg"},
		},
	}
	
	req := AuthenticatedGameRequest{
		CreateGameRequest: gameReq,
		ActorName:         ctx.AdminName,
		ActorPassword:     ctx.AdminPassword,
	}
	
	// Create the game
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	ctx.Handler.createGame(rr, httpReq)
	
	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create game, status: %d, body: %s", rr.Code, rr.Body.String())
	}
	
	var game models.GameWithDetails
	json.NewDecoder(rr.Body).Decode(&game)
	
	// Verify image metadata is included in game response
	if len(game.Images) != 2 {
		t.Errorf("Expected 2 images in game response, got %d", len(game.Images))
	}
	
	// Both images should be converted to WebP
	if game.Images[0].MimeType != "image/webp" {
		t.Errorf("Expected first image mime type 'image/webp', got '%s'", game.Images[0].MimeType)
	}
	
	if game.Images[1].MimeType != "image/webp" {
		t.Errorf("Expected second image mime type 'image/webp', got '%s'", game.Images[1].MimeType)
	}
	
	// Test GET /images/{id} for first image
	t.Run("Get first image data", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/images/%d", game.Images[0].ID), nil)
		rr := httptest.NewRecorder()
		
		router := mux.NewRouter()
		router.HandleFunc("/images/{id}", ctx.Handler.getImage).Methods("GET")
		router.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		
		if rr.Header().Get("Content-Type") != "image/webp" {
			t.Errorf("Expected Content-Type 'image/webp', got '%s'", rr.Header().Get("Content-Type"))
		}
		
		if rr.Header().Get("Cache-Control") != "public, max-age=3600" {
			t.Errorf("Expected Cache-Control header")
		}
		
		if len(rr.Body.Bytes()) == 0 {
			t.Error("Response body is empty")
		}
	})
	
	// Test GET /images/{id} for second image
	t.Run("Get second image data", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/images/%d", game.Images[1].ID), nil)
		rr := httptest.NewRecorder()
		
		router := mux.NewRouter()
		router.HandleFunc("/images/{id}", ctx.Handler.getImage).Methods("GET")
		router.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		
		if rr.Header().Get("Content-Type") != "image/webp" {
			t.Errorf("Expected Content-Type 'image/webp', got '%s'", rr.Header().Get("Content-Type"))
		}
		
		if len(rr.Body.Bytes()) == 0 {
			t.Error("Response body is empty")
		}
	})
	
	// Test non-existent image
	t.Run("Get non-existent image data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/images/99999", nil)
		rr := httptest.NewRecorder()
		
		router := mux.NewRouter()
		router.HandleFunc("/images/{id}", ctx.Handler.getImage).Methods("GET")
		router.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})
	
	// Test invalid image ID
	t.Run("Get image with invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/images/invalid", nil)
		rr := httptest.NewRecorder()
		
		router := mux.NewRouter()
		router.HandleFunc("/images/{id}", ctx.Handler.getImage).Methods("GET")
		router.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}

func TestGameWithNote(t *testing.T) {
	ctx := setupGameFixture(t)
	
	gameNote := "This was a great game!"
	gameReq := models.CreateGameRequest{
		Name:        "Game with Note",
		Date:        "2024-01-16",
		Map:         models.Ptr("Tharsis"),
		Generations: models.Ptr(10),
		Note:        &gameNote,
		Expansions:  models.Ptr(models.Expansions{"base": true}),
		Players: []models.PlayerRequest{
			{Name: "Alice", Corporation: "Ecoline", TerraformingRating: 20},
		},
	}
	
	req := AuthenticatedGameRequest{
		CreateGameRequest: gameReq,
		ActorName:         ctx.AdminName,
		ActorPassword:     ctx.AdminPassword,
	}
	
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	ctx.Handler.createGame(rr, httpReq)
	
	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create game, status: %d, body: %s", rr.Code, rr.Body.String())
	}
	
	var game models.GameWithDetails
	json.NewDecoder(rr.Body).Decode(&game)
	
	// Verify note was saved
	if game.Game.Note == nil || *game.Game.Note != gameNote {
		t.Errorf("Expected note '%s', got %v", gameNote, game.Game.Note)
	}
}

func TestUpdateGameMilestonesAndAwards(t *testing.T) {
	ctx := setupGameFixture(t)
	
	// Step 1: Create initial game with milestones and awards
	winnerIndex := 0
	initialGameReq := models.CreateGameRequest{
		Name:        "Test Update Game",
		Date:        "2024-01-15",
		Map:         models.Ptr("Tharsis"),
		Generations: models.Ptr(10),
		Expansions: models.Ptr(models.Expansions{
			"Base Game":     true,
			"Corporate Era": false,
		}),
		Players: []models.PlayerRequest{
			{
				Name:               "Alice",
				Corporation:        "Ecoline",
				TerraformingRating: 25,
				Cities:             5,
				Greeneries:         7,
				Cards:              15,
				TurmoilPoints:      0,
			},
			{
				Name:               "Bob",
				Corporation:        "Tharsis Republic",
				TerraformingRating: 23,
				Cities:             4,
				Greeneries:         5,
				Cards:              12,
				TurmoilPoints:      0,
			},
		},
		Milestones: []models.MilestoneRequest{
			{Name: "Terraformer", WinnerGamePlayerIndex: &winnerIndex},
			{Name: "Mayor", WinnerGamePlayerIndex: nil},
		},
		Awards: []models.AwardRequest{
			{
				Name: "Landlord",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 0, Placement: models.PlacementFirst},
					{PlayerIndex: 1, Placement: models.PlacementSecond},
				},
			},
			{
				Name:       "Banker",
				Placements: []models.PlacementRequest{},
			},
		},
	}
	
	req := AuthenticatedGameRequest{
		CreateGameRequest: initialGameReq,
		ActorName:         ctx.AdminName,
		ActorPassword:     ctx.AdminPassword,
	}
	
	// Create the game
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	ctx.Handler.createGame(rr, httpReq)
	
	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create game, status: %d, body: %s", rr.Code, rr.Body.String())
	}
	
	var createdGame models.GameWithDetails
	err := json.NewDecoder(rr.Body).Decode(&createdGame)
	if err != nil {
		t.Fatalf("Failed to parse created game: %v", err)
	}
	gameID := createdGame.Game.GameID
	
	// Verify initial state
	if len(createdGame.Milestones) != 2 {
		t.Errorf("Expected 2 milestones, got %d", len(createdGame.Milestones))
	}
	if len(createdGame.Awards) != 2 {
		t.Errorf("Expected 2 awards, got %d", len(createdGame.Awards))
	}
	
	// Count funded awards (those with placements)
	fundedCount := 0
	for _, award := range createdGame.Awards {
		if award.Name == "Landlord" {
			// Check placements via Placements
			placementsForAward := 0
			for _, placement := range createdGame.Placements {
				if placement.AwardID == award.ID {
					placementsForAward++
				}
			}
			if placementsForAward > 0 {
				fundedCount++
			}
		}
	}
	if fundedCount != 1 {
		t.Errorf("Initially expected 1 funded award, got %d", fundedCount)
	}
	
	// Step 2: Update the game with different milestones and awards
	winnerIndex1 := 1
	updatedGameReq := models.CreateGameRequest{
		Name:        "Updated Test Game",
		Date:        "2024-01-15",
		Map:         models.Ptr("Tharsis"),
		Generations: models.Ptr(11),
		Expansions: models.Ptr(models.Expansions{
			"Base Game":     true,
			"Corporate Era": true,
		}),
		Players: []models.PlayerRequest{
			{
				Name:               "Alice",
				Corporation:        "Ecoline",
				TerraformingRating: 30,
				Cities:             6,
				Greeneries:         8,
				Cards:              18,
				TurmoilPoints:      0,
			},
			{
				Name:               "Bob",
				Corporation:        "Tharsis Republic",
				TerraformingRating: 28,
				Cities:             5,
				Greeneries:         6,
				Cards:              15,
				TurmoilPoints:      0,
			},
		},
		Milestones: []models.MilestoneRequest{
			{Name: "Terraformer", WinnerGamePlayerIndex: &winnerIndex},  // Same winner
			{Name: "Mayor", WinnerGamePlayerIndex: &winnerIndex1},       // Now achieved by Bob
			{Name: "Gardener", WinnerGamePlayerIndex: &winnerIndex},     // New milestone
		},
		Awards: []models.AwardRequest{
			{
				Name: "Landlord",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 1, Placement: models.PlacementFirst},   // Swapped
					{PlayerIndex: 0, Placement: models.PlacementSecond},  // Swapped
				},
			},
			{
				Name: "Banker",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 0, Placement: models.PlacementFirst}, // Now funded
				},
			},
			{
				Name: "Scientist", // New award
				Placements: []models.PlacementRequest{
					{PlayerIndex: 1, Placement: models.PlacementFirst},
				},
			},
		},
	}
	
	updateReq := AuthenticatedGameRequest{
		CreateGameRequest: updatedGameReq,
		ActorName:         ctx.AdminName,
		ActorPassword:     ctx.AdminPassword,
	}
	
	// Update the game
	body, _ = json.Marshal(updateReq)
	httpReq, _ = http.NewRequest("PUT", fmt.Sprintf("/games/%d", gameID), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", gameID)})
	
	rr = httptest.NewRecorder()
	ctx.Handler.updateGame(rr, httpReq)
	
	if rr.Code != http.StatusOK {
		t.Fatalf("Failed to update game, status: %d, body: %s", rr.Code, rr.Body.String())
	}
	
	var updatedGame models.GameWithDetails
	err = json.NewDecoder(rr.Body).Decode(&updatedGame)
	if err != nil {
		t.Fatalf("Failed to parse updated game: %v", err)
	}
	
	// Step 3: Verify the updates
	
	// Verify milestones
	if len(updatedGame.Milestones) != 3 {
		t.Errorf("Expected 3 milestones after update, got %d", len(updatedGame.Milestones))
	}
	
	milestoneNames := make(map[string]bool)
	for _, m := range updatedGame.Milestones {
		milestoneNames[m.Name] = true
		
		// Check Mayor milestone winner changed
		if m.Name == "Mayor" {
			if m.WinnerGamePlayerID == nil {
				t.Error("Mayor should have a winner after update")
			} else {
				// Find the winner player
				for _, gp := range updatedGame.GamePlayers {
					if gp.ID == *m.WinnerGamePlayerID {
						player, _ := ctx.Repo.GetPlayerByID(gp.PlayerID)
						if player.Name != "Bob" {
							t.Errorf("Mayor winner should be Bob, got %s", player.Name)
						}
					}
				}
			}
		}
	}
	
	if !milestoneNames["Terraformer"] {
		t.Error("Terraformer milestone missing after update")
	}
	if !milestoneNames["Mayor"] {
		t.Error("Mayor milestone missing after update")
	}
	if !milestoneNames["Gardener"] {
		t.Error("New milestone Gardener should be added")
	}
	
	// Verify awards
	if len(updatedGame.Awards) != 3 {
		t.Errorf("Expected 3 awards after update, got %d", len(updatedGame.Awards))
	}
	
	awardNames := make(map[string]bool)
	for _, a := range updatedGame.Awards {
		awardNames[a.Name] = true
	}
	
	if !awardNames["Landlord"] {
		t.Error("Landlord award missing after update")
	}
	if !awardNames["Banker"] {
		t.Error("Banker award missing after update")
	}
	if !awardNames["Scientist"] {
		t.Error("New award Scientist should be added")
	}
	
	// Verify award placements changed
	for _, award := range updatedGame.Awards {
		if award.Name == "Landlord" {
			// Count placements for this award
			placementCount := 0
			for _, p := range updatedGame.Placements {
				if p.AwardID == award.ID {
					placementCount++
					// Check that placements are swapped
					for _, gp := range updatedGame.GamePlayers {
						if gp.ID == p.GamePlayerID {
							player, _ := ctx.Repo.GetPlayerByID(gp.PlayerID)
							if player.Name == "Bob" && p.Placement != models.PlacementFirst {
								t.Error("Bob should have first place in Landlord after update")
							}
							if player.Name == "Alice" && p.Placement != models.PlacementSecond {
								t.Error("Alice should have second place in Landlord after update")
							}
						}
					}
				}
			}
			if placementCount != 2 {
				t.Errorf("Landlord should still have 2 placements, got %d", placementCount)
			}
		}
		
		if award.Name == "Banker" {
			// Check that Banker is now funded
			placementCount := 0
			for _, p := range updatedGame.Placements {
				if p.AwardID == award.ID {
					placementCount++
				}
			}
			if placementCount != 1 {
				t.Errorf("Banker should now have 1 placement, got %d", placementCount)
			}
		}
	}
	
	// Verify other game fields updated
	if updatedGame.Game.Generations == nil || *updatedGame.Game.Generations != 11 {
		var gen int
		if updatedGame.Game.Generations != nil {
			gen = *updatedGame.Game.Generations
		}
		t.Errorf("Generations should be updated to 11, got %d", gen)
	}

	if updatedGame.Game.Expansions == nil || !(*updatedGame.Game.Expansions)["Corporate Era"] {
		t.Error("Corporate Era expansion should be enabled after update")
	}
	
	// Verify player scores updated
	for _, gp := range updatedGame.GamePlayers {
		player, _ := ctx.Repo.GetPlayerByID(gp.PlayerID)
		if player.Name == "Alice" {
			if gp.TerraformingRating != 30 {
				t.Errorf("Alice's TR should be 30, got %d", gp.TerraformingRating)
			}
		}
		if player.Name == "Bob" {
			if gp.TerraformingRating != 28 {
				t.Errorf("Bob's TR should be 28, got %d", gp.TerraformingRating)
			}
		}
	}
}

