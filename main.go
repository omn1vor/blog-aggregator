package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/omn1vor/blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB                  *database.Queries
	NumberOfFeedsToRead int
	ReadingInterval     time.Duration
	PostsLimit          int
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_CON_STRING")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error while connection with DB: ", err)
	}

	dbQueries := database.New(db)
	apiConfig := apiConfig{
		DB:                  dbQueries,
		NumberOfFeedsToRead: 10,
		ReadingInterval:     60 * time.Second,
		PostsLimit:          20,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{}))

	v1 := chi.NewRouter()
	v1.Get("/readiness", readinessHandler)
	v1.Get("/err", errorHandler)
	v1.Post("/users", apiConfig.createUser)
	v1.Get("/users", apiConfig.middlewareAuth(apiConfig.getUser))
	v1.Post("/feeds", apiConfig.middlewareAuth(apiConfig.createFeed))
	v1.Get("/feeds", apiConfig.getFeeds)
	v1.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.createFeedFollow))
	v1.Delete("/feed_follows/{id}", apiConfig.middlewareAuth(apiConfig.deleteFeedFollow))
	v1.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.getFeedFollows))
	v1.Get("/posts", apiConfig.middlewareAuth(apiConfig.getPosts))

	router.Mount("/v1", v1)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Starting reading")
	StartReadingFeeds(ctx, apiConfig)

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
