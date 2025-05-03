package mysqlrepositories

import (
	"github.com/robertvitoriano/penguin-server/internal/models"
	"gorm.io/gorm"
)

type PlayerRepository struct {
	Db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{Db: db}
}

func (r *PlayerRepository) CreatePlayer(newPlayer *models.Player) error {
	if err := r.Db.Create(newPlayer).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerRepository) GetPlayers() ([]models.Player, error) {
	var players []models.Player
	if err := r.Db.Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}
