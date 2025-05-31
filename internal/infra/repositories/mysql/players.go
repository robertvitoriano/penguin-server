package mysql

import (
	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
	"gorm.io/gorm"
)

type PlayerMysqlRepository struct {
	Db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerMysqlRepository {
	return &PlayerMysqlRepository{Db: db}
}

func (r *PlayerMysqlRepository) CreatePlayer(newPlayer *entities.Player) error {
	if err := r.Db.Create(newPlayer).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerMysqlRepository) GetPlayers() ([]entities.Player, error) {
	var players []entities.Player
	if err := r.Db.Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}
