package redis

import (
	"fmt"

	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
)

var Players = []*entities.Player{}

func NewPlayerRepository() *PlayerRedisRepository {
	return &PlayerRedisRepository{}
}
func (p *PlayerRedisRepository) List() ([]*entities.Player, error) {
	return Players, nil
}

type PlayerRedisRepository struct {
}

func (p *PlayerRedisRepository) Save(newPlayer *entities.Player) error {
	Players = append(Players, newPlayer)
	return nil
}

func (p *PlayerRedisRepository) RemoveByID(id string) (*entities.Player, error) {
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

func (p *PlayerRedisRepository) FindByUsername(username string) (*entities.Player, error) {
	for _, player := range Players {

		if player.Username == username {
			return player, nil
		}

	}
	return nil, fmt.Errorf("player not found")
}
func (p *PlayerRedisRepository) FindByID(id string) (*entities.Player, error) {
	for _, player := range Players {

		if player.ID == id {
			return player, nil
		}

	}
	return nil, fmt.Errorf("Player not found")
}
