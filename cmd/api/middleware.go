package main

import (
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
