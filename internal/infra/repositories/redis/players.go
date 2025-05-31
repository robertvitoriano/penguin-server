package redis

import (
	"fmt"

	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
)

var Players = []*entities.Player{}

func GetPlayers() []*entities.Player {
	return Players
}

func CreatePlayer(newPlayer *entities.Player) []*entities.Player {
	Players = append(Players, newPlayer)

	return Players
}

func RemoveByID(id string) (*entities.Player, error) {
	newSlice := []*entities.Player{}
	removedPlayer := entities.Player{}

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

func FindPlayerByUsername(username string) (entities.Player, error) {
	for _, player := range Players {

		if player.Username == username {
			return *player, nil
		}

	}
	return entities.Player{}, fmt.Errorf("player not found")
}
