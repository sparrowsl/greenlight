package main

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sparrowsl/greenlight/internal/data"
	"github.com/sparrowsl/greenlight/internal/validator"
	"golang.org/x/time/rate"
)

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mux     sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mux.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Minute*3 {
					delete(clients, ip)
				}
			}

			mux.Unlock()
		}
	}()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(request.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(writer, request, err)
				return
			}

			mux.Lock()

			if _, exists := clients[ip]; !exists {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mux.Unlock()
				app.rateLimitExceededResponse(writer, request)
				return
			}

			mux.Unlock()
		}

		next.ServeHTTP(writer, request)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Authorization")

		authorizationHeader := request.Header.Get("Authorization")

		if authorizationHeader == "" {
			request = app.contextSetUser(request, data.AnonymousUser)
			next.ServeHTTP(writer, request)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		token := headerParts[1]

		v := validator.New()
		if data.ValidateTokenPlainText(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(writer, request)
			default:
				app.serverErrorResponse(writer, request, err)
			}

			return
		}

		request = app.contextSetUser(request, user)

		next.ServeHTTP(writer, request)
	})
}
