package database

import (
	"errors"
	"fmt"
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

var (
	ErrorDatabaseConnection = errors.New("Failed to connect to the database")
)

func NewConnection() (db *gorm.DB, err error) {
	args := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		config.Get().DB.Host,
		config.Get().DB.Port,
		config.Get().DB.DbName,
		config.Get().DB.Username,
		config.Get().DB.Password,
		config.Get().DB.SSLMode)

	db, err = gorm.Open(config.Get().DB.Dialect, args)
	if err != nil {
		logger.LogError(err)
		err = ErrorDatabaseConnection
	}

	if config.Get().Debug == true {
		db.LogMode(config.Get().DB.LogMode)
	}

	return
}

func Close(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		log.Panic("Failed to close database connection: ", err)
	}
}
