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

func (r *PlayerMysqlRepository) Save(newPlayer *entities.Player) error {
	if err := r.Db.Create(newPlayer).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlayerMysqlRepository) List() ([]*entities.Player, error) {
	var players []*entities.Player
	if err := r.Db.Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}
func (r *PlayerMysqlRepository) FindByID(id string) (*entities.Player, error) {
	var player entities.Player
	if err := r.Db.Where("id = ?", id).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}
func (r *PlayerMysqlRepository) RemoveByID(id string) (*entities.Player, error) {
	var player entities.Player

	if err := r.Db.Where("id = ?", id).First(&player).Error; err != nil {
		return nil, err
	}

	if err := r.Db.Delete(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}
func (r *PlayerMysqlRepository) FindByUsername(username string) (*entities.Player, error) {
	var player entities.Player
	if err := r.Db.Where("username = ?", username).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}
