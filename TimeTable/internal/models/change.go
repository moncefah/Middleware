package models

type EventChange struct {
	EventID    string               `json:"eventId"`
	UID        string               `json:"uid"`
	Summary    string               `json:"summary"`
	Start      string               `json:"start"`
	End        string               `json:"end"`
	Location   string               `json:"location"`
	ChangeType string               `json:"changeType"`
	Changes    map[string][2]string `json:"changes"`
}
