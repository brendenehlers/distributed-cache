package server

import (
	"log"
	"net/http"
	"time"
)

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			dur := time.Since(start)
			log.Printf("[%s] %s %dμs", r.Method, r.URL.Path, dur.Microseconds())
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
