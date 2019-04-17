package database

import (
	"fmt"
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func NewConnection() (db *gorm.DB) {
	args := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		config.Get().DB.Host,
		config.Get().DB.Port,
		config.Get().DB.DbName,
		config.Get().DB.Username,
		config.Get().DB.Password,
		config.Get().DB.SSLMode)

	db, err := gorm.Open(config.Get().DB.Dialect, args)
	if err != nil {
		log.Panic("Failed to connect database: ", err, args)
	}

	return
}

func Close(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		log.Panic("Failed to close database connection: ", err)
	}
}
