package middleware

import (
	"context"
	"github.com/dukhyungkim/gonuboard/config"
	"github.com/dukhyungkim/gonuboard/util"
	"net/http"
	"net/url"
	"strings"
)

func MainMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !shouldRunMiddleware(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		request := r.Context().Value(KeyRequest).(util.Request)

		//engine := os.Getenv("DB_ENGINE")
		//dbConn, err := db.NewDB(engine)
		//if err != nil {
		//	fmt.Println("TODO print new db error")
		//	return
		//}

		if !strings.HasPrefix(r.URL.Path, "/install") {
			if config.NotExistENV {
				renderAlertTemplate(w, request, ".env 파일이 없습니다. 설치를 진행해 주세요.", http.StatusBadRequest, "/install")
				return
			}

			// TODO check config table
		} else {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func shouldRunMiddleware(path string) bool {
	switch path {
	case "/generate_token":
		return false
	}
	return true
}

func renderAlertTemplate(w http.ResponseWriter, request util.Request, message string, statusCode int, url string) {
	tpl, err := util.AlertTemplate(request, message, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(tpl)
}

type CtxKey string

const (
	KeyRequest CtxKey = "request"
)

func RequestMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		request := util.NewRequest(r)
		ctx := context.WithValue(r.Context(), KeyRequest, request)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func UrlForMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		setUrlMapForInstall(r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func setUrlMapForInstall(r *http.Request) {
	request := r.Context().Value(KeyRequest).(util.Request)
	installURL, _ := url.JoinPath(request.BaseURL, "install")
	util.UrlMap.Store("install_license", installURL+"/license")
	util.UrlMap.Store("install_form", installURL+"/form")
	util.UrlMap.Store("install", installURL)
}
