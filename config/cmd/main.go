package main

import (
	_ "middleware/config/api"
	"middleware/config/internal/controllers/alerts"
	"middleware/config/internal/controllers/mails"
	"middleware/config/internal/controllers/resources"
	"middleware/config/internal/helpers"
	_ "middleware/config/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/mails", func(r chi.Router) {
		r.Get("/", mails.GetMails)
		r.Post("/", mails.CreateMail)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", mails.GetMail)
			r.Put("/", mails.UpdateMail)
			r.Delete("/", mails.DeleteMail)
		})
	})

	r.Route("/resources", func(r chi.Router) {
		r.Get("/", resources.GetResources)
		r.Post("/", resources.CreateResource)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", resources.GetResource)
			r.Put("/", resources.UpdateResource)
			r.Delete("/", resources.DeleteResource)

			r.Get("/alerts", alerts.GetAlerts)
			r.Post("/alerts/{mail_id}", alerts.AddAlert)
			r.Delete("/alerts/{mail_id}", alerts.RemoveAlert)
		})
	})

	logrus.Info("[INFO] Web server started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS Mail (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			mail TEXT NOT NULL UNIQUE
		);`,
		`CREATE TABLE IF NOT EXISTS Ressource (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS Alert (
			mailId INTEGER NOT NULL,
			ressourceId INTEGER NOT NULL,
			PRIMARY KEY (mailId, ressourceId),
			FOREIGN KEY (mailId) REFERENCES Mail(id) ON DELETE CASCADE,
			FOREIGN KEY (ressourceId) REFERENCES Ressource(id) ON DELETE CASCADE
		);`,
	}
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}
