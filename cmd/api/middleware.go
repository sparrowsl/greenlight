package main

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

func (app *application) rateLimit(next http.Handler) http.Handler {
	var (
		mux     sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(writer, request, err)
			return
		}

		mux.Lock()

		if _, exists := clients[ip]; !exists {
			clients[ip] = rate.NewLimiter(2, 4)
		}

		if !clients[ip].Allow() {
			mux.Unlock()
			app.rateLimitExceededResponse(writer, request)
			return
		}

		mux.Unlock()

		next.ServeHTTP(writer, request)
	})
}
