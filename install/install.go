package install

import (
	"net/http"

	"github.com/dukhyungkim/gonuboard/util"
	"github.com/dukhyungkim/gonuboard/version"
	"github.com/go-chi/chi/v5"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func DefaultRouter(r chi.Router) {
	r.Get("/", indexHandler())
}

func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const templatePath = "install/templates/main.html"
		data := exec.NewContext(map[string]any{
			"default_version": version.Version,
			"python_version":  version.RuntimeVersion,
			"fastapi_version": version.RouterVersion,
		})

		err := util.RenderTemplate(w, templatePath, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
