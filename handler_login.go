package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		log.Println("error creating user:", err)
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		log.Println("error creating user:", err)
		respondWithError(w, http.StatusUnauthorized, "User not found")
		return
	}

	expiresIn := 3600

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Second * time.Duration(expiresIn))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem creating JWT")
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		log.Println("error creating refresh token:", err)
		respondWithError(w, http.StatusInternalServerError, "Problem creating refresh token")
		return
	}

	respBody := User {
		ID:			user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
		Token: 		token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, respBody)

}