package models

type Item struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Position    *Position `json:"position" gorm:"embedded"`
	Type        string    `json:"type"`
	CollectedBy *int      `json:"collectedBy"`
	Image       *string   `json:"image"`
}
