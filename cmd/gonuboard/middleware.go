package main

import (
	"context"
	"github.com/dukhyungkim/gonuboard/util"
	"net/http"
)

func mainMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//currentIP := r.
		request := util.NewRequest(r)
		ctx := context.WithValue(r.Context(), "request", request)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
