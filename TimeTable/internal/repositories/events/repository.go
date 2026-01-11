package events

import (
	"database/sql"
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
			agenda_id,
			uid,
			name,
			start,
			end,
			location,
			checksum,
			last_seen
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
			idStr       string
			agendaID    string
			startStr    string
			endStr      string
			lastSeenStr string
		)

		err := rows.Scan(
			&idStr,
			&agendaID,
			&e.UID,
			&e.Name,
			&startStr,
			&endStr,
			&e.Location,
			&e.Checksum,
			&lastSeenStr,
		)
		if err != nil {
			logrus.Error("scan error: ", err)
			return nil, err
		}

		// conversions
		e.ID, _ = uuid.FromString(idStr)
		if agendaID != "" {
			e.AgendaID, _ = uuid.FromString(agendaID)
		}
		e.Start, _ = helpers.ParseTime(startStr)
		e.End, _ = helpers.ParseTime(endStr)
		e.LastSeen, _ = helpers.ParseTime(lastSeenStr)

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
			agenda_id,
			uid,
			name,
			start,
			end,
			location,
			checksum,
			last_seen
		FROM events
		WHERE id = ?
	`, id.String())

	var e models.Event

	var (
		idStr       string
		agendaID    string
		startStr    string
		endStr      string
		lastSeenStr string
	)

	err = row.Scan(
		&idStr,
		&agendaID,
		&e.UID,
		&e.Name,
		&startStr,
		&endStr,
		&e.Location,
		&e.Checksum,
		&lastSeenStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	e.ID, _ = uuid.FromString(idStr)
	if agendaID != "" {
		e.AgendaID, _ = uuid.FromString(agendaID)
	}
	e.Start, _ = helpers.ParseTime(startStr)
	e.End, _ = helpers.ParseTime(endStr)
	e.LastSeen, _ = helpers.ParseTime(lastSeenStr)

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
			agenda_id,
			uid,
			name,
			start,
			end,
			location,
			checksum,
			last_seen
		FROM events
		WHERE uid = ?
	`, uid)

	var e models.Event

	var (
		idStr       string
		agendaID    string
		startStr    string
		endStr      string
		lastSeenStr string
	)

	err = row.Scan(
		&idStr,
		&agendaID,
		&e.UID,
		&e.Name,
		&startStr,
		&endStr,
		&e.Location,
		&e.Checksum,
		&lastSeenStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	e.ID, _ = uuid.FromString(idStr)
	if agendaID != "" {
		e.AgendaID, _ = uuid.FromString(agendaID)
	}
	e.Start, _ = helpers.ParseTime(startStr)
	e.End, _ = helpers.ParseTime(endStr)
	e.LastSeen, _ = helpers.ParseTime(lastSeenStr)

	return &e, nil
}

// CreateEvent insère un événement en base
func CreateEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec(`
		INSERT INTO events (
			id,
			agenda_id,
			uid,
			name,
			start,
			end,
			location,
			checksum,
			last_seen
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		event.ID.String(),
		agendaIDValue(event.AgendaID),
		event.UID,
		event.Name,
		event.Start.Format(time.RFC3339),
		event.End.Format(time.RFC3339),
		event.Location,
		event.Checksum,
		event.LastSeen.Format(time.RFC3339),
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

	_, err = db.Exec(`
		UPDATE events
		SET agenda_id = ?,
			uid = ?,
			name = ?,
			start = ?,
			end = ?,
			location = ?,
			checksum = ?,
			last_seen = ?
		WHERE id = ?
	`,
		agendaIDValue(event.AgendaID),
		event.UID,
		event.Name,
		event.Start.Format(time.RFC3339),
		event.End.Format(time.RFC3339),
		event.Location,
		event.Checksum,
		event.LastSeen.Format(time.RFC3339),
		event.ID.String(),
	)
	return err
}

func agendaIDValue(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}
