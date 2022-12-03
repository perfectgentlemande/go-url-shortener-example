package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/logger"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/storage/dburl"
)

func resolve(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("resolve"))
}
func shorten(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shorten"))
}

func main() {
	ctx := context.Background()
	log := logger.DefaultLogger()

	dbURL := dburl.NewDatabase(&dburl.Config{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer dbURL.Close()

	err := dbURL.Ping(ctx)
	if err != nil {
		log.WithError(err).Error("cannot ping dbURL")
		return
	}

	r := chi.NewRouter()
	r.Get("/:url", resolve)
	r.Post("/api/v1", shorten)
	http.ListenAndServe(":3000", r)
}
