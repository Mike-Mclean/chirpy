package main

import (
	"net/http"
	"chirpy/internal/auth"
	"log"
	"time"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("error getting bearer token:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate")
		return
	}

	log.Println("bearertoken:", bearer_token)

	user_id, err := cfg.db.GetUserFromRefreshToken(r.Context(), bearer_token)
	if err != nil {
		log.Println("error looking up user from bearer token:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to find user")
		return
	}

	jwt, err := auth.MakeJWT(user_id, cfg.secret, time.Hour)
	if err != nil {
		log.Println("error creating jwt:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to make jwt")
		return
	}

	response := TokenResponse{
		Token: jwt,
	}

	respondWithJSON(w, http.StatusOK, response)

}