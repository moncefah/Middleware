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

	alertCtrl "github.com/moncefah/TimeTableAlerter/internal/controllers/alerts"
	alertRepo "github.com/moncefah/TimeTableAlerter/internal/repositories/alerts"
	alertServ "github.com/moncefah/TimeTableAlerter/internal/services/alerts"
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

	alertRepository := alertRepo.NewRepository(db)
	alertService := alertServ.NewService(alertRepository)
	alertControl := alertCtrl.NewController(alertService)

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
		r.Get("/", alertControl.GetAlerts)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(alerts.Context)
			r.Get("/", alertControl.GetAlert)

		})
		r.Post("/", alertControl.CreateAlert)
	})

	logrus.Info("[INFO] Web server started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))
}
