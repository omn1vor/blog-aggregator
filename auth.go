package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/omn1vor/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey"))
		if apiKey == "" {
			respondWithError(w, http.StatusUnauthorized, "Expecting an API key")
			return
		}

		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "User not found")
				return
			} else if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while getting user") {
				return
			}
		}
		handler(w, r, user)
	})

}
