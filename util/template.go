package util

import (
	"log"
	"net/http"
	"sync"

	"github.com/dukhyungkim/gonuboard/version"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

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

type TemplateEngine struct {
	defaultCtx *exec.Context
}

func NewTemplateEngine() *TemplateEngine {
	defaultCtx := exec.NewContext(map[string]interface{}{
		"default_version": version.Version,
		"theme_asset":     themeAsset,
		"url_for":         urlFor,
	})
	return &TemplateEngine{defaultCtx: defaultCtx}
}

func (t *TemplateEngine) Instance(name string, data any) render.Render {
	template, err := gonja.FromFile(name)
	if err != nil {
		panic(err)
	}

	templateData := t.defaultCtx
	if d, ok := data.(*exec.Context); ok {
		templateData.Update(d)
	}

	return Renderer{
		template: template,
		data:     templateData,
	}
}

type Renderer struct {
	template *exec.Template
	data     *exec.Context
}

func (r Renderer) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	err := r.template.Execute(w, r.data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r Renderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

func RenderAlertTemplate(c *gin.Context, request Request, message string, statusCode int, redirect string) {
	data := exec.NewContext(map[string]any{
		"request": request.ToMap(),
		"errors":  []string{message},
		"url":     redirect,
	})

	c.HTML(statusCode, "templates/basic/alert.html", data)
}
