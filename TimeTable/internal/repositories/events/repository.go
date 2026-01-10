package events

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/moncefah/TimeTableAlerter/internal/models"
	"github.com/sirupsen/logrus"
)

func GetAllEvents() ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query(`
		SELECT 
			id,
			agenda_ids,
			uid,
			description,
			name,
			start,
			end,
			location,
			last_update
		FROM events
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event

		var (
			idStr        string
			agendaIDsStr string
			startStr     string
			endStr       string
			updateStr    string
		)

		err := rows.Scan(
			&idStr,
			&agendaIDsStr,
			&e.UID,
			&e.Description,
			&e.Name,
			&startStr,
			&endStr,
			&e.Location,
			&updateStr,
		)
		if err != nil {
			logrus.Error("scan error: ", err)
			return nil, err
		}

		// conversions
		e.ID, _ = uuid.FromString(idStr)
		_ = json.Unmarshal([]byte(agendaIDsStr), &e.AgendaIDs)
		e.Start, _ = helpers.ParseTime(startStr)
		e.End, _ = helpers.ParseTime(endStr)
		e.LastUpdate, _ = helpers.ParseTime(updateStr)

		events = append(events, e)
	}

	return events, nil
}

// GetEventById récupère un événement selon son ID
func GetEventById(id uuid.UUID) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	row := db.QueryRow(`
		SELECT 
			id,
			agenda_ids,
			uid,
			description,
			name,
			start,
			end,
			location,
			last_update
		FROM events
		WHERE id = ?
	`, id.String())

	var e models.Event

	var (
		idStr        string
		agendaIDsStr string
		startStr     string
		endStr       string
		updateStr    string
	)

	err = row.Scan(
		&idStr,
		&agendaIDsStr,
		&e.UID,
		&e.Description,
		&e.Name,
		&startStr,
		&endStr,
		&e.Location,
		&updateStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	e.ID, _ = uuid.FromString(idStr)
	_ = json.Unmarshal([]byte(agendaIDsStr), &e.AgendaIDs)
	e.Start, _ = helpers.ParseTime(startStr)
	e.End, _ = helpers.ParseTime(endStr)
	e.LastUpdate, _ = helpers.ParseTime(updateStr)

	return &e, nil
}

// GetEventByUID récupère un événement selon son UID
func GetEventByUID(uid string) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	row := db.QueryRow(`
		SELECT
			id,
			agenda_ids,
			uid,
			description,
			name,
			start,
			end,
			location,
			last_update
		FROM events
		WHERE uid = ?
	`, uid)

	var e models.Event

	var (
		idStr        string
		agendaIDsStr string
		startStr     string
		endStr       string
		updateStr    string
	)

	err = row.Scan(
		&idStr,
		&agendaIDsStr,
		&e.UID,
		&e.Description,
		&e.Name,
		&startStr,
		&endStr,
		&e.Location,
		&updateStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	e.ID, _ = uuid.FromString(idStr)
	_ = json.Unmarshal([]byte(agendaIDsStr), &e.AgendaIDs)
	e.Start, _ = helpers.ParseTime(startStr)
	e.End, _ = helpers.ParseTime(endStr)
	e.LastUpdate, _ = helpers.ParseTime(updateStr)

	return &e, nil
}

// CreateEvent insère un événement en base
func CreateEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	agendaIDs, err := json.Marshal(event.AgendaIDs)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO events (
			id,
			agenda_ids,
			uid,
			description,
			name,
			start,
			end,
			location,
			last_update
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		event.ID.String(),
		string(agendaIDs),
		event.UID,
		event.Description,
		event.Name,
		event.Start.Format(time.RFC3339),
		event.End.Format(time.RFC3339),
		event.Location,
		event.LastUpdate.Format(time.RFC3339),
	)
	return err
}

// UpdateEvent met à jour un événement en base
func UpdateEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	agendaIDs, err := json.Marshal(event.AgendaIDs)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE events
		SET agenda_ids = ?,
			uid = ?,
			description = ?,
			name = ?,
			start = ?,
			end = ?,
			location = ?,
			last_update = ?
		WHERE id = ?
	`,
		string(agendaIDs),
		event.UID,
		event.Description,
		event.Name,
		event.Start.Format(time.RFC3339),
		event.End.Format(time.RFC3339),
		event.Location,
		event.LastUpdate.Format(time.RFC3339),
		event.ID.String(),
	)
	return err
}
