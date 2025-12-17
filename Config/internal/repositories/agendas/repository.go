package agendas

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetAllAgendas retrieves all agendas from database
func (r *Repository) GetAllAgendas() ([]models.Agenda, error) {
	rows, err := r.db.Query(`
		SELECT *
		FROM agendas
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendas []models.Agenda

	for rows.Next() {
		var agenda models.Agenda
		if err := rows.Scan(&agenda.ID, &agenda.Name, &agenda.UcaID); err != nil {
			return nil, err
		}
		agendas = append(agendas, agenda)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return agendas, nil
}

// GetAgendaById retrieves a single agenda by its ID
func (r *Repository) GetAgendaById(id uuid.UUID) (*models.Agenda, error) {
	var agenda models.Agenda

	err := r.db.QueryRow(`
		SELECT id, name, uca_id
		FROM agendas
		WHERE id = $1
	`, id).Scan(&agenda.ID, &agenda.Name, &agenda.UcaID)

	if err != nil {
		return nil, err
	}

	return &agenda, nil
}

// CreateAgenda inserts a new agenda into database
func (r *Repository) CreateAgenda(agenda *models.Agenda) error {
	_, err := r.db.Exec(`
		INSERT INTO agendas (id, name, uca_id)
		VALUES ($1, $2, $3)
	`, agenda.ID, agenda.Name, agenda.UcaID)

	return err
}
func (r *Repository) UpdateAgenda(agenda *models.Agenda) error {
	_, err := r.db.Exec(`
	UPDATE agendas
	SET name = ?, uca_id = ?
	WHERE id = ?
`,
		agenda.Name,
		agenda.UcaID,
		agenda.ID,
	)
	return err

}

func (r *Repository) DeleteAgenda(id *uuid.UUID) error {
	_, err := r.db.Exec(`
	DELETE FROM agendas
	
	WHERE id = ?
`,
		id,
	)
	return err

}
