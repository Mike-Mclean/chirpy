package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type userParams struct {
	Email 				string 	`json:"email"`
	Password 			string 	`json:"password"`
}

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not hash password")
		return
	}

	newUserParams := database.CreateUserParams {
		Email: params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		log.Println("error creating user:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respBody := User {
		ID:			user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
	}

	respondWithJSON(w, http.StatusCreated, respBody)
}

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

	respBody := User {
		ID:			user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
		Token: 		token,
		RefreshToken: auth.MakeRefreshToken(),
	}

	respondWithJSON(w, http.StatusOK, respBody)

}