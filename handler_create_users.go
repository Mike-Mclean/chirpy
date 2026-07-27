package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	type email struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := email{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
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