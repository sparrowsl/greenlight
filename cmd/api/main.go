package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

// Return information about the API, including current version
// and operating environment - dev, prod or staging
func (app *application) checkHealth(writer http.ResponseWriter, request *http.Request) {

}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 5000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.checkHealth)

	server := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf(":%d", cfg.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	fmt.Printf("starting %s server on %s\n", cfg.env, server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
