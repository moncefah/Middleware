package dto

import "github.com/google/uuid"

type CreateALertRequest struct {
	AgendaId uuid.UUID `json:"agenda_id"`
	Email    string    `json:"email"`
}
