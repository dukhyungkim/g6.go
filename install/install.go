package install

import (
	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func DefaultRouter(r chi.Router) {
	r.Get("/", indexHandler())
}

func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := pongo2.FromFile("install/templates/main.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := pongo2.Context{}

		err = tpl.ExecuteWriter(data, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//_, _ = w.Write([]byte("install index"))
	}
}
