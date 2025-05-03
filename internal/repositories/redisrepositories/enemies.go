package redisrepositories

import (
	"fmt"

	"github.com/robertvitoriano/penguin-server/internal/models"
)

var Enemies = []*models.Enemy{}

func GetEnemies() []*models.Enemy {
	return Enemies
}

func CreateEnemy(newEnemy *models.Enemy) []*models.Enemy {
	Enemies = append(Enemies, newEnemy)

	return Enemies
}

func KillEnemy(id string) (*models.Enemy, error) {
	newSlice := []*models.Enemy{}
	removedEnemy := models.Enemy{}

	for _, enemy := range Enemies {
		if enemy.ID == id {
			removedEnemy = *enemy
			continue
		}
		newSlice = append(newSlice, enemy)
	}
	Enemies = newSlice

	if removedEnemy.ID == "" {
		return nil, fmt.Errorf("PLAYER NOT FOUND")
	}

	return &removedEnemy, nil
}
