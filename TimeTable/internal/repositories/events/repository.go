package events

import (
	"database/sql"
	"encoding/json"
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
