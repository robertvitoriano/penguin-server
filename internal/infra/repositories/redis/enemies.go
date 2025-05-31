package redis

import (
	"fmt"

	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
)

var Enemies = []*entities.Enemy{}

func GetEnemies() []*entities.Enemy {
	return Enemies
}

func CreateEnemy(newEnemy *entities.Enemy) []*entities.Enemy {
	Enemies = append(Enemies, newEnemy)

	return Enemies
}

func KillEnemy(id string) (*entities.Enemy, error) {
	newSlice := []*entities.Enemy{}
	removedEnemy := entities.Enemy{}

	for _, enemy := range Enemies {
		if enemy.ID == nil {
			removedEnemy = *enemy
			continue
		}
		newSlice = append(newSlice, enemy)
	}
	Enemies = newSlice

	if removedEnemy.ID == nil {
		return nil, fmt.Errorf("PLAYER NOT FOUND")
	}

	return &removedEnemy, nil
}
