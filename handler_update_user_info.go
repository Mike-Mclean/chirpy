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

	userID, err := auth.ValidateJWT(bearer_token, cfg.secret)
	if err != nil {
		log.Println("error validating JWT:", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate")
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

	updatedInfo := database.UpdateUserInfoParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.UpdateUserInfo(r.Context(), updatedInfo)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not update user information")
		return
	}

	respBody := User {
		ID:			user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
	}

	respondWithJSON(w, http.StatusOK, respBody)

}