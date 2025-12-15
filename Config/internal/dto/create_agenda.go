package dto

type CreateAgendaRequest struct {
	Name  string `json:"name"`
	UcaID string `json:"ucaId"`
}
