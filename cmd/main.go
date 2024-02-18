package main

import (
	dbDelete "film-list/internal/http-server/handler/film/delete"
	"film-list/internal/http-server/handler/film/get"
	"film-list/internal/http-server/handler/film/save"
	"film-list/internal/logger"
	"film-list/internal/storage/mongo"

	"log"
	"log/slog"
	"net/http"
)

// TODO: добавить красивое логирование
// TODO: выложить на github
func main() {
	logger := logger.New()

	// TODO: добавить реализацию на PostgreSQL
	storage, err := mongo.New(logger)
	if err != nil {
		slog.Error("failed to create mongo client:", err)
	}

	http.HandleFunc("GET /", get.New(storage, logger))
	http.HandleFunc("POST /add-film", save.New(storage, logger))
	http.HandleFunc("DELETE /film", dbDelete.New(storage, logger))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
