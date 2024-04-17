package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	shutdownError := make(chan error)

	// Start a background go routine to check for signal terminations
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Println("shutting down server", map[string]string{"signal": s.String()})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		shutdownError <- server.Shutdown(ctx)
	}()

	app.logger.Printf("starting %s server on %s\n", app.config.env, server.Addr)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Println("stopped server", map[string]string{
		"addr": server.Addr,
	})

	return nil
}
