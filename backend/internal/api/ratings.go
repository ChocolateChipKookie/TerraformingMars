package api

import (
	"log"
	"net/http"
	"strconv"

	"terraforming-mars-backend/internal/rating"

	"github.com/gorilla/mux"
)

// historyEntry is one row of /players/{id}/rating-history, joined with game metadata.
type historyEntry struct {
	rating.Entry
	GameName string `json:"game_name"`
}

// GET /players/{id}/rating-history — that player's per-game rating timeline.
func (h *Handler) getPlayerRatingHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	snap, err := h.rating.Snapshot()
	if err != nil {
		log.Printf("rating snapshot failed: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to compute ratings")
		return
	}

	entries := snap.Timeline[playerID]
	if len(entries) == 0 {
		h.sendJSON(w, http.StatusOK, []historyEntry{})
		return
	}

	// Build a game_id → game_name lookup from the games list.
	games, err := h.repo.GetAllGames()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to load games")
		return
	}
	nameByGameID := make(map[int]string, len(games))
	for _, g := range games {
		nameByGameID[g.GameID] = g.Name
	}

	out := make([]historyEntry, len(entries))
	for i, e := range entries {
		out[i] = historyEntry{Entry: e, GameName: nameByGameID[e.GameID]}
	}

	h.sendJSON(w, http.StatusOK, out)
}
