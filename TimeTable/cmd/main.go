package main

import (
	"github.com/go-chi/chi/v5"
	events_consumers "github.com/moncefah/TimeTableAlerter/internal/consumers/events"
	events "github.com/moncefah/TimeTableAlerter/internal/controllers/events"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	_ "github.com/moncefah/TimeTableAlerter/internal/models"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	if err := helpers.InitNats(natsURL); err != nil {
		logrus.Fatalf("error while connecting to nats: %s", err.Error())
	}
	defer helpers.CloseNats()

	go func() {
		consumer, err := events_consumers.EventConsumer()
		if err != nil {
			logrus.Warnf("error during nats consumer creation : %v", err)
			return
		}
		if err := events_consumers.Consume(*consumer); err != nil {
			logrus.Warnf("error during nats consume : %v", err)
		}
	}()

	r := chi.NewRouter()
	r.Route("/events", func(r chi.Router) { // route /events
		r.Get("/", events.GetEvents)          // GET /users
		r.Route("/{id}", func(r chi.Router) { // route /events/{id}
			r.Use(events.Context)       // Use Context method to get event ID
			r.Get("/", events.GetEvent) // GET /events/{id}
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
		`CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY NOT NULL UNIQUE,
		agenda_ids TEXT NOT NULL DEFAULT '[]',  -- JSON array
		uid TEXT NOT NULL,
		description TEXT,
		name TEXT NOT NULL,
		start TEXT NOT NULL,                   -- RFC3339
		end TEXT NOT NULL,                     -- RFC3339
		location TEXT,
		last_update TEXT NOT NULL              -- RFC3339
	);`,
	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}
