package mysqlrepositories

import "gorm.io/gorm"

type EnemiesRepository struct {
	Db *gorm.DB
}

func NewEnemiesRepository(db *gorm.DB) *EnemiesRepository {
	return &EnemiesRepository{
		Db: db,
	}
}

func (r *EnemiesRepository) CreateEnemy() {

}
