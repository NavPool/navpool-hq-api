package account

import (
	"github.com/NavPool/navpool-hq-api/internal/service/address/model"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;" json:"id"`
	Username    string          `gorm:"unique;not null" json:"username,omitempty"`
	Password    string          `json:"-"`
	Active      bool            `json:"active,omitempty"`
	LastLoginAt *time.Time      `json:"lastlogin_at,omitempty"`
	DeletedAt   *time.Time      `sql:"index" json:"deleted_at,omitempty"`
	CreatedAt   *time.Time      `json:"created_at,omitempty"`
	UpdatedAt   *time.Time      `json:"update_at,omitempty"`
	Addresses   []model.Address `json:"addresses,omitempty"`
	TwoFactor   *TwoFactor      `json:"two_factor,omitempty"`
}

func (u *User) TwoFactorExists() bool {
	return u.TwoFactor != nil
}

func (u *User) TwoFactorActive() bool {
	return u.TwoFactorExists() && u.TwoFactor.Active
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
