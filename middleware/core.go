package middleware

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
)

// Set the Content Type Header
func ContentTypeHeader(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func CORS(c *web.C, h http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH"},
	})

	return cors.Handler(h)
}
