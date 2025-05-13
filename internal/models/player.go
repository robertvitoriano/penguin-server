package models

type Player struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Color    string    `json:"color"`
	Position *Position `json:"position" gorm:"embedded"`
	Score    *int      `json:"score"`
	Image    *string   `json:"image"`
	Size     *Size     `json:"size" gorm:"embedded"`
}
