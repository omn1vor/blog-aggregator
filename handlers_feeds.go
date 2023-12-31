package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/omn1vor/blog-aggregator/internal/database"
)

type feedResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type feedRequest struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	var params feedRequest

	err := json.NewDecoder(r.Body).Decode(&params)
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while decoding request body") {
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while creating feed") {
		return
	}

	follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while creating feed follow") {
		return
	}

	response := struct {
		Feed       feedResponse       `json:"feed"`
		FeedFollow feedFollowResponse `json:"feed_follow"`
	}{
		Feed: feedResponse{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
		},
		FeedFollow: feedFollowResponse{
			ID:        follow.ID,
			UserID:    follow.UserID,
			FeedID:    follow.FeedID,
			CreatedAt: follow.CreatedAt,
			UpdatedAt: follow.UpdatedAt,
		},
	}

	respondWithJson(w, http.StatusOK, response)
}

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while getting feeds") {
		return
	}

	response := []feedResponse{}
	for _, feed := range feeds {
		response = append(response, feedResponse{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, response)
}
