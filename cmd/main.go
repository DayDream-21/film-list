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

// TODO: выложить на github
func main() {
	log := charmLogger.New()

	// TODO: парсить конфигурацию из файла, на основе конфигурации выбирать тип логгирования и ту БД и к которой мы подключаемся.
	//  Сделать два режима продакшен и разработки. В проде мы подключаемся к удаленной БД, а в разработке к локальной
	// TODO: добавить реализацию на PostgreSQL (может разбить на разные storage, типо sqlStorage и nosqlStorage?)
	storage, err := mongo.New(log)
	if err != nil {
		slog.Error("failed to create mongo client:", err)
	}

	http.HandleFunc("GET /", get.New(storage, log))
	http.HandleFunc("POST /add-film", save.New(storage, log))
	http.HandleFunc("DELETE /film", dbDelete.New(storage, log))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("failed to start server:", err)
	}
}
