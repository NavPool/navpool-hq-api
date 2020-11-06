package service

import (
	"github.com/NavPool/navpool-hq-api/internal/database"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AccountService struct {
	db *database.Database
}

func NewAccountService(db *database.Database) *AccountService {
	return &AccountService{db}
}

func (s *AccountService) CreateUser(username string, password string) (*account.User, error) {
	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer database.Close(db)

	password, err = s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &account.User{Username: username, Password: password, Active: true, TwoFactor: &account.TwoFactor{Active: false}}
	err = db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AccountService) GetUserByUsernamePassword(username string, password string, relationships ...string) (*account.User, error) {
	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer database.Close(db)

	user := new(account.User)
	if db.Take(&user, "username = ?", username).RecordNotFound() {
		log.Printf("User not found: %s\n", username)
		err = gorm.ErrRecordNotFound
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logrus.WithError(err).Errorf("Password not valid for user: %s", username)
		return nil, gorm.ErrRecordNotFound
	}

	s.retrieveRelationships(db, user, relationships...)

	return user, nil
}

func (s *AccountService) UsernameExists(username string) (bool, error) {
	db, err := s.db.Connect()
	if err != nil {
		return false, err
	}
	defer database.Close(db)

	var count int
	err = db.Table("users").Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 1, err
}

func (s *AccountService) GetUserByClaim(claimUser account.User, relationships ...string) (*account.User, error) {
	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer database.Close(db)

	user := new(account.User)
	db.Where(&account.User{ID: claimUser.ID}).First(&user)
	s.retrieveRelationships(db, user, relationships...)

	return user, nil
}

func (s *AccountService) UpdateUser(user *account.User) error {
	db, err := s.db.Connect()
	if err != nil {
		return err
	}
	defer database.Close(db)

	db.Save(&user)

	return nil
}

func (s *AccountService) DeleteUser(username string, soft bool) error {
	db, err := s.db.Connect()
	if err != nil {
		return err
	}
	defer database.Close(db)

	if !soft {
		db = db.Unscoped()
	}

	return db.Delete(account.User{}, "username = ?", username).Error
}

func (s *AccountService) GetUserCount() (int, error) {
	db, err := s.db.Connect()
	if err != nil {
		return 0, err
	}
	defer database.Close(db)

	count := 0
	err = db.Table("users").Count(&count).Error

	return count, err
}

func (s *AccountService) retrieveRelationships(db *gorm.DB, user *account.User, relationships ...string) {
	set := make(map[string]bool)
	for _, v := range relationships {
		set[v] = true
	}

	if set["TwoFactor"] {
		var twoFactor account.TwoFactor
		user.TwoFactor = &twoFactor

		db.Model(&user).Related(&twoFactor)
	}
}

func (s *AccountService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
