package mysql

import (
	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
	"gorm.io/gorm"
)

type EnemiesMySqlRepository struct {
	Db *gorm.DB
}
type EnemyQuery struct {
	ID   int
	Name string
}

func NewEnemiesRepository(db *gorm.DB) *EnemiesMySqlRepository {
	return &EnemiesMySqlRepository{
		Db: db,
	}
}

func (r *EnemiesMySqlRepository) CreateEnemy(newEnemy *entities.Enemy) error {
	if err := r.Db.Create(newEnemy).Error; err != nil {
		return err
	}
	return nil
}

func (r *EnemiesMySqlRepository) GetEnemies() ([]entities.Enemy, error) {
	var enemies []entities.Enemy
	if err := r.Db.Find(&enemies).Error; err != nil {
		return nil, err
	}

	return enemies, nil

}

func (r *EnemiesMySqlRepository) FindEnemy(criteria EnemyQuery) (*entities.Enemy, error) {
	var enemy *entities.Enemy

	err := r.Db.Where(criteria).First(&enemy).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return enemy, nil

}
