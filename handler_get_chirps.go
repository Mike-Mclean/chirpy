package main

import (
	"log"
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	allChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Println("error retrieving chirps:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp")
		return
	}

	var chirps []Chirp
	for _, chirp := range allChirps {
		i := Chirp{
			ID: 		chirp.ID,
			CreatedAt: 	chirp.CreatedAt,
			UpdatedAt: 	chirp.UpdatedAt,
			Body: 		chirp.Body,
			UserID: 	chirp.UserID,
		}

		chirps = append(chirps, i)

	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirp_id, err := uuid.Parse(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to find chirp")
		return
	}

	respBody := Chirp {
		ID: 		chirp.ID,
		CreatedAt: 	chirp.CreatedAt,
		UpdatedAt: 	chirp.UpdatedAt,
		Body: 		chirp.Body,
		UserID: 	chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, respBody)
}