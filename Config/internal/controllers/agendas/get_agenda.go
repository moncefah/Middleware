package agendas

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
)

func (c *Controller) GetAgenda(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agendaId, _ := ctx.Value("agendaId").(uuid.UUID) // getting key set in context.go

	agenda, err := c.service.GetAgendaById(agendaId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(agenda)
	_, _ = w.Write(body)
	return
}
