package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/omn1vor/blog-aggregator/internal/database"
)

type postsResponse struct {
	ID          uuid.UUID `json:"id"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cfg *apiConfig) getPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	var limit int
	limitParameter := r.URL.Query().Get("limit")
	if limitParameter != "" {
		lim, err := strconv.ParseInt(limitParameter, 10, 32)
		if err != nil {
			log.Print("Could not parse posts limit from url parameter", limitParameter)
		} else {
			limit = int(lim)
		}

	}
	if limit == 0 {
		limit = cfg.PostsLimit
	}
	params := database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := cfg.DB.GetPostsByUser(r.Context(), params)
	if checkErrorAndRespond(err, w, http.StatusInternalServerError, "Error while getting posts") {
		return
	}

	response := []postsResponse{}
	for _, post := range posts {
		response = append(response, postsResponse{
			ID:          post.ID,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description.String,
			PublishedAt: post.PublishedAt.Time.UTC(),
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			FeedID:      post.FeedID,
		})
	}

	respondWithJson(w, http.StatusOK, response)
}
