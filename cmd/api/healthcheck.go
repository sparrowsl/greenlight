package main

import (
	"fmt"
	"net/http"
)

// Return information about the API, including current version
// and operating environment - dev, prod or staging
func (app *application) checkHealth(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "status: available")
	fmt.Fprintf(writer, "environment: %s\n", app.config.env)
	fmt.Fprintf(writer, "version: %s\n", version)
}
