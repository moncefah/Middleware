package alerts

import (
	"encoding/json"
	"github.com/moncefah/TimeTableAlerter/internal/dto"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"net/http"
)

func (c *Controller) CreateAlert(w http.ResponseWriter, r *http.Request) {
	var alert dto.CreateAlertRequest

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	// alert.ID == uuid.Nil here (zero value)

	if err := c.service.CreateAlert(&alert); err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
