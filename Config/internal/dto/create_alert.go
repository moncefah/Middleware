package dto

import "github.com/gofrs/uuid"

type CreateAlertRequest struct {
	AgendaId uuid.UUID `json:"agenda_id"`
	Email    string    `json:"email"`
}
