package mysqlrepositories

import (
	"github.com/robertvitoriano/penguin-server/internal/models"
	"gorm.io/gorm"
)

type EnemiesRepository struct {
	Db *gorm.DB
}

func NewEnemiesRepository(db *gorm.DB) *EnemiesRepository {
	return &EnemiesRepository{
		Db: db,
	}
}

func (r *EnemiesRepository) CreateEnemy(newEnemy *models.Enemy) error {
	if err := r.Db.Create(newEnemy).Error; err != nil {
		return err
	}
	return nil
}

func (r *EnemiesRepository) GetEnemies() ([]models.Enemy, error) {
	var enemies []models.Enemy
	if err := r.Db.Find(&enemies).Error; err != nil {
		return nil, err
	}

	return enemies, nil

}
