package models

type Recommendation struct {
	Ticker           string  `json:"ticker"`
	Company          string  `json:"company"`
	TargetFrom       string  `json:"target_from"` // Leer como STRING
	TargetTo         string  `json:"target_to"`   // Leer como STRING
	CambioPorcentual float64 `json:"cambio_porcentual"`
}
