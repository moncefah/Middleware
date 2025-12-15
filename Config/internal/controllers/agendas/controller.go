package agendas

import agendasService "github.com/moncefah/TimeTableAlerter/internal/services/agendas"

type Controller struct {
	service *agendasService.Service
}

func NewController(service *agendasService.Service) *Controller {
	return &Controller{service: service}
}
