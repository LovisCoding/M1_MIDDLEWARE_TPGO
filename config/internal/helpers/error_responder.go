package helpers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/config/internal/models"
	"net/http"
)

// RespondError
// This function is for handling different error types
func RespondError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	if _, isErr := err.(*models.ErrorNotFound); isErr {
		status = http.StatusNotFound
	}

	if _, isErr := err.(*models.ErrorUnprocessableEntity); isErr {
		status = http.StatusUnprocessableEntity
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	// if error is not generic, we can send message
	if status != http.StatusInternalServerError {
		body, _ := json.Marshal(err)
		w.Write(body)
	}

	// logging error
	logrus.WithError(err).Printf("An error occured with http code %d", status)
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
