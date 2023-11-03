package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.Encode(struct {
		Error string `json:"error"`
	}{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Can't encode to JSON") {
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func checkErrorAndRespond(err error, w http.ResponseWriter, code int, msg string) bool {
	if err == nil {
		return false
	}
	respondWithError(w, code, msg)
	log.Println(err.Error())
	return true
}
