package auth

import (
	"github.com/NavPool/navpool-hq-api/database"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GetUserByUsernamePassword(username string, password string) (interface{}, error) {
	db := database.NewConnection()
	defer database.Close(db)

	var user User
	if db.Take(&user, "username = ?", username).RecordNotFound() {
		log.Printf("User not found: %s\n", username)
		return nil, gorm.ErrRecordNotFound
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password not valid for user: %s\n", username)
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}
