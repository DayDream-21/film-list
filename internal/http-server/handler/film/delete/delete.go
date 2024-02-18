package delete

import (
	"log/slog"
	"net/http"
)

type FilmDeleter interface {
	DeleteFilm(id string) (int64, error)
}

func New(filmDeleter FilmDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			slog.Warn("id is missing in the request")
			return
		}

		slog.Info("Received delete data:", "id", id)

		filmDeleter.DeleteFilm(id)
	}
}
