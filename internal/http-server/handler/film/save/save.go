package save

import (
	"film-list/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

type FilmSaver interface {
	SaveFilm(film dto.Film) (string, error)
}

func New(filmSaver FilmSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Just to see beautiful spinner
		time.Sleep(300 * time.Millisecond)

		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		if title == "" || director == "" {
			slog.Error("title or director is missing in the request")

			http.Error(w, "title or director is missing in the request", http.StatusBadRequest)

			return
		}

		slog.Info("received post data:", "title", title, "director", director)

		newFilm := dto.Film{
			Title:    title,
			Director: director,
		}

		idHex, err := filmSaver.SaveFilm(newFilm)
		if err != nil {
			slog.Error("failed to save film:", "error", err)

			http.Error(w, "failed to save film", http.StatusInternalServerError)

			return
		}

		tmpl := template.Must(template.ParseFiles("index.html"))

		newFilm.ID, err = primitive.ObjectIDFromHex(idHex)
		if err != nil {
			slog.Error("failed to convert id to ObjectID:", "error", err)

			http.Error(w, "failed to convert id to ObjectID", http.StatusInternalServerError)

			return
		}

		err = tmpl.ExecuteTemplate(w, "film-list-element", newFilm)
		if err != nil {
			slog.Error("failed to execute template:", "error", err)

			http.Error(w, "failed to execute template", http.StatusInternalServerError)

			return
		}
	}
}
