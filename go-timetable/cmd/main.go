package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/go-timetable/internal/controllers/events"
	"middleware/go-timetable/internal/controllers/users"
	"middleware/go-timetable/internal/helpers"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()

	// Routes pour USERS (l'exemple)
	r.Route("/users", func(r chi.Router) {
		r.Get("/", users.GetUsers)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(users.Context)
			r.Get("/", users.GetUser)
		})
	})

	// Routes pour EVENTS (votre partie Timetable)
	r.Route("/events", func(r chi.Router) {
		r.Get("/", events.GetEvents)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(events.Context)
			r.Get("/", events.GetEvent)
		})
	})

	port := os.Getenv("API_TIMETABLE_PORT")
	if port == "" {
		port = "8081"
	}

	logrus.Infof("[INFO] API started. Now listening on *:%s", port)
	logrus.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func init() {
	// Initialiser la base users.db
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}

	schemes := []string{
		// Table Users
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL
		);`,
		// Tables Events pour Timetable
		`CREATE TABLE IF NOT EXISTS Event (
			id TEXT PRIMARY KEY,
			summary TEXT NOT NULL,
			location TEXT,
			description TEXT,
			startTime DATETIME NOT NULL,
			endTime DATETIME NOT NULL,
			created DATETIME NOT NULL DEFAULT (datetime('now')),
			lastModifiedTime DATETIME NOT NULL DEFAULT (datetime('now')),
			sequence INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS Resource_Event (
			ressourceId INTEGER NOT NULL,
			eventId TEXT NOT NULL,
			PRIMARY KEY (ressourceId, eventId),
			FOREIGN KEY (eventId) REFERENCES Event(id) ON DELETE CASCADE
		);`,
	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}

	helpers.CloseDB(db)
	logrus.Info("Database tables initialized successfully")
}
