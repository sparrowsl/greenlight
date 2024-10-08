package main

import (
	"errors"
	"expvar"
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

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := app.contextGetUser(request)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (app *application) requireActivatedUser(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := app.contextGetUser(request)

		if !user.Activated {
			app.inactiveAccountResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})

	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := app.contextGetUser(request)

		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(writer, request, err)
			return
		}

		if !permissions.Include(code) {
			app.notPermittedResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})

	return app.requireActivatedUser(fn).(http.HandlerFunc)
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Origin")
		writer.Header().Add("Vary", "Access-Control-Request-Method")

		origin := request.Header.Get("Origin")

		if origin != "" {
			for _, i := range app.config.cors.trustedOrigins {
				if origin == i {
					writer.Header().Set("Access-Control-Allow-Origin", origin)

					if request.Method == http.MethodOptions && request.Header.Get("Access-Control-Request-Method") != "" {
						writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						writer.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(writer, request)
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
	var (
		totalRequestsReceived           = expvar.NewInt("total_request_received")
		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
	)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		totalRequestsReceived.Add(1)

		next.ServeHTTP(writer, request)

		totalResponsesSent.Add(1)

		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})

}
