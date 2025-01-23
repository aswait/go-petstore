package middleware

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)

func UserUnloggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "invalid or missing token", http.StatusForbidden)
			return
		}

		if token == nil || jwt.Validate(token, jwt.WithClaimValue("username", username)) != nil {
			http.Error(w, "permission error", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UnloggedInDelete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api_key")
		var tokenString string

		tokenString = apiKey

		token, err := jwt.Parse([]byte(tokenString))
		if err != nil || jwt.Validate(token) != nil {
			http.Error(w, "invalid token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UnloggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Authorization := r.Header.Get("Authorization")

		tokenString := Authorization

		token, err := jwt.Parse([]byte(tokenString))

		if err != nil || jwt.Validate(token) != nil {
			http.Error(w, "invalid token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
