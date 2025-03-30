package repositories

import "github.com/robertvitoriano/penguin-server/models"

var Players = []models.Player{}

func GetPlayers() []models.Player {
	return Players
}

func CreatePlayer(newPlayer models.Player) []models.Player {
	Players = append(Players, newPlayer)

	return Players
}
