package alerts

import (
	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/moncefah/TimeTableAlerter/internal/models"
)

func GetAllAlerts() ([]models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM alerts")
	helpers.CloseDB(db)
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
func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * FROM alerts WHERE id=?", id.String())
	helpers.CloseDB(db)

	var data models.Alert
	err = row.Scan(&data.ID, &data.AgendaID, &data.Email)
	if err != nil {
		return nil, err
	}
	return &data, err
}
