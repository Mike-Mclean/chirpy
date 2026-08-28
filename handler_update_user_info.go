package main

import (
	"net/http"
	"chirpy/internal/auth"
	"log"
	"encoding/json"
	"chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	bearer_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("error getting bearer token:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate")
		return
	}

	userID, err = cfg.db.GetUserFromRefreshToken(r.Context(), bearer_token)
	if err != nil {
		log.Println("error looking up user from bearer token:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to find user")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err = decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not hash password")
		return
	}

	updatedUserParams := database.CreateUserParams {
		Email: params.Email,
		HashedPassword: hashedPassword,
	}

	

}