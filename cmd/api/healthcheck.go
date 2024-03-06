package main

import (
	"fmt"
	"net/http"
)

// Return information about the API, including current version
// and operating environment - dev, prod or staging
func (app *application) checkHealth(writer http.ResponseWriter, request *http.Request) {
	js := `{"status": available, "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(js))
}
