package main

import (
	"context"
	"net/http"

	"gopkg.in/mgo.v2"
)

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

// Return the API key of a Context.
func APIKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

// Wrapper Handler for API Key.
func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
		fn(w, r.WithContext(ctx))
	}
}

// Wrapper Handler for CORS.
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

// Hardcode API Key
func isValidAPIKey(key string) bool {
	return key == "abc123"
}

// API Server.
type Server struct {
	db *mgo.Session
}

func main() {

}
