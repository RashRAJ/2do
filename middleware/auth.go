package middleware

import (
	"log"
	"net/http"
	"time"
	// Import any other packages you need
)

// Auth middleware function
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication logic here
		// For example, check for JWT token in the Authorization header

		// If authentication fails
		// w.WriteHeader(http.StatusUnauthorized)
		// json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		// return

		// If authentication succeeds, call the next handler
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
