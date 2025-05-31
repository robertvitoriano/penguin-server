package entities

import "time"

type Player struct {
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Color          string     `json:"color"`
	Position       *Position  `gorm:"embedded"`
	LastTimeOnline *time.Time `json:"lastTimeOnline"`
}
