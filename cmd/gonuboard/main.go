package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dukhyungkim/gonuboard/config"
	"github.com/dukhyungkim/gonuboard/db"
	"github.com/dukhyungkim/gonuboard/install"
	mw "github.com/dukhyungkim/gonuboard/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"

	"github.com/dukhyungkim/gonuboard/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	err = Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func Run() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	staticServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static", staticServer))
	templatesServer := http.FileServer(http.Dir("templates"))
	r.Handle("/templates/*", http.StripPrefix("/templates", templatesServer))

	r.Group(func(r chi.Router) {
		r.Use(mw.RequestMiddleware)
		r.Use(mw.MainMiddleware)
		r.Use(mw.UrlForMiddleware)

		r.Get("/", defaultHandler)
		r.Post("/generate_token", generateToken)
		r.Route("/install", install.DefaultRouter)
	})

	addr := ":8080"
	fmt.Printf("running on %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		return err
	}
	return nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	const templatePath = "templates/basic/index.html"
	request := r.Context().Value(mw.KeyRequest)
	data := exec.NewContext(map[string]any{
		"request": request.(util.Request).ToMap(),
	})

	util.RenderTemplate(w, templatePath, data)
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

func generateToken(w http.ResponseWriter, r *http.Request) {
	tokenHex, err := util.TokenHex(16)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = render.Render(w, r, NewTokenResponse(tokenHex))
}
