package models

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}
type Player struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Color    string   `json:"color"`
	Position Position `json:"position"`
}
