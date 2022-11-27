package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/:url", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("resolve"))
	})
	r.Post("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("shorten"))
	})
	http.ListenAndServe(":3000", r)
}
