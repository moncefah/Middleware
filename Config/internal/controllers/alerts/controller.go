package alerts

import "github.com/moncefah/TimeTableAlerter/internal/services/alerts"

type Controller struct {
	service *alerts.Service
}

func NewController(service *alerts.Service) *Controller {
	return &Controller{service: service}
}
