package middleware

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code and body
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	// Capture the response body
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code and body
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			written:        false,
			body:           &bytes.Buffer{},
		}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log the request details
		duration := time.Since(start)

		// If the request failed (status >= 400), include the response body
		if wrapped.statusCode >= 400 {
			log.Printf("[%s] %s - Status: %d - Duration: %v - Response: %s",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
				wrapped.body.String(),
			)
		} else {
			log.Printf("[%s] %s - Status: %d - Duration: %v",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
			)
		}
	})
}
