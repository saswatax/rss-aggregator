package server

import (
	"context"
	"log"
	"net/http"

	"github.com/saswatax/rss-aggregator/internal/server/services"
)

func (s *Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := services.GetBearToken(r)
		if err != nil {
			encode(w, http.StatusUnauthorized, err)
			return
		}

		claims, err := services.ParseTokenClaims(token, s.Config.Env.JWT_SECRET)
		if err != nil {
			log.Println(err)
			encode(w, http.StatusUnauthorized, "invalid authentication token")
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID{}, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
