package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/murtazapatel89100/RSS-Aggregartor/internal/database"
	"github.com/murtazapatel89100/RSS-Aggregartor/internal/handler"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the env")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("No DB URL found in the env")
	}

	conection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	config := handler.ApiConfig{DB: database.New(conection)}

	router := chi.NewRouter()

	router.Use(cors.Handler((cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})))

	v1router := chi.NewRouter()
	v1router.Get("/health", handler.HandlerReadiness)
	v1router.Get("/error", handler.HandlerError)
	v1router.Post("/users/create", config.HandlerCreateUser)
	v1router.Post("/feeds/fetch", config.HandlerGetFeeds)

	v1router.With(config.MiddlewareAuth).Get("/users/fetch", config.HandlerGetUser)
	v1router.With(config.MiddlewareAuth).Post("/feeds/create", config.HandlerCreateFeed)
	v1router.With(config.MiddlewareAuth).Post("/feeds-follow/create", config.HandlerCreateFeedFollow)

	router.Mount("/v1", v1router)

	fmt.Printf("Server is atrting on port %v", portString)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
