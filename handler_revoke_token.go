package main

import (
	"net/http"
	"log"
	"chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("error getting bearer token:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), bearer_token)
	if err != nil {
		log.Println("error revokeing refresh token:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}