package alerts

import (
	"database/sql"
	"fmt"
	"github.com/moncefah/TimeTableAlerter/internal/dto"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/models"
	repository "github.com/moncefah/TimeTableAlerter/internal/repositories/alerts"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{repository: repository}
}
func (s *Service) GetAllAlert() ([]models.Alert, error) {
	var err error
	// calling repository
	alerts, err := s.repository.GetAllAlerts()
	// managing errors
	if err != nil {
		logrus.Errorf("error retrieving users : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving agendas",
		}
	}

	return alerts, nil
}

func (s *Service) GetAlertById(id uuid.UUID) (*models.Alert, error) {
	alert, err := s.repository.GetAlertById(id)
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
func (s *Service) CreateAlert(alertReqDto *dto.CreateAlertRequest) error {
	newID, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("error generating uuid: %s", err.Error())
		return &models.ErrorGeneric{
			Message: "Failed to generate agenda ID",
		}
	}
	agenda := models.Alert{
		ID:       newID,
		AgendaID: alertReqDto.AgendaId,
		Email:    alertReqDto.Email,
	}

	if err := s.repository.CreateAlert(&agenda); err != nil {
		logrus.Errorf("error creating alert: %s", err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while creating alert",
		}
	}

	return nil
}
