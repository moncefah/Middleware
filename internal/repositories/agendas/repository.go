package agendas

import (
	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/moncefah/TimeTableAlerter/internal/models"
)

func GetAllAgendas() ([]models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM agendas")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	agendas := []models.Agenda{}
	for rows.Next() {
		var data models.Agenda
		err = rows.Scan(&data.ID, &data.Name, &data.UcaID)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, data)
	}
	_ = rows.Close()

	return agendas, err

}
func GetAgendaById(id uuid.UUID) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * FROM agendas WHERE id=?", id.String())
	helpers.CloseDB(db)

	var data models.Agenda
	err = row.Scan(&data.ID, &data.Name, &data.UcaID)
	if err != nil {
		return nil, err
	}
	return &data, err
}

func CreateAgenda(agenda *models.Agenda) (*models.Agenda, error) {

	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(
		`INSERT INTO agendas (id, name, uca_id) VALUES (?, ?, ?)`,
		agenda.ID,
		agenda.Name,
		agenda.UcaID,
	)
	if err != nil {
		return nil, err
	}
	return agenda, nil
}
