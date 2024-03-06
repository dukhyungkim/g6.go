package lib

import "net/http"

func SendHTML(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(html))
}
