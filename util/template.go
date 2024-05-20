package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/model"
	"github.com/dukhyungkim/gonuboard/version"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

const (
	templates      = "templates"
	adminTemplates = "admin/templates"

	defaultTheme = "basic"
)

var (
	templatesDir      = getThemePath()
	adminTemplatesDir = getAdminThemePath()
)

func getCurrentTheme() string {
	dbConn := db.GetInstance()
	if dbConn == nil {
		return defaultTheme
	}

	var theme string
	if err := dbConn.Model(model.Config{}).Select("cf_theme").Scan(&theme).Error; err != nil {
		return defaultTheme
	}
	return theme
}

func getThemePath() string {
	const defaultThemePath = templates + "/" + defaultTheme
	theme := getCurrentTheme()
	themePath := templates + "/" + theme

	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return defaultThemePath
	}
	return themePath
}

func getAdminThemePath() string {
	const defaultThemePath = adminTemplates + "/" + defaultTheme

	theme := os.Getenv("ADMIN_THEME")

	if _, err := os.Stat(theme); os.IsNotExist(err) {
		return defaultThemePath
	}
	return theme
}

func themeAsset(r map[string]any, assetPath string) string {
	theme := getCurrentTheme()

	request := MapToStruct[Request](r)
	var mobileDir string
	if request.State.IsMobile {
		mobileDir = "/mobile"
	}
	fmt.Println(fmt.Sprintf("/theme_static/%s%s/%s", theme, mobileDir, assetPath))
	return fmt.Sprintf("/theme_static/%s%s/%s", theme, mobileDir, assetPath)
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

func RenderAlertTemplate(c *gin.Context, request Request, message string, statusCode int, redirect string) {
	data := exec.NewContext(map[string]any{
		"request": request.ToMap(),
		"errors":  []string{message},
		"url":     redirect,
	})

	c.HTML(statusCode, "templates/basic/alert.html", data)
}
