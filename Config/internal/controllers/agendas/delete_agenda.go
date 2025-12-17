package agendas

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"net/http"
)

func (c *Controller) DeleteAgenda(w http.ResponseWriter, r *http.Request) {
	type Id struct {
		ID string `json:"id"`
	}

	id_to_delete := Id{}
	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&id_to_delete); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}
	print(id_to_delete.ID == "")

	// agenda.ID == uuid.Nil here (zero value)
	uuid_id, _ := uuid.FromString(id_to_delete.ID)

	if err := c.service.DeleteAgenda(&uuid_id); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
