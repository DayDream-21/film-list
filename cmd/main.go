package main

import (
	dbDelete "film-list/internal/http-server/handler/film/delete"
	"film-list/internal/http-server/handler/film/get"
	"film-list/internal/http-server/handler/film/save"
	"film-list/internal/storage/mongo"

	"log"
	"log/slog"
	"net/http"
)

// TODO: добавить красивое логирование
// TODO: выложить на github
func main() {
	// TODO: добавить реализацию на PostgreSQL
	storage, err := mongo.New()
	if err != nil {
		slog.Error("failed to create mongo client:", err)
	}

	http.HandleFunc("GET /", get.New(storage))
	http.HandleFunc("POST /add-film", save.New(storage))
	http.HandleFunc("DELETE /film", dbDelete.New(storage))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
