package api

import (
	"context"
	"moh.gov.bz/mch/emtct/internal/auth"
	"net/http"
	"strings"

	"github.com/uris77/auth0"
)

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

// VerifyAuthToken is a middleware that verifies an auth0 token.
func VerifyAuthToken(jwkUrl, aud, iss string, auth0Client auth0.Auth0) Middleware {
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
			email, _ := jwtToken.Get("email")
			ctx := context.WithValue(r.Context(), "user", auth.JwtToken{Email: email.(string)})
			r = r.WithContext(ctx)
			f(w, r)
		}
	}
}

func VerifyToken(userStore auth.UserStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// OPTIONS request might not include the Authorization header.
			// We don't need to verify a token for OPTIONS.
			if r.Method == http.MethodOptions {
				f(w, r)
				return
			}
			h := r.Header
			bearer := h.Get("Authorization")
			if len(strings.Trim(bearer, "")) == 0 {
				// No Authorization Token was provided
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			bearerParts := strings.Split(bearer, " ")
			if bearerParts[0] != "Bearer" {
				// Wrong header format... return error
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token := bearerParts[1]
			jwtToken, err := userStore.VerifyToken(token)
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
