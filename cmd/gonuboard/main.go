package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dukhyungkim/gonuboard/install"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var (
	NeedInstall = false
)

func main() {
	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			NeedInstall = true
		} else {
			log.Fatalln(err)
		}
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", defaultHandler)
	r.Route("/install", install.DefaultRouter)

	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	addr := ":8080"
	fmt.Printf("running on %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	if NeedInstall {
		http.Redirect(w, r, "/install", http.StatusMovedPermanently)
		return
	}

	_, _ = w.Write([]byte("hello gnuboard"))
}
