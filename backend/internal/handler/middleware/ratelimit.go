package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	clients = make(map[string]*rate.Limiter)
	mu      sync.Mutex
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = rate.NewLimiter(rate.Every(1*time.Minute), 5)
		}

		if !clients[ip].Allow() {
			mu.Unlock()
			http.Error(w, "You have made too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}