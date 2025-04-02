package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/saswatax/rss-aggregator/internal/config"
	"github.com/saswatax/rss-aggregator/internal/database"
)

type Server struct {
	Config *config.Config
	DB     *database.Queries
}

func NewServer(s Server) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.RegisterRoutes(r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + s.Config.Env.PORT,
	}

	return srv
}
