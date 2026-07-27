package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"
	"github.com/google/uuid"
	"chirpy/internal/database"
)


func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirpRequest{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 500, "Something went wrong")
		return
	}

	const maxChirpLen = 140
	if len(params.Body) > maxChirpLen {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	newChirpParams := database.CreateChirpParams{
		Body: removeProfane(params.Body),
		UserID: params.UserID,
	}

	newChirp, err := cfg.db.CreateChirp(r.Context(), newChirpParams)
	if err != nil {
		log.Println("error creating chirp:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	respBody := Chirp {
		ID: 		newChirp.ID,
		CreatedAt: 	newChirp.CreatedAt,
		UpdatedAt: 	newChirp.UpdatedAt,
		Body: 		newChirp.Body,
		UserID: 	newChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, respBody)

}

func removeProfane (msg string) string {
	msgDetails := strings.Split(msg, " ")

	profane := map[string]bool {
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	for i := 0; i < len(msgDetails); i++ {
		word := strings.ToLower(msgDetails[i])
		if profane[word] {
			msgDetails[i] = "****"
		}
	}
	cleanMsg := strings.Join(msgDetails, " ")
	return cleanMsg
}