package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Handler:      app.routes(),
		Addr:         fmt.Sprintf(":%d", app.config.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	app.logger.Printf("starting %s server on %s\n", app.config.env, server.Addr)
	return server.ListenAndServe()
}
