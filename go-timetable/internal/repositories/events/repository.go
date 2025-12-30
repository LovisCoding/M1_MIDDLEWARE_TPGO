package events

import (
	"database/sql"
	"middleware/go-timetable/internal/helpers"
	"middleware/go-timetable/internal/models"
)

// GetAllEvents récupère les événements, optionnellement filtrés par ressource
func GetAllEvents(resourceId string) ([]models.Event, error) {
	db, err := helpers.OpenTimetableDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := "SELECT e.id, e.summary, e.location, e.description, e.startTime, e.endTime, e.created, e.lastModifiedTime, e.sequence FROM Event e"
	var args []interface{}

	// Si un resourceId est fourni, on fait la jointure
	if resourceId != "" {
		// Attention à l'espace au début de la chaîne pour éviter "Event eJOIN"
		query += " JOIN Resource_Event re ON e.id = re.eventId WHERE re.ressourceId = ?"
		args = append(args, resourceId)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var event models.Event
		var location, description sql.NullString

		// Scan doit correspondre à l'ordre du SELECT
		err = rows.Scan(&event.Id, &event.Summary, &location, &description, &event.StartTime, &event.EndTime, &event.Created, &event.LastModifiedTime, &event.Sequence)
		if err != nil {
			return nil, err
		}

		if location.Valid {
			event.Location = &location.String
		}
		if description.Valid {
			event.Description = &description.String
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventById(id string) (*models.Event, error) {
	db, err := helpers.OpenTimetableDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := "SELECT id, summary, location, description, startTime, endTime, created, lastModifiedTime, sequence FROM Event WHERE id=?"
	row := db.QueryRow(query, id)

	var event models.Event
	var location, description sql.NullString

	err = row.Scan(&event.Id, &event.Summary, &location, &description, &event.StartTime, &event.EndTime, &event.Created, &event.LastModifiedTime, &event.Sequence)
	if err != nil {
		return nil, err
	}

	if location.Valid {
		event.Location = &location.String
	}
	if description.Valid {
		event.Description = &description.String
	}

	return &event, nil
}