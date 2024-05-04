package util

import (
	"github.com/dukhyungkim/gonuboard/lib"
	"github.com/dukhyungkim/gonuboard/version"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"log"
	"net/http"
	"sync"
)

var UserTemplate = newUserTemplateProcessor()
var AdminTemplate = newAdminTemplateProcessor()

type TemplateProcessor struct {
	ctx              *exec.Context
	contextProcessor func(request *Request) *exec.Context
}

func newUserTemplateProcessor() *TemplateProcessor {
	defaultCtx := gonja.DefaultContext
	defaultCtx.Set("default_version", version.Version)
	defaultCtx.Set("theme_asset", themeAsset)
	defaultCtx.Set("url_for", urlFor)

	return &TemplateProcessor{ctx: defaultCtx}
}

func newAdminTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{}
}

func themeAsset(_ map[string]any, assetPath string) string {
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

func processDefaultContext(request *Request) *exec.Context {
	return exec.NewContext(map[string]interface{}{
		"current_login_count": lib.GetCurrentLoginCount(request),
	})
}

func RenderTemplate(w http.ResponseWriter, path string, data *exec.Context) {
	tpl, err := gonja.FromFile(path)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, data)
	if err != nil {
		log.Println(err)
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

func RenderAlertTemplate(w http.ResponseWriter, request Request, message string, statusCode int, url string) {
	tpl, err := AlertTemplate(request, message, url)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(tpl)
}
