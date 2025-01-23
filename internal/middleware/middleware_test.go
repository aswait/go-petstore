package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"
)

func TestUserUnloggedIn(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		token          string
		expectedStatus int
	}{
		{
			name:           "invalid token",
			username:       "testuser",
			token:          "invalid-token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "missing token",
			username:       "testuser",
			token:          "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(UserUnloggedIn)
			r.Get("/{username}", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/"+tt.username, nil)
			if tt.token != "" {
				tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
				_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"username": tt.username})
				req.Header.Set("Authorization", "Bearer "+tokenString)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestUnloggedInDelete(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string
		expectedStatus int
	}{
		{
			name:           "invalid token",
			apiKey:         "invalid-token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "missing token",
			apiKey:         "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(UnloggedInDelete)
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req, _ := http.NewRequest("DELETE", "/", nil)
			if tt.apiKey != "" {
				req.Header.Set("api_key", tt.apiKey)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestUnloggedIn(t *testing.T) {
	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "invalid token",
			authorization:  "invalid-token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "missing token",
			authorization:  "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(UnloggedIn)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
