package alerts

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/models"
	repository "github.com/moncefah/TimeTableAlerter/internal/repositories/alerts"
	"github.com/sirupsen/logrus"
)

func GetAllAlert() ([]models.Alert, error) {
	var err error
	// calling repository
	alerts, err := repository.GetAllAlerts()
	// managing errors
	if err != nil {
		logrus.Errorf("error retrieving users : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving agendas",
		}
	}

	return alerts, nil
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	alert, err := repository.GetAlertById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "agenda not found",
			}
		}
		logrus.Errorf("error retrieving user %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving user %s", id.String()),
		}
	}

	return alert, err
}
