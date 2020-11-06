package account

import uuid "github.com/satori/go.uuid"

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
