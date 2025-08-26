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

func TestGetPlayers(t *testing.T) {
	ctx := setupTestFixture(t)
	
	// Test GET /players
	req, _ := http.NewRequest("GET", "/players", nil)
	rr := httptest.NewRecorder()
	ctx.Handler.getPlayers(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	
	var players []models.Player
	err := json.NewDecoder(rr.Body).Decode(&players)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if len(players) != 4 { // admin + Alice + Bob + Charlie
		t.Errorf("Expected 4 players, got %d", len(players))
	}
	
	// Verify player names, roles, and created_by
	playersByName := make(map[string]models.Player)
	for _, p := range players {
		playersByName[p.Name] = p
	}
	
	// Check admin (created_by should be nil for system admin)
	if admin, ok := playersByName["admin"]; !ok {
		t.Error("Admin player not found")
	} else {
		if admin.Role != models.RoleAdmin {
			t.Errorf("Admin role: expected '%s', got '%s'", models.RoleAdmin, admin.Role)
		}
		if admin.CreatedBy != nil {
			t.Errorf("Admin created_by: expected nil, got %v", admin.CreatedBy)
		}
	}
	
	// Check Alice (created by admin)
	if aliceFromAPI, ok := playersByName["Alice"]; !ok {
		t.Error("Alice player not found")
	} else {
		if aliceFromAPI.Role != models.RolePlayer {
			t.Errorf("Alice role: expected '%s', got '%s'", models.RolePlayer, aliceFromAPI.Role)
		}
		if aliceFromAPI.CreatedBy == nil || *aliceFromAPI.CreatedBy != ctx.Admin.ID {
			t.Errorf("Alice created_by: expected %d, got %v", ctx.Admin.ID, aliceFromAPI.CreatedBy)
		}
		if aliceFromAPI.ID != ctx.Alice.ID {
			t.Errorf("Alice ID mismatch: expected %d, got %d", ctx.Alice.ID, aliceFromAPI.ID)
		}
	}
	
	// Check Bob (created by admin)
	if bobFromAPI, ok := playersByName["Bob"]; !ok {
		t.Error("Bob player not found")
	} else {
		if bobFromAPI.Role != models.RoleUser {
			t.Errorf("Bob role: expected '%s', got '%s'", models.RoleUser, bobFromAPI.Role)
		}
		if bobFromAPI.CreatedBy == nil || *bobFromAPI.CreatedBy != ctx.Admin.ID {
			t.Errorf("Bob created_by: expected %d, got %v", ctx.Admin.ID, bobFromAPI.CreatedBy)
		}
		if bobFromAPI.ID != ctx.Bob.ID {
			t.Errorf("Bob ID mismatch: expected %d, got %d", ctx.Bob.ID, bobFromAPI.ID)
		}
	}
	
	// Check Charlie (created by Bob)
	if charlieFromAPI, ok := playersByName["Charlie"]; !ok {
		t.Error("Charlie player not found")
	} else {
		if charlieFromAPI.Role != models.RolePlayer {
			t.Errorf("Charlie role: expected '%s', got '%s'", models.RolePlayer, charlieFromAPI.Role)
		}
		if charlieFromAPI.CreatedBy == nil || *charlieFromAPI.CreatedBy != ctx.Bob.ID {
			t.Errorf("Charlie created_by: expected %d (Bob's ID), got %v", ctx.Bob.ID, charlieFromAPI.CreatedBy)
		}
		if charlieFromAPI.ID != ctx.Charlie.ID {
			t.Errorf("Charlie ID mismatch: expected %d, got %d", ctx.Charlie.ID, charlieFromAPI.ID)
		}
	}
}

func TestGetPlayer(t *testing.T) {
	ctx := setupTestAPI(t)
	
	// Create a test player
	player, _ := ctx.Repo.CreatePlayer("Alice", nil, models.RolePlayer, *ctx.Admin)
	
	// Test GET /players/{id} with valid ID
	t.Run("Valid player ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/players/%d", player.ID), nil)
		// Need to use mux to set the URL vars
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", player.ID)})
		
		rr := httptest.NewRecorder()
		ctx.Handler.getPlayer(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, status)
		}
		
		var returnedPlayer models.Player
		err := json.NewDecoder(rr.Body).Decode(&returnedPlayer)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		
		if returnedPlayer.ID != player.ID {
			t.Errorf("Expected player ID %d, got %d", player.ID, returnedPlayer.ID)
		}
		if returnedPlayer.Name != "Alice" {
			t.Errorf("Expected player name 'Alice', got '%s'", returnedPlayer.Name)
		}
		if returnedPlayer.Role != models.RolePlayer {
			t.Errorf("Expected role '%s', got '%s'", models.RolePlayer, returnedPlayer.Role)
		}
		if returnedPlayer.CreatedBy == nil || *returnedPlayer.CreatedBy != ctx.Admin.ID {
			t.Errorf("Expected created_by %d, got %v", ctx.Admin.ID, returnedPlayer.CreatedBy)
		}
	})
	
	// Test with non-existent ID
	t.Run("Non-existent player ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/players/99999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "99999"})
		
		rr := httptest.NewRecorder()
		ctx.Handler.getPlayer(rr, req)
		
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test with invalid ID format
	t.Run("Invalid player ID format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/players/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		
		rr := httptest.NewRecorder()
		ctx.Handler.getPlayer(rr, req)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
}

func TestCreatePlayer(t *testing.T) {
	ctx := setupTestAPI(t)
	
	// Test successful player creation
	t.Run("Create player successfully", func(t *testing.T) {
		req := CreatePlayerRequest{
			Name:          "Alice",
			Password:      nil,
			Role:          "player",
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
		}
		
		var player models.Player
		err := json.NewDecoder(rr.Body).Decode(&player)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		
		if player.Name != "Alice" {
			t.Errorf("Expected name 'Alice', got '%s'", player.Name)
		}
		if player.Role != models.RolePlayer {
			t.Errorf("Expected role '%s', got '%s'", models.RolePlayer, player.Role)
		}
		if player.CreatedBy == nil || *player.CreatedBy != ctx.Admin.ID {
			t.Errorf("Expected created_by %d, got %v", ctx.Admin.ID, player.CreatedBy)
		}
	})
	
	// Test creating user with password
	t.Run("Create user with password", func(t *testing.T) {
		userPass := "bobpass123"
		req := CreatePlayerRequest{
			Name:          "Bob",
			Password:      &userPass,
			Role:          "user",
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
		}
		
		var player models.Player
		json.NewDecoder(rr.Body).Decode(&player)
		
		if player.Role != models.RoleUser {
			t.Errorf("Expected role '%s', got '%s'", models.RoleUser, player.Role)
		}
	})
	
	// Test invalid actor credentials
	t.Run("Invalid actor credentials", func(t *testing.T) {
		req := CreatePlayerRequest{
			Name:          "Charlie",
			Password:      nil,
			Role:          "player",
			ActorName:     ctx.AdminName,
			ActorPassword: "wrongpassword",
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test invalid role
	t.Run("Invalid role", func(t *testing.T) {
		req := CreatePlayerRequest{
			Name:          "Dave",
			Password:      nil,
			Role:          "invalidrole",
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test user requires password
	t.Run("User role requires password", func(t *testing.T) {
		req := CreatePlayerRequest{
			Name:          "Eve",
			Password:      nil,
			Role:          "user",
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test player cannot have password
	t.Run("Player role cannot have password", func(t *testing.T) {
		badPass := "shouldnotwork"
		req := CreatePlayerRequest{
			Name:          "Frank",
			Password:      &badPass,
			Role:          "player",
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test that user cannot create another user (only admin can)
	t.Run("User cannot create user role", func(t *testing.T) {
		// First create a user
		userPass := "userpass123"
		ctx.Repo.CreatePlayer("TestUser", &userPass, models.RoleUser, *ctx.Admin)
		
		// Now try to use that user to create another user
		newUserPass := "newuserpass123"
		req := CreatePlayerRequest{
			Name:          "ShouldFail",
			Password:      &newUserPass,
			Role:          "user",
			ActorName:     "TestUser",
			ActorPassword: userPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test that user CAN create player role (this should work)
	t.Run("User can create player role", func(t *testing.T) {
		userPass := "userpass123"
		
		req := CreatePlayerRequest{
			Name:          "PlayerByUser",
			Password:      nil,
			Role:          "player",
			ActorName:     "TestUser", // Use the user created in previous test
			ActorPassword: userPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.createPlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
		}
		
		var player models.Player
		json.NewDecoder(rr.Body).Decode(&player)
		
		if player.Name != "PlayerByUser" {
			t.Errorf("Expected name 'PlayerByUser', got '%s'", player.Name)
		}
		if player.Role != models.RolePlayer {
			t.Errorf("Expected role '%s', got '%s'", models.RolePlayer, player.Role)
		}
	})
}

func TestUpdatePlayer(t *testing.T) {
	// Test admin updating any player
	t.Run("Admin can update any player", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Alice Updated"),
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Alice.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Alice.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, status)
		}
		
		var player models.Player
		json.NewDecoder(rr.Body).Decode(&player)
		
		if player.Name != "Alice Updated" {
			t.Errorf("Expected name 'Alice Updated', got '%s'", player.Name)
		}
	})
	
	// Test user updating themselves (password only, not name)
	t.Run("User can update themselves", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		newPassword := "bobsnewpass456"
		req := UpdatePlayerRequest{
			Password:      &newPassword,
			ActorName:     "Bob",
			ActorPassword: ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Bob.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Bob.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, status)
		}
		
		var player models.Player
		json.NewDecoder(rr.Body).Decode(&player)
		
		// Name should remain unchanged
		if player.Name != "Bob" {
			t.Errorf("Expected name 'Bob' (unchanged), got '%s'", player.Name)
		}
		
		// Verify new password works
		_, err := ctx.Repo.AuthenticatePlayer("Bob", newPassword)
		if err != nil {
			t.Errorf("New password should work for authentication: %v", err)
		}
	})
	
	// Test user cannot update player names even for players they created
	t.Run("User cannot update player names even for players they created", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Charlie Updated"),
			ActorName:     "Bob",
			ActorPassword: ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Charlie.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Charlie.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test user cannot update other players
	t.Run("User cannot update players they didn't create", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Alice Forbidden"),
			ActorName:     "Bob",
			ActorPassword: ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Alice.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Alice.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test player cannot update anything
	t.Run("Player cannot update anything", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Should Fail"),
			ActorName:     "Alice",
			ActorPassword: "", // Players don't have passwords
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Alice.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Alice.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
		}
	})
	
	// Test invalid actor credentials
	t.Run("Invalid actor credentials", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Should Fail"),
			ActorName:     "Bob",
			ActorPassword: "wrongpassword",
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Bob.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Bob.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
		}
	})
	
	// Test non-existent player
	t.Run("Non-existent player ID", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Should Fail"),
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", "/players/99999", bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": "99999"})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
		}
	})
	
	// Test invalid player ID format
	t.Run("Invalid player ID format", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Should Fail"),
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", "/players/invalid", bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": "invalid"})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
	})
	
	// Test role changes (only admin can change roles)
	t.Run("Admin can change player role", func(t *testing.T) {
		ctx := setupTestFixture(t)
		newUserPass := "alicenewpass123"
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Alice"),
			Role:          stringPtr("user"),
			Password:      &newUserPass,
			ActorName:     ctx.AdminName,
			ActorPassword: ctx.AdminPassword,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Alice.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Alice.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, status)
		}
		
		var player models.Player
		json.NewDecoder(rr.Body).Decode(&player)
		
		if player.Role != models.RoleUser {
			t.Errorf("Expected role '%s', got '%s'", models.RoleUser, player.Role)
		}
	})
	
	// Test user cannot change roles
	t.Run("User cannot change roles", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Bob"),
			Role:          stringPtr("admin"),
			ActorName:     "Bob",
			ActorPassword: ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Bob.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Bob.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
	})
	
	// Test cannot rename to existing player name
	t.Run("Cannot rename to existing player name", func(t *testing.T) {
		ctx := setupTestFixture(t)
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Alice"), // Alice already exists
			ActorName:     "Bob",
			ActorPassword: ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Bob.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Bob.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test user can change their own password
	t.Run("User can change their own password", func(t *testing.T) {
		ctx := setupTestFixture(t)
		newPassword := "bobsnewpassword456"
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Bob"),
			Password:      &newPassword,
			ActorName:     "Bob",
			ActorPassword: ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Bob.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Bob.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, status)
		}
		
		// Verify new password works by trying to authenticate with it
		_, err := ctx.Repo.AuthenticatePlayer("Bob", newPassword)
		if err != nil {
			t.Errorf("New password should work for authentication: %v", err)
		}
	})
	
	// Test cannot set password for player role
	t.Run("Cannot set password for player role", func(t *testing.T) {
		ctx := setupTestFixture(t)
		badPassword := "shouldnotwork"
		
		req := UpdatePlayerRequest{
			Name:          stringPtr("Charlie"),
			Password:      &badPassword,
			ActorName:     "Bob",
			ActorPassword: ctx.UserPass,
		}
		
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("PUT", fmt.Sprintf("/players/%d", ctx.Charlie.ID), bytes.NewBuffer(body))
		httpReq = mux.SetURLVars(httpReq, map[string]string{"id": fmt.Sprintf("%d", ctx.Charlie.ID)})
		httpReq.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		ctx.Handler.updatePlayer(rr, httpReq)
		
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
		
		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})
}