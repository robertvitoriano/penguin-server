package database

import (
	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Db     *gorm.DB
	Dsn    string
	DbType string
}

func NewMysqlDabase() *Database {
	return &Database{}
}

func (d *Database) Connect() (*gorm.DB, error) {

	var err error

	d.Db, err = gorm.Open(mysql.Open(d.Dsn))

	if err != nil {
		return nil, err
	}

	d.Db.AutoMigrate(&entities.Player{}, &entities.Enemy{}, &entities.Item{})
	return d.Db, nil

}
