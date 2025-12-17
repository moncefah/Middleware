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

func (r *Repository) CreateAlert(alert *models.Alert) error {
	_, err := r.db.Exec(`
		INSERT INTO alerts (id, agenda_id, email)
		VALUES (?, ?, ?)
	`,
		alert.ID.String(),
		alert.AgendaID.String(),
		alert.Email,
	)

	return err
}


