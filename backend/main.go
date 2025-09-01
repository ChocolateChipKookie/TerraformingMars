package main

import (
	"log"
	"net/http"

	"terraforming-mars-backend/internal/api"
	"terraforming-mars-backend/internal/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize database with file path
	db, err := database.Init("data/terraforming_mars.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create API handler
	apiHandler := api.NewHandler(db)

	// Setup routes
	router := mux.NewRouter()
	
	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiHandler.RegisterRoutes(apiRouter)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"}, // React and Vite dev servers
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Start server
	handler := c.Handler(router)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}