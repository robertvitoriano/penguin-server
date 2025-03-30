package models

type Position struct {
	X float32
	Y float32
}
type Player struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Color    string   `json:"color"`
	Position Position `json:"position"`
}
