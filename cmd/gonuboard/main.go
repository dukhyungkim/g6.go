package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dukhyungkim/gonuboard/config"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/install"
	"github.com/dukhyungkim/gonuboard/lib"
	mw "github.com/dukhyungkim/gonuboard/middleware"
	"github.com/dukhyungkim/gonuboard/util"
	"github.com/gin-gonic/gin"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	parseFlags()

	if FlagVersion {
		printVersion()
		return
	}

	if FlagHelp {
		flag.Usage()
		return
	}

	err := config.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalln(err)
	}

	engine := config.Global.DbEngine
	_, err = db.NewDB(engine)
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.HTMLRender = util.NewTemplateEngine()
	if err = Run(r); err != nil {
		log.Fatal(err)
	}
}

func Run(r *gin.Engine) error {
	r.Static("/static", "static")
	r.Static("/templates", "templates")

	r.Group("/")

	r.Use(mw.RequestMiddleware())
	r.Use(mw.MainMiddleware())
	r.Use(mw.UrlForMiddleware())

	r.GET("/", defaultHandler)
	r.POST("/generate_token", generateToken)

	install.DefaultRouter(r)

	addr := ":8080"
	fmt.Printf("running on %s\n", addr)
	return r.Run(addr)
}

func defaultHandler(c *gin.Context) {
	const templatePath = "templates/basic/index.html"
	request := c.MustGet(mw.KeyRequest).(util.Request)
	data := exec.NewContext(map[string]interface{}{
		"request": request.ToMap(),
	})

	c.HTML(http.StatusOK, templatePath, data)
}

type TokenResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func NewTokenResponse(token string) TokenResponse {
	return TokenResponse{
		Success: true,
		Token:   token,
	}
}

func generateToken(c *gin.Context) {
	tokenHex, err := util.TokenHex(16)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, lib.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, NewTokenResponse(tokenHex))
}
