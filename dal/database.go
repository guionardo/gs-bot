package dal

import (
	"github.com/guionardo/gs-bot/configuration"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

func GetDatabase(cfg configuration.RepositoryConfiguration) (*gorm.DB, error) {
	if database != nil {
		return database, nil
	}
	db, err := gorm.Open(sqlite.Open(cfg.ConnectionString), &gorm.Config{})
	if err == nil {
		database = db
	}
	return database, err
}
