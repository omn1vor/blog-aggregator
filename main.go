package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{}))

	v1 := chi.NewRouter()
	v1.Get("/readiness", readinessHandler)
	v1.Get("/err", errorHandler)

	router.Mount("/v1", v1)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	respondWithJson(w, http.StatusOK, response)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

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
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't encode to JSON: "+err.Error())
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
