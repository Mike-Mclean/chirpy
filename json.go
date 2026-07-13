package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string){
	type returnError struct {
		Err string `json:"error"`
	}

	respondWithJSON(w, code, returnError {
		Err: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any){
	resp, erro := json.Marshal(payload)
	if erro != nil {
		log.Printf("error marshalling JSON: %s", erro)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}