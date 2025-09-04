package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"terraforming-mars-backend/internal/models"
)

// Request structure for creating a player
type CreatePlayerRequest struct {
	// Player details
	Name     string  `json:"name"`
	Password *string `json:"password,omitempty"`
	Role     string  `json:"role"`

	// Actor authentication (who is creating this player)
	ActorName     string `json:"actor_name"`
	ActorPassword string `json:"actor_password"`
}

// Request structure for updating a player
type UpdatePlayerRequest struct {
	// Player details (all optional for partial updates)
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`

	// Actor authentication (who is updating this player)
	ActorName     string `json:"actor_name"`
	ActorPassword string `json:"actor_password"`
}

// GET /players - List all players
func (h *Handler) getPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.repo.GetAllPlayers()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to fetch players")
		return
	}

	h.sendJSON(w, http.StatusOK, players)
}

// GET /players/{id} - Get a specific player
func (h *Handler) getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	player, err := h.repo.GetPlayerByID(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Player not found")
		return
	}

	h.sendJSON(w, http.StatusOK, player)
}

// POST /players - Create a new player
func (h *Handler) createPlayer(w http.ResponseWriter, r *http.Request) {
	var req CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if this is the first player in the system
	players, err := h.repo.GetAllPlayers()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to check existing players")
		return
	}

	// If no players exist, create the first one as admin without authentication
	if len(players) == 0 {
		// Force the role to admin for the first player
		player, err := h.repo.CreateSystemAdmin(req.Name, req.Password)
		if err != nil {
			h.sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.sendJSON(w, http.StatusCreated, player)
		return
	}

	// Authenticate the actor
	actor, err := h.repo.AuthenticatePlayer(req.ActorName, req.ActorPassword)
	if err != nil {
		h.sendError(w, http.StatusUnauthorized, "Invalid actor credentials")
		return
	}

	// Create the player (repository will handle role validation)
	player, err := h.repo.CreatePlayer(req.Name, req.Password, models.PlayerRole(req.Role), *actor)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.sendJSON(w, http.StatusCreated, player)
}

// PUT /players/{id} - Update an existing player
func (h *Handler) updatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	var req UpdatePlayerRequest
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

	// Get current player to preserve unchanged fields
	currentPlayer, err := h.repo.GetPlayerByID(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Player not found")
		return
	}

	// Use provided name or keep existing
	name := currentPlayer.Name
	if req.Name != nil {
		name = *req.Name
	}

	// Convert role string to PlayerRole if provided
	var role *models.PlayerRole
	if req.Role != nil {
		r := models.PlayerRole(*req.Role)
		role = &r
	}

	// Update the player
	player, err := h.repo.UpdatePlayer(id, name, req.Password, role, *actor)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.sendJSON(w, http.StatusOK, player)
}

