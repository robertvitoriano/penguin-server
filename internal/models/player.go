package models

type Position struct {
	X        *float64 `json:"x"`
	Y        *float64 `json:"y" `
	PlayerId *int     `json:"playerId"`
}
type Player struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Color    string    `json:"color"`
	Position *Position `json:"position"`
}
