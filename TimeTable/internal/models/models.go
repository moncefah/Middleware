package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Event struct {
	ID       uuid.UUID `json:"id"`       // Internal DB ID
	AgendaID uuid.UUID `json:"agendaId"` // One agenda only
	UID      string    `json:"uid"`      // Stable iCal / ADE UID
	Name     string    `json:"name"`     // Course name
	Start    time.Time `json:"start"`    // Start datetime
	End      time.Time `json:"end"`      // End datetime
	Location string    `json:"location"` // Room
	Checksum string    `json:"checksum"` // Used for change detection
	LastSeen time.Time `json:"lastSeen"` // When scheduler last observed it
}
