package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes(r *chi.Mux) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server running"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", s.Signup)
		r.Post("/login", s.Login)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(s.Authenticate)
		r.Get("/me", s.GetProfile)
	})

	r.Route("/feeds", func(r chi.Router) {
		r.Use(s.Authenticate)
		r.Post("/", s.CreateFeed)
		r.Get("/", s.GetFeeds)
	})

	r.Route("/feed-follows", func(r chi.Router) {
		r.Use(s.Authenticate)
		r.Post("/", s.CreateFeedFollow)
		r.Get("/", s.GetFeedFollows)
		r.Delete("/{id}", s.DeleteFeedFollow)
	})

	r.Route("/posts", func(r chi.Router) {
		r.Use(s.Authenticate)
		r.Get("/", s.GetPosts)
	})
}
