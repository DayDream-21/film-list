package main

import (
	"film-list/internal/config"
	dbDelete "film-list/internal/http-server/handler/film/delete"
	"film-list/internal/http-server/handler/film/get"
	"film-list/internal/http-server/handler/film/save"
	"film-list/internal/logger"
	"film-list/internal/storage/mongo"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	// TODO: парсить конфигурацию из файла, на основе конфигурации выбирать ту БД и к которой мы подключаемся.
	//  В проде мы подключаемся к удаленной БД, а в разработке к локальной
	// TODO: добавить реализацию на PostgreSQL (может разбить на разные storage, типо sqlStorage и nosqlStorage?)
	storage, err := mongo.New(log)
	if err != nil {
		slog.Error("failed to create mongo client:", err)
	}

	http.HandleFunc("GET /", get.New(storage, log))
	http.HandleFunc("POST /add-film", save.New(storage, log))
	http.HandleFunc("DELETE /film", dbDelete.New(storage, log))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      nil,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = server.ListenAndServe(); err != nil {
		slog.Error("failed to start server:", err)
	}
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = charmLogger.New()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		fallthrough
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
