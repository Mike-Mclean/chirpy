package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
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
