package account

import (
	"github.com/NavPool/navpool-hq-api/service/address/model"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

type Login struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	TwoFactor string `form:"twoFactor" json:"twoFactor"`
}

type Register struct {
	Username        string `form:"username" json:"username" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm" binding:"required"`
}

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

type TwoFactor struct {
	ID       uint      `gorm:"primary_key" json:"-"`
	UserID   uuid.UUID `gorm:"unique;type:uuid;column:user_id;not null;" json:"-"`
	Active   bool      `json:"active"`
	Secret   *string   `json:"-"`
	LastUsed *int      `json:"-"`
}

func (t *TwoFactor) Disable() *TwoFactor {
	t.Active = false
	t.Secret = nil

	return t
}

func (t *TwoFactor) Enable() *TwoFactor {
	t.Active = true

	return t
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
