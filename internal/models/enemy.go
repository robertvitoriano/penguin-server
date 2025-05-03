package models

type Enemy struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Color    string    `json:"color"`
	Position *Position `json:"position"`
}
