package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"log"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) printMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(message))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete users")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Content string `json:"body"`
	}
	type returnValid struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 500, "Something went wrong")
		return
	}

	const maxChirpLen = 140
	if len(params.Content) > maxChirpLen {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, 200, returnValid{
		Cleaned_body: removeProfane(params.Content),
	})
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
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to database")
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