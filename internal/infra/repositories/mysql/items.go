package mysql

import (
	"github.com/robertvitoriano/penguin-server/internal/models"
	"gorm.io/gorm"
)

type ItemsMysqlRepository struct {
	Db *gorm.DB
}
type ItemQuery struct {
	ID   int
	Type string
}

func NewItemsRepository(db *gorm.DB) *ItemsMysqlRepository {
	return &ItemsMysqlRepository{
		Db: db,
	}
}

func (r *ItemsMysqlRepository) CreateItem(newItem *models.Item) error {
	if err := r.Db.Create(newItem).Error; err != nil {
		return err
	}
	return nil
}

func (r *ItemsMysqlRepository) GetItems() ([]models.Enemy, error) {
	var items []models.Enemy
	if err := r.Db.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}
func (r *ItemsMysqlRepository) FindItem(criteria ItemQuery) (*models.Item, error) {
	var item *models.Item

	err := r.Db.Where(criteria).First(&item).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return item, nil

}
