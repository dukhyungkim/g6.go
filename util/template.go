package util

import (
	"net/http"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func init() {
	gonja.DefaultContext.Set("theme_asset", themeAsset)
}

func themeAsset(r map[string]any, assetPath string) string {
	return "templates/basic/static/" + assetPath
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

func AlertTemplate(req Request, message string, redirect string) ([]byte, error) {
	tpl, err := gonja.FromFile("templates/basic/alert.html")
	if err != nil {
		return nil, err
	}
	data := exec.NewContext(map[string]any{
		"request": req.ToMap(),
		"errors":  []string{message},
		"url":     redirect,
	})

	bytes, err := tpl.ExecuteToBytes(data)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
