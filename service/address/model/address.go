package model

import (
	"github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

type Address struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	UserID             uuid.UUID  `gorm:"type:uuid;column:user_id;not null;" json:"-"`
	SpendingAddress    string     `gorm:"not null;unique_index:idx_address" json:"spending_address"`
	StakingAddress     string     `gorm:"not null" json:"staking_address"`
	ColdStakingAddress string     `gorm:"not null" json:"cold_staking_address"`
	Balance            float64    `json:"balance"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"update_at,omitempty"`
}

func (address *Address) BeforeCreate(scope *gorm.Scope) error {
	id, err := uuid.NewV4()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return err
	}

	return scope.SetColumn("ID", id)
}
