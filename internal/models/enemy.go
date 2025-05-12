package models

type Enemy struct {
	ID       *int      `json:"id"`
	Name     string    `json:"name"`
	Position *Position `json:"position" gorm:"embedded"`
	Killed   bool      `json:"killed"`
	Size     *Size     `json:"size" gorm:"embedded"`
}
