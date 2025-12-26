package events

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"middleware/go-timetable/internal/models"
	repository "middleware/go-timetable/internal/repositories/events"
)

// GetAllEvents appelle le repo avec un resourceId optionnel
func GetAllEvents(resourceId string) ([]models.Event, error) {
	events, err := repository.GetAllEvents(resourceId)
	if err != nil {
		logrus.Errorf("error retrieving events: %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving events",
		}
	}

	return events, nil
}

func GetEventById(id string) (*models.Event, error) {
	event, err := repository.GetEventById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "event not found",
			}
		}
		logrus.Errorf("error retrieving event %s : %s", id, err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving event %s", id),
		}
	}

	return event, nil
}