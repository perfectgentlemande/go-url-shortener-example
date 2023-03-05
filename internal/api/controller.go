package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/perfectgentlemande/go-url-shortener-example/internal/service"
	"go.uber.org/fx"
)

type Config struct {
	AppPort string
	Domain  string
}

type Controller struct {
	srvc   *service.Service
	domain string
}

func New(srvc *service.Service, domain string) *Controller {
	return &Controller{
		srvc: srvc,
	}
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
func WriteError(ctx context.Context, w http.ResponseWriter, status int, message string) {
	err := RespondWithJSON(w, status, APIError{Message: message})
	if err != nil {
		log.Printf("cannot write response: %s\n", err)
	}
}
func WriteSuccessful(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	err := RespondWithJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("cannot write response: %s\n", err)
	}
}

func Provide(lifecycle fx.Lifecycle, apiConf *Config, srvc *service.Service) (*http.Server, error) {
	c := New(srvc, apiConf.Domain)
	r := chi.NewRouter()

	HandlerFromMux(c, r)
	srv := &http.Server{
		Handler: r,
		Addr:    apiConf.AppPort,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					log.Println("Listening on:", srv.Addr)
					err := srv.ListenAndServe()
					if err != nil {
						log.Println("Server listening error:", err)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		},
	)

	return srv, nil
}
