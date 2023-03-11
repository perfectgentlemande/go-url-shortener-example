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

type loggerCtxKey struct{}

func WithLogger(ctx context.Context, log service.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, log)
}
func GetLogger(ctx context.Context) service.Logger {
	// no checks because Provide woill not run without logger
	return ctx.Value(loggerCtxKey{}).(service.Logger)
}

func newLoggingMiddleware(log service.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextLog := log.WithFields(
				map[string]interface{}{
					"method": r.Method,
					"path":   r.URL.Path,
				})
			handler.ServeHTTP(w, r.WithContext(WithLogger(r.Context(), nextLog)))
		})
	}
}

func Provide(lifecycle fx.Lifecycle, apiConf *Config, srvc *service.Service, log service.Logger) (*http.Server, error) {
	c := New(srvc, apiConf.Domain)
	r := chi.NewRouter()
	r.Use(newLoggingMiddleware(log))

	HandlerFromMux(c, r)
	srv := &http.Server{
		Handler: r,
		Addr:    apiConf.AppPort,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					log.WithField("addr", srv.Addr).Info("listening")
					err := srv.ListenAndServe()
					if err != nil {
						log.WithError(err).Error("server listening error")
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
