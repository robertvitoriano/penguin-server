package mysqlrepositories

import (
	"github.com/robertvitoriano/penguin-server/internal/models"
	"gorm.io/gorm"
)

type ItemsRepository struct {
	Db *gorm.DB
}

func NewItemsRepository(db *gorm.DB) *EnemiesRepository {
	return &EnemiesRepository{
		Db: db,
	}
}

func (r *EnemiesRepository) CreateItem(newEnemy *models.Enemy) error {
	if err := r.Db.Create(newEnemy).Error; err != nil {
		return err
	}
	return nil
}

func (r *EnemiesRepository) GetItems() ([]models.Enemy, error) {
	var items []models.Enemy
	if err := r.Db.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}
