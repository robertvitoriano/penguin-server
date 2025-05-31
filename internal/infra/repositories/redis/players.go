package redis

import (
	"fmt"

	"github.com/robertvitoriano/penguin-server/internal/models"
)

var Players = []*models.Player{}

func GetPlayers() []*models.Player {
	return Players
}

func CreatePlayer(newPlayer *models.Player) []*models.Player {
	Players = append(Players, newPlayer)

	return Players
}

func RemoveByID(id string) (*models.Player, error) {
	newSlice := []*models.Player{}
	removedPlayer := models.Player{}

	for _, player := range Players {
		if player.ID == id {
			removedPlayer = *player
			continue
		}
		newSlice = append(newSlice, player)
	}
	Players = newSlice

	if removedPlayer.ID == "" {
		return nil, fmt.Errorf("PLAYER NOT FOUND")
	}

	return &removedPlayer, nil
}

func FindPlayerByUsername(username string) (models.Player, error) {
	for _, player := range Players {

		if player.Username == username {
			return *player, nil
		}

	}
	return models.Player{}, fmt.Errorf("player not found")
}
