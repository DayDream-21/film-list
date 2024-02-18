package get

import (
	"film-list/internal/dto"
	"html/template"
	"net/http"
)

type FilmGetter interface {
	GetFilms() ([]dto.Film, error)
}

func New(filmGetter FilmGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filmSlice []dto.Film

		filmSlice, err := filmGetter.GetFilms()
		if err != nil {
			return
		}

		films := make(map[string][]dto.Film)

		films["Films"] = filmSlice

		tmpl := template.Must(template.ParseFiles("index.html"))

		tmpl.Execute(w, films)
	}
}
