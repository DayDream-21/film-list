package delete

import (
	"github.com/charmbracelet/log"
	"net/http"
)

type FilmDeleter interface {
	DeleteFilm(id string) (int64, error)
}

func New(filmDeleter FilmDeleter, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			log.Error("id is missing in the request")

			http.Error(w, "id is missing in the request", http.StatusBadRequest)

			return
		}

		log.Info("received delete data:", "id", id)

		_, err := filmDeleter.DeleteFilm(id)
		if err != nil {
			log.Error("failed to delete film:", "error", err)

			http.Error(w, "failed to delete film", http.StatusInternalServerError)

			return
		}
	}
}
