package models

// AgendaItem represents an agenda item
type AgendaItem struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Time     string `json:"time"`
	Location string `json:"location"`
}
