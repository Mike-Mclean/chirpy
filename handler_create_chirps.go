package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"
	"chirpy/internal/database"
	"chirpy/internal/auth"
)


func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
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

	type chirpRequest struct {
		Body 	string 			`json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirpRequest{}
	err = decoder.Decode(&params)
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
		UserID: userID,
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
		UserID: 	userID,
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