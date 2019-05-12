package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type VoteType string
type VoteChoice string

const (
	VoteTypeProposal       VoteType = "PROPOSAL"
	VoteTypePaymentRequest VoteType = "PAYMENT_REQUEST"
)

type Vote struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	UserID    uuid.UUID  `gorm:"type:uuid;column:user_id;not null;" json:"-"`
	Type      VoteType   `json:"type"`
	Hash      string     `json:"hash"`
	Choice    VoteChoice `json:"vote"`
	CreatedAt *time.Time `json:"_"`
	UpdatedAt *time.Time `json:"_"`
	Committed bool       `json:"committed"`
}
