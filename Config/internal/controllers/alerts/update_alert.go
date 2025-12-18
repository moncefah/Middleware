package alerts

import (
	"encoding/json"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/moncefah/TimeTableAlerter/internal/models"
	"net/http"
)

func (c *Controller) UpdateAlert(w http.ResponseWriter, r *http.Request) {
	var alert models.Alert
	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	// agenda.ID == uuid.Nil here (zero value)

	if err := c.service.UpdateAlert(&alert); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
