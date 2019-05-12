package account

import (
	"github.com/NavPool/navpool-hq-api/database"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CreateUser(username string, password string) (err error) {
	db, err := database.NewConnection()
	if err != nil {
		return err
	}
	defer database.Close(db)

	password, err = hashPassword(password)
	if err != nil {
		return
	}

	user := &User{Username: username, Password: password, Active: true, TwoFactor: &TwoFactor{Active: false}}
	return db.Create(user).Error
}

func GetUserByUsernamePassword(username string, password string, relationships ...string) (user User, err error) {
	db, err := database.NewConnection()
	if err != nil {
		return
	}
	defer database.Close(db)

	if db.Take(&user, "username = ?", username).RecordNotFound() {
		log.Printf("User not found: %s\n", username)
		err = gorm.ErrRecordNotFound
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password not valid for user: %s\n", username)
		err = gorm.ErrRecordNotFound
		return
	}

	retrieveRelationships(db, &user, relationships...)

	return user, nil
}

func GetUserByClaim(claimUser User, relationships ...string) (user User, err error) {
	db, err := database.NewConnection()
	if err != nil {
		log.Print(err)
		err = ErrUnableToValidateUser
		return
	}
	defer database.Close(db)

	db.Where(&User{ID: claimUser.ID}).First(&user)

	retrieveRelationships(db, &user, relationships...)

	return
}

func UpdateUser(user User) (err error) {
	db, err := database.NewConnection()
	if err != nil {
		return
	}
	defer database.Close(db)

	db.Save(&user)

	return nil
}

func DeleteUser(username string, soft bool) (err error) {
	db, err := database.NewConnection()
	if err != nil {
		return err
	}
	defer database.Close(db)

	if !soft {
		db = db.Unscoped()
	}

	return db.Delete(User{}, "username = ?", username).Error
}

func retrieveRelationships(db *gorm.DB, user *User, relationships ...string) {
	set := make(map[string]bool)
	for _, v := range relationships {
		set[v] = true
	}

	if set["TwoFactor"] {
		var twoFactor TwoFactor
		user.TwoFactor = &twoFactor

		db.Model(&user).Related(&twoFactor)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
