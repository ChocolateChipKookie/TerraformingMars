package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"terraforming-mars-backend/internal/models"
)

// Request structure for creating/updating a game with authentication
type AuthenticatedGameRequest struct {
	models.CreateGameRequest
	
	// Actor authentication (who is creating/updating this game)
	ActorName     string `json:"actor_name"`
	ActorPassword string `json:"actor_password"`
}

// GET /games - List all games (latest revisions only)
func (h *Handler) getGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.repo.GetAllGames()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to fetch games")
		return
	}
	
	h.sendJSON(w, http.StatusOK, games)
}

// GET /games/{id} - Get a specific game with all details
func (h *Handler) getGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	
	game, err := h.repo.GetGameByID(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Game not found")
		return
	}
	
	h.sendJSON(w, http.StatusOK, game)
}

// POST /games - Create a new game
func (h *Handler) createGame(w http.ResponseWriter, r *http.Request) {
	var req AuthenticatedGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Authenticate the actor
	actor, err := h.repo.AuthenticatePlayer(req.ActorName, req.ActorPassword)
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, "Invalid actor credentials")
		return
	}
	
	// Create the game
	game, err := h.repo.CreateGame(req.CreateGameRequest, *actor)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	h.sendJSON(w, http.StatusCreated, game)
}

// PUT /games/{id} - Update a game (creates a new revision)
func (h *Handler) updateGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}
	
	var req AuthenticatedGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Authenticate the actor
	actor, err := h.repo.AuthenticatePlayer(req.ActorName, req.ActorPassword)
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, "Invalid actor credentials")
		return
	}
	
	// Update the game (creates new revision)
	game, err := h.repo.UpdateGame(id, req.CreateGameRequest, *actor)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	h.sendJSON(w, http.StatusOK, game)
}

// GET /images/{id} - Get image data by image ID
func (h *Handler) getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}
	
	imageData, mimeType, err := h.repo.GetGameImageData(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Image not found")
		return
	}
	
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(imageData)
}