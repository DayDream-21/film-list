package main

import (
	dbDelete "film-list/internal/http-server/handler/film/delete"
	"film-list/internal/http-server/handler/film/get"
	"film-list/internal/http-server/handler/film/save"
	"film-list/internal/storage/myMongo"

	"log"
	"log/slog"
	"net/http"
)

// TODO: доработать обработку ошибок, добавить логирование, добавить информативности при возникновении ошибок
// TODO: выложить на github
func main() {
	storage, err := myMongo.New()
	if err != nil {
		slog.Error("Failed to create mongo client:", err)
	}

	http.HandleFunc("GET /", get.New(storage))
	http.HandleFunc("POST /add-film", save.New(storage))
	http.HandleFunc("DELETE /film", dbDelete.New(storage))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
