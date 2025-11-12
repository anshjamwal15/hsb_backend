package entities

type PCOSQuestion struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Type     string   `json:"type"` // multiple_choice, yes_no, scale
	Options  []string `json:"options,omitempty"`
}
