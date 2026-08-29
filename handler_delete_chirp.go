package main

import (
	"net/http"
	"log"
	"chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("error getting bearer token:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate")
		return
	}

	userID, err := auth.ValidateJWT(bearer_token, cfg.secret)
	if err != nil {
		log.Println("error validating JWT:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate")
		return
	}

	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not retrieve chirpID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find chirp from ID")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Unauthorized User")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}