package http

import (
	"context"
	"net/http"

	"github.com/uris77/auth0"
)

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			f(w, r)
		}
	}
}

// EnableCors enables CORS
func EnableCors() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Referer, Connection")
			f(w, r)
		}
	}
}

// VerifyToken is a middleware that verifies an auth0 token.
func VerifyToken(jwkUrl, aud, iss string, auth0Client auth0.Auth0) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// OPTIONS request might not include the Authorization header.
			// We don't need to verify a token for OPTIONS.
			if r.Method == "OPTIONS" {
				f(w, r)
				return
			}
			token := r.Header.Get("Authorization")
			jwtToken, err := auth0Client.Validate(jwkUrl, aud, iss, token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", jwtToken)
			r = r.WithContext(ctx)
			f(w, r)
		}
	}
}
