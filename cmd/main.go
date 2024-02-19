package main

import (
	dbDelete "film-list/internal/http-server/handler/film/delete"
	"film-list/internal/http-server/handler/film/get"
	"film-list/internal/http-server/handler/film/save"
	"film-list/internal/logger"
	"film-list/internal/storage/mongo"
	"log/slog"
	"net/http"
)

// TODO: добавить красивое логирование
// TODO: выложить на github
func main() {
	log := logger.New()

	// TODO: добавить реализацию на PostgreSQL (может разбить на разные storage, типо sqlStorage и nosqlStorage?)
	storage, err := mongo.New(log)
	if err != nil {
		slog.Error("failed to create mongo client:", err)
	}

	http.HandleFunc("GET /", get.New(storage, log))
	http.HandleFunc("POST /add-film", save.New(storage, log))
	http.HandleFunc("DELETE /film", dbDelete.New(storage, log))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
