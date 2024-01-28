package util

import (
	"net/http"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func init() {
	gonja.DefaultContext.Set("theme_asset", themeAsset)
	gonja.DefaultContext.Set("url_for", urlFor)
}

func themeAsset(r map[string]any, assetPath string) string {
	return "templates/basic/static/" + assetPath
}
func urlFor(basePath string, assetPath string) string {
	return basePath + assetPath
}

func RenderTemplate(w http.ResponseWriter, path string, data *exec.Context) error {
	tpl, err := gonja.FromFile(path)
	if err != nil {
		return err
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
