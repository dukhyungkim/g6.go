package main

import (
	"context"
	"github.com/dukhyungkim/gonuboard/util"
	"net/http"
	"strings"
)

func mainMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		request := r.Context().Value(KeyRequest).(util.Request)

		if !strings.HasPrefix(r.URL.Path, "/install") {
			if NeedInstall {
				renderAlertTemplate(w, request)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func renderAlertTemplate(w http.ResponseWriter, request util.Request) {
	tpl, err := util.AlertTemplate(request, ".env 파일이 없습니다. 설치를 진행해 주세요.", "/install")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(tpl)
}

type CtxKey string

const (
	KeyRequest CtxKey = "request"
)

func requestMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		request := util.NewRequest(r)
		ctx := context.WithValue(r.Context(), KeyRequest, request)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
