package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"terraforming-mars-backend/internal/database"
)

func main() {
	var dbPath string
	var adminName string
	var adminPassword string

	flag.StringVar(&dbPath, "db", "data/terraforming_mars.db", "Database file path")
	flag.StringVar(&adminName, "admin", "", "Admin username (required)")
	flag.StringVar(&adminPassword, "password", "", "Admin password (required)")
	flag.Parse()

	if adminName == "" || adminPassword == "" {
		fmt.Println("Usage: init-db -admin <username> -password <password> [-db <path>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check if parent directory exists
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Parent directory does not exist: %s\nPlease create it first or specify a different path with -db flag", dir)
	}

	// Initialize database
	db, err := database.Init(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create repository
	repo := database.NewRepository(db)

	// Create admin user using the existing CreateSystemAdmin function
	admin, err := repo.CreateSystemAdmin(adminName, &adminPassword)
	if err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	fmt.Printf("Successfully initialized database at %s\n", dbPath)
	fmt.Printf("Created admin user: %s (ID: %d)\n", admin.Name, admin.ID)
}