package api

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	"terraforming-mars-backend/internal/database"
	"terraforming-mars-backend/internal/models"
)

// TestContext holds common test data
type TestContext struct {
	Handler       *Handler
	Repo          *database.Repository
	Admin         *models.Player
	AdminName     string
	AdminPassword string
}

// setupTestAPI creates a test database and API handler with admin user
func setupTestAPI(t *testing.T) *TestContext {
	// Create in-memory database
	db, err := database.Init(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	// Create repository and handler
	repo := database.NewRepository(db)
	handler := NewHandler(db)

	// Create a system admin for testing
	adminName := "admin"
	adminPassword := "adminpass123"
	admin, err := repo.CreateSystemAdmin(adminName, &adminPassword)
	if err != nil {
		t.Fatalf("Failed to create system admin: %v", err)
	}

	return &TestContext{
		Handler:       handler,
		Repo:          repo,
		Admin:         admin,
		AdminName:     adminName,
		AdminPassword: adminPassword,
	}
}

// TestFixture holds both context and test players
type TestFixture struct {
	*TestContext
	Alice    *models.Player
	Bob      *models.Player
	Charlie  *models.Player
	UserPass string
}

// TestGameFixture holds context, players, and test games
type TestGameFixture struct {
	*TestFixture
	Dave  *models.Player
	Eve   *models.Player
	Game1 *models.GameWithDetails
	Game2 *models.GameWithDetails
}

// setupTestFixture creates a fresh test database, context, and standard test players
func setupTestFixture(t *testing.T) *TestFixture {
	ctx := setupTestAPI(t)
	userPass := "userpass123"

	alice, _ := ctx.Repo.CreatePlayer("Alice", nil, models.RolePlayer, *ctx.Admin)
	bob, _ := ctx.Repo.CreatePlayer("Bob", &userPass, models.RoleUser, *ctx.Admin)
	charlie, _ := ctx.Repo.CreatePlayer("Charlie", nil, models.RolePlayer, *bob)

	return &TestFixture{
		TestContext: ctx,
		Alice:       alice,
		Bob:         bob,
		Charlie:     charlie,
		UserPass:    userPass,
	}
}

// setupGameFixture creates a fresh test database, context, players, and standard test games
func setupGameFixture(t *testing.T) *TestGameFixture {
	fixture := setupTestFixture(t)

	// Create additional users for testing game permissions
	dave, _ := fixture.Repo.CreatePlayer("Dave", &fixture.UserPass, models.RoleUser, *fixture.Admin)
	eve, _ := fixture.Repo.CreatePlayer("Eve", &fixture.UserPass, models.RoleUser, *fixture.Admin)

	// Create test game 1 - created by admin with multiple players
	gameReq1 := models.CreateGameRequest{
		Name:        "Admin Game",
		Date:        "2024-01-15",
		Map:         "Tharsis",
		Generations: 10,
		Expansions:  models.Expansions{"Venus Next": true, "Prelude": true},
		Players: []models.PlayerRequest{
			{
				Name:               "Alice",
				Corporation:        "Tharsis Republic",
				TerraformingRating: 20,
				Cities:             2,
				Greeneries:         3,
				Cards:              10,
				TurmoilPoints:      5,
			},
			{
				Name:               "Bob",
				Corporation:        "Saturn Systems",
				TerraformingRating: 25,
				Cities:             1,
				Greeneries:         4,
				Cards:              8,
				TurmoilPoints:      3,
			},
			{
				Name:               "Charlie",
				Corporation:        "Mining Guild",
				TerraformingRating: 18,
				Cities:             3,
				Greeneries:         2,
				Cards:              12,
				TurmoilPoints:      7,
			},
		},
		Milestones: []models.MilestoneRequest{
			{
				Name:                  "Terraformer",
				WinnerGamePlayerIndex: models.Ptr(0), // Alice (index 0)
			},
			{
				Name:                  "Mayor",
				WinnerGamePlayerIndex: models.Ptr(2), // Charlie (index 2)
			},
			{
				Name:                  "Gardener",
				WinnerGamePlayerIndex: models.Ptr(1), // Bob (index 1)
			},
		},
		Awards: []models.AwardRequest{
			{
				Name: "Landlord",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 2, Placement: models.PlacementFirst},  // Charlie - 1st place
					{PlayerIndex: 0, Placement: models.PlacementSecond}, // Alice - 2nd place
				},
			},
			{
				Name: "Scientist",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 2, Placement: models.PlacementFirst},  // Charlie - 1st place
					{PlayerIndex: 1, Placement: models.PlacementSecond}, // Bob - 2nd place
				},
			},
			{
				Name: "Thermalist",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 1, Placement: models.PlacementFirst}, // Bob - 1st place
				},
			},
		},
	}

	// Create test game 2 - created by Bob (user) with 2 other users participating
	// This is key for testing that participants cannot modify games they didn't create
	gameReq2 := models.CreateGameRequest{
		Name:        "User Game With Multiple Users",
		Date:        "2024-01-20",
		Map:         "Hellas",
		Generations: 12,
		Expansions:  models.Expansions{"Colonies": true},
		Players: []models.PlayerRequest{
			{
				Name:               "Bob",
				Corporation:        "Helion",
				TerraformingRating: 22,
				Cities:             2,
				Greeneries:         5,
				Cards:              9,
				TurmoilPoints:      4,
			},
			{
				Name:               "Dave",
				Corporation:        "Inventrix",
				TerraformingRating: 19,
				Cities:             1,
				Greeneries:         3,
				Cards:              11,
				TurmoilPoints:      6,
			},
			{
				Name:               "Eve",
				Corporation:        "Robinson Industries",
				TerraformingRating: 21,
				Cities:             4,
				Greeneries:         1,
				Cards:              7,
				TurmoilPoints:      2,
			},
		},
		Milestones: []models.MilestoneRequest{
			{
				Name:                  "Diversifier",
				WinnerGamePlayerIndex: models.Ptr(1), // Dave (index 1)
			},
			{
				Name:                  "Tactician",
				WinnerGamePlayerIndex: models.Ptr(0), // Bob (index 0)
			},
			{
				Name:                  "Polar Explorer",
				WinnerGamePlayerIndex: models.Ptr(2), // Eve (index 2)
			},
		},
		Awards: []models.AwardRequest{
			{
				Name: "Banker",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 1, Placement: models.PlacementFirst},  // Dave - 1st place
					{PlayerIndex: 0, Placement: models.PlacementSecond}, // Bob - 2nd place
				},
			},
			{
				Name: "Space Baron",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 2, Placement: models.PlacementFirst},  // Eve - 1st place
					{PlayerIndex: 1, Placement: models.PlacementSecond}, // Dave - 2nd place
				},
			},
			{
				Name: "Excentric",
				Placements: []models.PlacementRequest{
					{PlayerIndex: 0, Placement: models.PlacementFirst}, // Bob - 1st place
				},
			},
		},
	}

	game1, _ := fixture.Repo.CreateGame(&models.ParsedGameRequest{Normal: &gameReq1}, *fixture.Admin)
	game2, _ := fixture.Repo.CreateGame(&models.ParsedGameRequest{Normal: &gameReq2}, *fixture.Bob)

	return &TestGameFixture{
		TestFixture: fixture,
		Dave:        dave,
		Eve:         eve,
		Game1:       game1,
		Game2:       game2,
	}
}

// createTestImage creates a minimal valid image for testing in the specified format
func createTestImage(mimeType string) []byte {
	// Create a 1x1 pixel image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // Red pixel

	var buf bytes.Buffer

	switch mimeType {
	case "image/jpeg":
		jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	case "image/png":
		png.Encode(&buf, img)
	default:
		// Default to PNG
		png.Encode(&buf, img)
	}

	return buf.Bytes()
}
