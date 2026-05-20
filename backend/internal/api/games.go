package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"terraforming-mars-backend/internal/imageutil"
	"terraforming-mars-backend/internal/models"
)

// RequestAuthentication contains the actor credentials for authenticated requests
type RequestAuthentication struct {
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

// playerRatingDelta is the per-participant rating slice attached to a game-detail response.
type playerRatingDelta struct {
	PlayerID     int     `json:"player_id"`
	Ordinal      float64 `json:"ordinal"`
	DeltaOrdinal float64 `json:"delta_ordinal"`
}

// gameDetailResponse wraps GameWithDetails with the rating data the frontend uses
// to render per-player delta on the game view page (no extra round-trip).
type gameDetailResponse struct {
	*models.GameWithDetails
	Ratings []playerRatingDelta `json:"ratings,omitempty"`
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

	resp := gameDetailResponse{GameWithDetails: game}

	// Attach rating deltas if the snapshot is available; ratings are best-effort —
	// a snapshot failure shouldn't break the game-detail load.
	if snap, err := h.rating.Snapshot(); err == nil {
		entries := snap.ForGame(game.Game.GameID)
		for _, gp := range game.GamePlayers {
			if e, ok := entries[gp.PlayerID]; ok {
				resp.Ratings = append(resp.Ratings, playerRatingDelta{
					PlayerID:     gp.PlayerID,
					Ordinal:      e.Ordinal,
					DeltaOrdinal: e.DeltaOrdinal,
				})
			}
		}
	} else {
		log.Printf("rating snapshot unavailable for game %d: %v", id, err)
	}

	h.sendJSON(w, http.StatusOK, resp)
}

// POST /games - Create a new game
func (h *Handler) createGame(w http.ResponseWriter, r *http.Request) {
	// Read the entire body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var auth RequestAuthentication
	if err := json.Unmarshal(bodyBytes, &auth); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	actor, err := h.repo.AuthenticatePlayer(auth.ActorName, auth.ActorPassword)
	if err != nil {
		log.Printf("Auth failure for actor %q: %v", auth.ActorName, err)
		h.sendError(w, http.StatusUnauthorized, "Invalid actor credentials")
		return
	}

	log.Printf("User %q: creating game", actor.Name)

	parsedReq, err := models.ParseGameRequest(bytes.NewReader(bodyBytes), false)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.ValidateGameRequest(parsedReq); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	images := parsedReq.GetImages()
	if len(images) > 0 {
		log.Printf("Processing %d image(s) for new game by %q", len(images), actor.Name)
		if err := h.processImages(images, nil); err != nil {
			h.sendError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	game, err := h.repo.CreateGame(parsedReq, *actor)
	if err != nil {
		log.Printf("Failed to create game: %v", err)
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create game: %v", err))
		return
	}
	h.rating.Invalidate()

	log.Printf("Game %d created by %q", game.Game.GameID, actor.Name)
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

	// Read the entire body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var auth RequestAuthentication
	if err := json.Unmarshal(bodyBytes, &auth); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	actor, err := h.repo.AuthenticatePlayer(auth.ActorName, auth.ActorPassword)
	if err != nil {
		log.Printf("Auth failure for actor %q: %v", auth.ActorName, err)
		h.sendError(w, http.StatusUnauthorized, "Invalid actor credentials")
		return
	}

	log.Printf("User %q: updating game %d", actor.Name, id)

	parsedReq, err := models.ParseGameRequest(bytes.NewReader(bodyBytes), true)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.ValidateGameUpdateRequest(id, parsedReq); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	images := parsedReq.GetImages()
	if len(images) > 0 {
		log.Printf("Processing %d image(s) for game %d update by %q", len(images), id, actor.Name)
		if err := h.processImages(images, &id); err != nil {
			h.sendError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	game, err := h.repo.UpdateGame(id, parsedReq, *actor)
	if err != nil {
		log.Printf("Failed to update game %d: %v", id, err)
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update game: %v", err))
		return
	}
	h.rating.Invalidate()

	log.Printf("Game %d updated by %q (revision %d)", game.Game.GameID, actor.Name, game.Game.Revision)
	h.sendJSON(w, http.StatusOK, game)
}

// DELETE /games/{id} - Delete a game
func (h *Handler) deleteGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var auth RequestAuthentication
	if err := json.Unmarshal(bodyBytes, &auth); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	actor, err := h.repo.AuthenticatePlayer(auth.ActorName, auth.ActorPassword)
	if err != nil {
		log.Printf("Auth failure for actor %q: %v", auth.ActorName, err)
		h.sendError(w, http.StatusUnauthorized, "Invalid actor credentials")
		return
	}

	log.Printf("User %q: deleting game %d", actor.Name, id)

	if err := h.repo.DeleteGame(id, *actor); err != nil {
		log.Printf("Failed to delete game %d: %v", id, err)
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete game: %v", err))
		return
	}
	h.rating.Invalidate()

	log.Printf("Game %d deleted by %q", id, actor.Name)
	w.WriteHeader(http.StatusNoContent)
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
	w.Write(imageData)
}

// processImages validates and processes all images in a slice
func (h *Handler) processImages(images []models.ImageRequest, gameID *int) error {
	if len(images) > 4 {
		return fmt.Errorf("too many images: maximum 4 images allowed, got %d", len(images))
	}

	// Process each image
	for i, imageReq := range images {
		if imageReq.ID != nil {
			// The image is a reference for an existing image
			_, _, err := h.repo.GetGameImageData(*imageReq.ID)
			if err != nil {
				return fmt.Errorf("referenced image %d not found", *imageReq.ID)
			}

			// We have to validate that the referenced image is for the same game
			if gameID != nil {
				imageGameID, err := h.repo.GetImageGameID(*imageReq.ID)
				if err != nil {
					return fmt.Errorf("could not verify ownership of image %d", *imageReq.ID)
				}
				if imageGameID != *gameID {
					return fmt.Errorf("image %d does not belong to this game", *imageReq.ID)
				}
			}
			continue
		}

		if len(imageReq.ImageData) == 0 || imageReq.MimeType == "" {
			return fmt.Errorf("image %d: missing image data or mime type", i+1)
		}

		processedData, finalMimeType, err := imageutil.ProcessImage(imageReq.ImageData, imageReq.MimeType)
		if err != nil {
			return fmt.Errorf("error processing image %d: %v", i+1, err)
		}

		images[i].ImageData = processedData
		images[i].MimeType = finalMimeType
	}

	return nil
}
