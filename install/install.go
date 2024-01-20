package install

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("install index"))
}
