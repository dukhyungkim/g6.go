package util

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/dukhyungkim/gonuboard/version"
	"github.com/labstack/echo/v4"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

var UserTemplate = newUserTemplateProcessor()
var AdminTemplate = newAdminTemplateProcessor()

type TemplateProcessor struct {
	ctx              *exec.Context
	contextProcessor func(request *Request) *exec.Context
}

func newUserTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{}
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

type TemplateRenderer struct {
	defaultCtx *exec.Context
	templates  *exec.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	defaultCtx := exec.NewContext(map[string]interface{}{
		"default_version": version.Version,
		"theme_asset":     themeAsset,
		"url_for":         urlFor,
	})
	return &TemplateRenderer{defaultCtx: defaultCtx}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tpl, err := gonja.FromFile(name)
	if err != nil {
		log.Println(err)
		return err
	}

	ctxData := t.defaultCtx
	if d, ok := data.(*exec.Context); ok {
		ctxData.Update(d)
	}

	err = tpl.Execute(w, ctxData)
	if err != nil {
		log.Println(err)
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
