package dto

import "github.com/NavPool/navpool-hq-api/internal/service/communityFund/model"

type Vote struct {
	Type   model.VoteType   `json:"type" binding:"required"`
	Hash   string           `json:"hash" binding:"required"`
	Choice model.VoteChoice `json:"vote" binding:"required"`
}
