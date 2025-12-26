package helpers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/go-timetable/internal/models"
	"net/http"
)

func RespondError(err error) (body []byte, status int) {
	status = http.StatusInternalServerError

	if _, isErr := err.(*models.ErrorNotFound); isErr {
		status = http.StatusNotFound
	}

	if _, isErr := err.(*models.ErrorUnprocessableEntity); isErr {
		status = http.StatusUnprocessableEntity
	}

	if status != http.StatusInternalServerError {
		body, _ = json.Marshal(err)
	}

	logrus.WithError(err).Printf("An error occurred with http code %d", status)

	return
}
