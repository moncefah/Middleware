package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/moncefah/TimeTableAlerter/internal/controllers/agendas"
	"github.com/moncefah/TimeTableAlerter/internal/controllers/alerts"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/sirupsen/logrus"
	"net/http"

	agendaCtrl "github.com/moncefah/TimeTableAlerter/internal/controllers/agendas"
	agendaRepo "github.com/moncefah/TimeTableAlerter/internal/repositories/agendas"
	agendaServ "github.com/moncefah/TimeTableAlerter/internal/services/agendas"
)

func main() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}
	defer helpers.CloseDB(db)

	if err := helpers.InitSchema(db); err != nil {
		logrus.Fatalf("error while Initiating database schemas : %s", err.Error())
	}

	agendaRepository := agendaRepo.NewRepository(db)
	agendaService := agendaServ.NewService(agendaRepository)
	agendaControl := agendaCtrl.NewController(agendaService)

	r := chi.NewRouter()

	r.Route("/agendas", func(r chi.Router) {
		r.Get("/", agendaControl.GetAgendas)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(agendas.Context)
			r.Get("/", agendaControl.GetAgenda)

		})
		r.Post("/", agendaControl.CreateAgenda)

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
