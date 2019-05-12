package account

import (
	"github.com/NavExplorer/navexplorer-api-go/service/address"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

type Login struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	TwoFactor string `form:"twoFactor" json:"twoFactor"`
}

type User struct {
	ID          uuid.UUID         `gorm:"type:uuid;primary_key;" json:"id"`
	Username    string            `gorm:"unique;not null" json:"username,omitempty"`
	Password    string            `json:"-"`
	Active      bool              `json:"active,omitempty"`
	LastLoginAt *time.Time        `json:"lastlogin_at,omitempty"`
	DeletedAt   *time.Time        `sql:"index" json:"deleted_at,omitempty"`
	CreatedAt   *time.Time        `json:"created_at,omitempty"`
	UpdatedAt   *time.Time        `json:"update_at,omitempty"`
	Addresses   []address.Address `json:"addresses,omitempty"`
	TwoFactor   *TwoFactor        `json:"two_factor,omitempty"`
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
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	return scope.SetColumn("ID", id)
}
