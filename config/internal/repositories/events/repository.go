package events

import (
	"database/sql"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
)

func Upsert(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// Check if event exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Event WHERE uid = ?)", event.UID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = db.Exec(`
			UPDATE Event 
			SET summary = ?, description = ?, location = ?, start = ?, end = ?
			WHERE uid = ?`,
			event.Summary, event.Description, event.Location, event.Start, event.End, event.UID)
	} else {
		_, err = db.Exec(`
			INSERT INTO Event (uid, summary, description, location, start, end)
			VALUES (?, ?, ?, ?, ?, ?)`,
			event.UID, event.Summary, event.Description, event.Location, event.Start, event.End)
	}

	return err
}

func GetByUID(uid string) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	var e models.Event
	err = db.QueryRow("SELECT uid, summary, description, location, start, end FROM Event WHERE uid = ?", uid).
		Scan(&e.UID, &e.Summary, &e.Description, &e.Location, &e.Start, &e.End)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &e, nil
}
