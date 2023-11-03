package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/omn1vor/blog-aggregator/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Name string `json:"name"`
	}
	var params userRequest

	err := json.NewDecoder(r.Body).Decode(&params)
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while decoding request body") {
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while creating user") {
		return
	}

	userResponse := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Name      string    `json:"name"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	}

	respondWithJson(w, http.StatusOK, userResponse)
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Expecting an API key")
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), strings.TrimPrefix(authHeader, "ApiKey "))
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while getting user") {
		return
	}

	userResponse := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}

	respondWithJson(w, http.StatusOK, userResponse)
}
