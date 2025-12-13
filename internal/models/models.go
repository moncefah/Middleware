package models

import (
	"github.com/gofrs/uuid"
)

type User struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}
type Agenda struct {
	ID    uuid.UUID // identifiant interne
	Name  string    // nom lisible (ex: "M1 ISIMA groupe A")
	UcaID string    // identifiant UCA de l’EDT
}
type Alert struct {
	ID       uuid.UUID // identifiant interne
	AgendaID uuid.UUID // lien vers l’agenda surveillé
	Email    string    // destinataire
}
