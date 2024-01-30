package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dukhyungkim/gonuboard/install"
	"net/http"
	"os"

	"github.com/dukhyungkim/gonuboard/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/nikolalohinski/gonja/v2/exec"
)

var (
	NeedInstall = false
)

func main() {
	parseFlags()

	if FlagVersion {
		printVersion()
		return
	}

	if FlagHelp {
		flag.Usage()
		return
	}

	err := loadEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		r.Use(requestMiddleware)
		r.Use(mainMiddleware)

		r.Get("/", defaultHandler)
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

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			NeedInstall = true
		} else {
			return err
		}
	}
	return nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	const templatePath = "templates/basic/index.html"
	request := r.Context().Value("request")
	data := exec.NewContext(map[string]any{
		"request": request.(util.Request).ToMap(),
	})

	util.RenderTemplate(w, templatePath, data)
}
