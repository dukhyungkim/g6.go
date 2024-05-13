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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e := echo.New()
	e.Renderer = util.NewTemplateRenderer()
	e.Logger.Fatal(Run(e))
}

func Run(e *echo.Echo) error {
	lib.NewLogger()
	e.Use(mw.LoggingMiddleware)
	e.Use(middleware.Recover())

	e.Static("/static/*", "static")
	e.Static("/templates/*", "templates")

	g := e.Group("/")

	g.Use(mw.RequestMiddleware)
	g.Use(mw.MainMiddleware)
	g.Use(mw.UrlForMiddleware)

	g.GET("/", defaultHandler)
	g.POST("/generate_token", generateToken)

	install.DefaultRouter(e)

	addr := ":8080"
	fmt.Printf("running on %s\n", addr)
	return e.Start(addr)
}

func defaultHandler(c echo.Context) error {
	const templatePath = "templates/basic/index.html"
	request := c.Get(mw.KeyRequest).(util.Request)
	data := exec.NewContext(map[string]interface{}{
		"request": request.ToMap(),
	})

	return c.Render(http.StatusOK, templatePath, data)
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

func (t TokenResponse) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func generateToken(c echo.Context) error {
	tokenHex, err := util.TokenHex(16)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, NewTokenResponse(tokenHex))
}
