package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/omn1vor/blog-aggregator/internal/database"
)

type feedFollowResponse struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	params := struct {
		FeedID uuid.UUID `json:"feed_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while decoding request body") {
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while creating feed follow") {
		return
	}

	response := feedFollowResponse{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	respondWithJson(w, http.StatusOK, response)
}

func (cfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID, err := uuid.Parse(chi.URLParam(r, "id"))
	if checkErrorAndRespond(err, w, http.StatusBadRequest, "Wrong ID") {
		return
	}

	_, err = cfg.DB.GetFeedFollow(r.Context(), database.GetFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "ID not found")
			return
		} else if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while getting feed follow") {
			return
		}
	}

	err = cfg.DB.DeleteFeedFollows(r.Context(), feedFollowID)
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while deleting feed follow") {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while getting feed follow") {
		return
	}

	response := []feedFollowResponse{}
	for _, follow := range follows {
		response = append(response, feedFollowResponse{
			ID:        follow.ID,
			CreatedAt: follow.CreatedAt,
			UpdatedAt: follow.UpdatedAt,
			UserID:    follow.UserID,
			FeedID:    follow.FeedID,
		})
	}

	respondWithJson(w, http.StatusOK, response)
}
