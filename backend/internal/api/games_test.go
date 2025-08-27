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
		if adminGame.Map != "Tharsis" {
			t.Errorf("Expected map 'Tharsis', got '%s'", adminGame.Map)
		}
		if adminGame.Generations != 10 {
			t.Errorf("Expected 10 generations, got %d", adminGame.Generations)
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
		if userGame.Map != "Hellas" {
			t.Errorf("Expected map 'Hellas', got '%s'", userGame.Map)
		}
		if userGame.Generations != 12 {
			t.Errorf("Expected 12 generations, got %d", userGame.Generations)
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
		Map:         "Tharsis",
		Generations: 8,
		Expansions:  models.Expansions{"Prelude": true},
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
				WinnerGamePlayerIndex: intPtr(0),
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
		if game.Game.Map != "Tharsis" {
			t.Errorf("Expected map 'Tharsis', got '%s'", game.Game.Map)
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
	image1 := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A}
	image2 := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	
	// Create a game with images
	gameReq := models.CreateGameRequest{
		Name:        "Game with Images",
		Date:        "2024-01-16",
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
	
	if game.Images[0].MimeType != "image/png" {
		t.Errorf("Expected first image mime type 'image/png', got '%s'", game.Images[0].MimeType)
	}
	
	if game.Images[1].MimeType != "image/jpeg" {
		t.Errorf("Expected second image mime type 'image/jpeg', got '%s'", game.Images[1].MimeType)
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
		
		if rr.Header().Get("Content-Type") != "image/png" {
			t.Errorf("Expected Content-Type 'image/png', got '%s'", rr.Header().Get("Content-Type"))
		}
		
		if rr.Header().Get("Cache-Control") != "public, max-age=3600" {
			t.Errorf("Expected Cache-Control header")
		}
		
		if !bytes.Equal(rr.Body.Bytes(), image1) {
			t.Error("Response body doesn't match expected image data")
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
		
		if rr.Header().Get("Content-Type") != "image/jpeg" {
			t.Errorf("Expected Content-Type 'image/jpeg', got '%s'", rr.Header().Get("Content-Type"))
		}
		
		if !bytes.Equal(rr.Body.Bytes(), image2) {
			t.Error("Response body doesn't match expected image data")
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
		Map:         "Tharsis",
		Generations: 10,
		Note:        &gameNote,
		Expansions:  models.Expansions{"base": true},
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

