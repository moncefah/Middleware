package alerts

import (
	"encoding/json"
	"net/http"

	"github.com/moncefah/TimeTableAlerter/internal/helpers"
)

func (c *Controller) GetAlerts(w http.ResponseWriter, _ *http.Request) {
	// calling service
	alerts, err := c.service.GetAllAlert()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alerts)
	_, _ = w.Write(body)
	return
}
