package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func resolve(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("resolve"))
}
func shorten(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shorten"))
}

func main() {
	r := chi.NewRouter()
	r.Get("/:url", resolve)
	r.Post("/api/v1", shorten)
	http.ListenAndServe(":3000", r)
}
