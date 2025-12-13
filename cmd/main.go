package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/moncefah/TimeTableAlerter/internal/controllers/users"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/sirupsen/logrus"

	"net/http"

	_ "github.com/moncefah/TimeTableAlerter/internal/models"
)

func main() {
	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) { // route /users
		r.Get("/", users.GetUsers)            // GET /users
		r.Route("/{id}", func(r chi.Router) { // route /users/{id}
			r.Use(users.Context)      // Use Context method to get user ID
			r.Get("/", users.GetUser) // GET /users/{id}
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
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL
		);`,
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
