package main

import (
	"net/http"
)

// Return information about the API, including current version
// and operating environment - dev, prod or staging
func (app *application) checkHealth(writer http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	if err := app.writeJSON(writer, http.StatusOK, data, nil); err != nil {
		app.logger.Print(err)
		http.Error(writer, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
