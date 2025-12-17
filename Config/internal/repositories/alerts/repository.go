package alerts

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

func (r *Repository) GetAllAlerts() ([]models.Alert, error) {

	rows, err := r.db.Query("SELECT * FROM alerts")
	if err != nil {
		return nil, err
	}

	alerts := []models.Alert{}
	for rows.Next() {
		var data models.Alert
		err = rows.Scan(&data.ID, &data.AgendaID, &data.Email)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, data)
	}
	_ = rows.Close()

	return alerts, err

}
func (r *Repository) GetAlertById(id uuid.UUID) (*models.Alert, error) {

	row := r.db.QueryRow("SELECT * FROM alerts WHERE id=?", id.String())

	var data models.Alert
	err := row.Scan(&data.ID, &data.AgendaID, &data.Email)
	if err != nil {
		return nil, err
	}
	return &data, err
}

func (r *Repository) CreateAlert(agenda *models.Agenda) error {
	_, err := r.db.Exec(`
		INSERT INTO agendas (id, name, uca_id)
		VALUES ($1, $2, $3)
	`, agenda.ID, agenda.Name, agenda.UcaID)

	return err
}
