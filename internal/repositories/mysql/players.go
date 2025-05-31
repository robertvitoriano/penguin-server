package mysql

import (
	"github.com/robertvitoriano/penguin-server/internal/models"
	"gorm.io/gorm"
)

type PlayerMysqlRepository struct {
	Db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerMysqlRepository {
	return &PlayerMysqlRepository{Db: db}
}

func (r *PlayerMysqlRepository) CreatePlayer(newPlayer *models.Player) error {
	if err := r.Db.Create(newPlayer).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerMysqlRepository) GetPlayers() ([]models.Player, error) {
	var players []models.Player
	if err := r.Db.Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}
