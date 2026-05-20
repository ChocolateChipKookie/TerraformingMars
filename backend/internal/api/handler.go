package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"terraforming-mars-backend/internal/database"
	"terraforming-mars-backend/internal/rating"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo   *database.Repository
	rating *rating.Service
}

func NewHandler(db *sql.DB) *Handler {
	repo := database.NewRepository(db)
	return &Handler{
		repo:   repo,
		rating: rating.NewService(repo.GetGamesForRating),
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Player routes - following REST conventions
	router.HandleFunc("/players", h.getPlayers).Methods("GET")                                 // List all players
	router.HandleFunc("/players", h.createPlayer).Methods("POST")                              // Create new player
	router.HandleFunc("/players/{id}/extended", h.getPlayerExtendedInfo).Methods("GET")        // Get player with extended info
	router.HandleFunc("/players/{id}/rating-history", h.getPlayerRatingHistory).Methods("GET") // Player's per-game rating timeline
	router.HandleFunc("/players/{id}", h.getPlayer).Methods("GET")                             // Get specific player
	router.HandleFunc("/players/{id}", h.updatePlayer).Methods("PUT")                          // Update existing player

	// Game routes - following REST conventions
	router.HandleFunc("/games", h.getGames).Methods("GET")           // List all games
	router.HandleFunc("/games", h.createGame).Methods("POST")        // Create new game
	router.HandleFunc("/games/{id}", h.getGame).Methods("GET")       // Get specific game
	router.HandleFunc("/games/{id}", h.updateGame).Methods("PUT")    // Update existing game (new revision)
	router.HandleFunc("/games/{id}", h.deleteGame).Methods("DELETE") // Delete game

	// Image routes
	router.HandleFunc("/images/{id}", h.getImage).Methods("GET") // Get image data by image ID
}

// Helper function to send JSON responses
func (h *Handler) sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function to send error responses
func (h *Handler) sendError(w http.ResponseWriter, status int, message string) {
	h.sendJSON(w, status, map[string]string{"error": message})
}
