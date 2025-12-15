package agendas

import (
	"encoding/json"
	"github.com/moncefah/TimeTableAlerter/internal/dto"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"net/http"
)

func (c *Controller) CreateAgenda(w http.ResponseWriter, r *http.Request) {
	var agenda dto.CreateAgendaRequest

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&agenda); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	// agenda.ID == uuid.Nil here (zero value)

	if err := c.service.CreateAgenda(&agenda); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
