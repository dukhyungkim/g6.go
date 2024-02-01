package util

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"net/http"
	"sync"
)

func init() {
	defaultCtx := gonja.DefaultContext
	defaultCtx.Set("default_version", "Go누보드6.0.0")
	defaultCtx.Set("theme_asset", themeAsset)
	defaultCtx.Set("url_for", urlFor)
}

func themeAsset(r map[string]any, assetPath string) string {
	return "templates/basic/static/" + assetPath
}

var UrlMap = sync.Map{}

func urlFor(assetPath string) string {
	value, ok := UrlMap.Load(assetPath)
	if !ok {
		return ""
	}
	return value.(string)
}

func RenderTemplate(w http.ResponseWriter, path string, data *exec.Context) {
	tpl, err := gonja.FromFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
