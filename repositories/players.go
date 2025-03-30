package repositories

import (
	"fmt"

	"github.com/robertvitoriano/penguin-server/models"
)

var Players = []*models.Player{}

func GetPlayers() []*models.Player {
	return Players
}

func CreatePlayer(newPlayer *models.Player) []*models.Player {
	Players = append(Players, newPlayer)

	return Players
}

func FindPlayerByUsername(username string) (models.Player, error) {
	for _, player := range Players {

		if player.Username == username {
			return *player, nil
		}

	}
	return models.Player{}, fmt.Errorf("player not found")
}
