package models

import (
	"github.com/gofrs/uuid"
)

type User struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}
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
