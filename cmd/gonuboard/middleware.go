package main

import "net/http"

func mainMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//currentIP := r.
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
