package models

import (
	"github.com/gofrs/uuid"
)

type Agenda struct {
	ID    uuid.UUID `json:"id"`    // identifiant interne
	Name  string    `json:"name"`  // nom lisible (ex: "M1 ISIMA groupe A")
	UcaID string    `json:"ucaId"` // identifiant UCA de l’EDT
}
type Alert struct {
	ID       uuid.UUID `json:"id"`       // identifiant interne
	AgendaID uuid.UUID `json:"agendaId"` // lien vers l’agenda surveillé
	Email    string    `json:"email"`    // destinataire
}
