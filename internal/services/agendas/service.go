package agendas

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/models"
	repository "github.com/moncefah/TimeTableAlerter/internal/repositories/agendas"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{repository: repository}
}

// GetAllAgendas handles business logic for retrieving agendas
func (s *Service) GetAllAgendas() ([]models.Agenda, error) {
	agendas, err := s.repository.GetAllAgendas()
	if err != nil {
		logrus.Errorf("error retrieving agendas: %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving agendas",
		}
	}

	return agendas, nil
}

// GetAgendaById retrieves one agenda and maps repository errors
func (s *Service) GetAgendaById(id uuid.UUID) (*models.Agenda, error) {
	agenda, err := s.repository.GetAgendaById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.ErrorNotFound{
				Message: "Agenda not found",
			}
		}

		logrus.Errorf("error retrieving agenda %s: %s", id, err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving agenda",
		}
	}

	return agenda, nil
}

// CreateAgenda applies business rules and creates agenda
func (s *Service) CreateAgenda(agenda *models.Agenda) error {
	newID, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("error generating uuid: %s", err.Error())
		return &models.ErrorGeneric{
			Message: "Failed to generate agenda ID",
		}
	}

	agenda.ID = newID

	if err := s.repository.CreateAgenda(agenda); err != nil {
		logrus.Errorf("error creating agenda: %s", err.Error())
		return &models.ErrorGeneric{
			Message: "Something went wrong while creating agenda",
		}
	}

	return nil
}
