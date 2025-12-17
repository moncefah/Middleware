package alerts

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
)

func (c *Controller) GetAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertId, _ := ctx.Value("alertId").(uuid.UUID) // getting key set in context.go

	alert, err := c.service.GetAlertById(alertId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alert)
	_, _ = w.Write(body)
	return
}
