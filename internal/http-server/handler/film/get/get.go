package get

import (
	"film-list/internal/dto"
	"github.com/charmbracelet/log"
	"html/template"
	"net/http"
)

type FilmGetter interface {
	GetFilms() ([]dto.Film, error)
}

func New(filmGetter FilmGetter, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filmSlice []dto.Film

		filmSlice, err := filmGetter.GetFilms()
		if err != nil {
			log.Error("failed to get films:", err)

			http.Error(w, "failed to get films", http.StatusInternalServerError)

			return
		}

		films := make(map[string][]dto.Film)

		films["Films"] = filmSlice

		tmpl := template.Must(template.ParseFiles("index.html"))

		err = tmpl.Execute(w, films)
		if err != nil {
			log.Error("failed to execute template:", err)

			http.Error(w, "failed to execute template", http.StatusInternalServerError)

			return
		}
	}
}
