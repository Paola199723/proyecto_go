package models

type ListResponse struct {
	ID          uint   `gorm:"primaryKey"`
	Ticker      string `json:"ticker"`
	Target_from string `json:"target_from"`
	Target_to   string `json:"target_to"`
	Company     string `json:"company"`
	Action      string `json:"action"`
	Brokerage   string `json:"brokerage"`
	Rating_from string `json:"rating_from"`
	Rating_to   string `json:"rating_to"`
	Time        string `json:"time"`
}
type StockResponse struct {
	Items    []ListResponse `json:"items"`
	NextPage string         `json:"next_page"`
}
