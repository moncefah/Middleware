package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/moncefah/TimeTableAlerter/internal/controllers/agendas"
	"github.com/moncefah/TimeTableAlerter/internal/controllers/alerts"

	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/sirupsen/logrus"

	"net/http"

	_ "github.com/moncefah/TimeTableAlerter/internal/models"
)

func main() {
	r := chi.NewRouter()

	r.Route("/agendas", func(r chi.Router) {
		r.Get("/", agendas.GetAgendas)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(agendas.Context)
			r.Get("/", agendas.GetAgenda)

		})
	})
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", alerts.GetAlerts)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(alerts.Context)
			r.Get("/", alerts.GetAlert)

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

		`CREATE TABLE IF NOT EXISTS agendas (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
    		uca_id VARCHAR(255) 
		);`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			agenda_id VARCHAR(255) NOT NULL ,
    		email VARCHAR(255) NOT NULL,
    		FOREIGN KEY (agenda_id) REFERENCES agendas(id)
		);`,
	}
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}
