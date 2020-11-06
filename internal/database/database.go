package database

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"log"
)

var (
	ErrDatabaseConnection = errors.New("Failed to connect to the database")
)

type Database struct {
	dialect string
	args    string
	debug   bool
}

func NewDatabase(dialect string, host string, port int, dbName, username, password, sslMode string, debug bool) *Database {
	args := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbName, username, password, sslMode)

	return &Database{
		dialect: dialect,
		args:    args,
		debug:   debug,
	}
}

func (d *Database) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(d.dialect, d.args)
	if err != nil {
		logrus.WithError(err).Error("Failed to open connection")
		return nil, ErrDatabaseConnection
	}

	db.LogMode(d.debug)

	return db, nil
}

func Close(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		log.Panic("Failed to close database connection: ", err)
	}
}
