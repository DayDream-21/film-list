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
			slog.Error("id is missing in the request")

			http.Error(w, "id is missing in the request", http.StatusBadRequest)

			return
		}

		slog.Info("received delete data:", "id", id)

		_, err := filmDeleter.DeleteFilm(id)
		if err != nil {
			slog.Error("failed to delete film:", "error", err)

			http.Error(w, "failed to delete film", http.StatusInternalServerError)

			return
		}
	}
}
