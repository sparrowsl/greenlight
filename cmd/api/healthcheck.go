package main

import (
	"net/http"
)

// Return information about the API, including current version
// and operating environment - dev, prod or staging
func (app *application) checkHealth(writer http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(writer, http.StatusOK, map[string]any{"status": "available", "system_info": data}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
