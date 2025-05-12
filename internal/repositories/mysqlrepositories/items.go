package mysqlrepositories

import (
	"github.com/robertvitoriano/penguin-server/internal/models"
	"gorm.io/gorm"
)

type ItemsRepository struct {
	Db *gorm.DB
}
type ItemQuery struct {
	ID   int
	Type string
}

func NewItemsRepository(db *gorm.DB) *ItemsRepository {
	return &ItemsRepository{
		Db: db,
	}
}

func (r *ItemsRepository) CreateItem(newItem *models.Item) error {
	if err := r.Db.Create(newItem).Error; err != nil {
		return err
	}
	return nil
}

func (r *ItemsRepository) GetItems() ([]models.Enemy, error) {
	var items []models.Enemy
	if err := r.Db.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}
func (r *ItemsRepository) FindItem(criteria ItemQuery) (*models.Item, error) {
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
