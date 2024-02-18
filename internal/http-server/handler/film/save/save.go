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
			slog.Warn("title or director is missing in the request")

			w.WriteHeader(http.StatusBadRequest)

			return
		}

		slog.Info("Received post data:", "title", title, "director", director)

		newFilm := dto.Film{
			Title:    title,
			Director: director,
		}

		idHex, err := filmSaver.SaveFilm(newFilm)
		if err != nil {
			return
		}

		tmpl := template.Must(template.ParseFiles("index.html"))

		newFilm.ID, _ = primitive.ObjectIDFromHex(idHex)

		tmpl.ExecuteTemplate(w, "film-list-element", newFilm)
	}
}
