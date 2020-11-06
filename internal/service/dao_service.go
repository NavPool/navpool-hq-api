package service

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/internal/database"
	"github.com/NavPool/navpool-hq-api/internal/navpool"
	"github.com/NavPool/navpool-hq-api/internal/resource/dto"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/NavPool/navpool-hq-api/internal/service/communityFund/model"
	log "github.com/sirupsen/logrus"
)

var (
	ErrorUnableToGetProposalVotes       = errors.New("Unable to retrieve proposal votes")
	ErrorUnableToGetPaymentRequestVotes = errors.New("Unable to retrieve payment request votes")
	ErrorUnableToMatchVote              = errors.New("Unable to match vote")
)

type DaoService struct {
	db        *database.Database
	addresses *AddressService
	pool      *navpool.PoolApi
}

func NewDaoService(db *database.Database, addresses *AddressService, pool *navpool.PoolApi) *DaoService {
	return &DaoService{db, addresses, pool}
}

func (s *DaoService) GetProposalVotes(user account.User) ([]*model.Vote, error) {
	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer database.Close(db)

	votes := make([]*model.Vote, 0)
	err = db.Where(&model.Vote{UserID: user.ID, Type: model.VoteTypeProposal}).Find(&votes).Error
	if err != nil {
		log.WithError(err).Error(ErrorUnableToGetProposalVotes.Error())
		return nil, ErrorUnableToGetProposalVotes
	}

	return votes, nil
}

func (s *DaoService) UpdateProposalVotes(voteDtos []*dto.Vote, user account.User) error {
	db, err := s.db.Connect()
	if err != nil {
		return err
	}
	defer database.Close(db)

	votes, err := s.GetProposalVotes(user)
	if err != nil {
		return err
	}

	tx := db.Begin()

	modifiedVotes := make([]*model.Vote, 0)
	for _, voteDto := range voteDtos {
		vote, err := matchedVote(voteDto.Hash, model.VoteTypeProposal, votes)
		if err == nil {
			err = tx.Model(&vote).Updates(model.Vote{Choice: voteDto.Choice, Committed: false}).Error
		} else {
			newVote := &model.Vote{
				UserID:    user.ID,
				Type:      model.VoteTypeProposal,
				Hash:      voteDto.Hash,
				Choice:    voteDto.Choice,
				Committed: false,
			}
			err = tx.Create(newVote).Error
			vote = newVote
		}

		if err != nil {
			log.WithError(err).Error(err.Error())
			tx.Rollback()
			return err
		}
		modifiedVotes = append(modifiedVotes, vote)
	}

	err = tx.Commit().Error
	if err != nil {
		log.WithError(err).Error(err.Error())
		return err
	}
	err = s.updatePoolVotes(modifiedVotes, user)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return err
	}

	return nil
}

func (s *DaoService) GetPaymentRequestVotes(user account.User) ([]*model.Vote, error) {
	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer database.Close(db)

	votes := make([]*model.Vote, 0)
	err = db.Where(&model.Vote{UserID: user.ID, Type: model.VoteTypePaymentRequest}).Find(&votes).Error
	if err != nil {
		log.WithError(err).Error(ErrorUnableToGetPaymentRequestVotes.Error())
		return nil, ErrorUnableToGetPaymentRequestVotes
	}

	return votes, nil
}

func (s *DaoService) UpdatePaymentRequestVotes(voteDtos []*dto.Vote, user account.User) error {
	db, err := s.db.Connect()
	if err != nil {
		return err
	}
	defer database.Close(db)

	votes, err := s.GetPaymentRequestVotes(user)
	if err != nil {
		return err
	}

	tx := db.Begin()

	modifiedVotes := make([]*model.Vote, 0)
	for _, voteDto := range voteDtos {
		vote, err := matchedVote(voteDto.Hash, model.VoteTypePaymentRequest, votes)
		if err == nil {
			err = tx.Model(&vote).Updates(model.Vote{Choice: voteDto.Choice, Committed: false}).Error
		} else {
			newVote := &model.Vote{
				UserID:    user.ID,
				Type:      model.VoteTypePaymentRequest,
				Hash:      voteDto.Hash,
				Choice:    voteDto.Choice,
				Committed: false,
			}
			err = tx.Create(newVote).Error
			vote = newVote
		}

		if err != nil {
			tx.Rollback()
			return err
		}
		modifiedVotes = append(modifiedVotes, vote)
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	err = s.updatePoolVotes(modifiedVotes, user)
	if err != nil {
		return err
	}

	return nil
}

func matchedVote(hash string, voteType model.VoteType, votes []*model.Vote) (*model.Vote, error) {
	for _, vote := range votes {
		if hash == vote.Hash && voteType == vote.Type {
			return vote, nil
		}
	}

	return nil, ErrorUnableToMatchVote
}

func (s *DaoService) updatePoolVotes(votes []*model.Vote, user account.User) error {
	db, err := s.db.Connect()
	if err != nil {
		return err
	}
	defer database.Close(db)

	addresses, err := s.addresses.GetAddresses(user)
	if err != nil {
		return err
	}

	voteOptions := map[model.VoteChoice]string{
		"YES":     "yes",
		"NO":      "no",
		"ABSTAIN": "remove",
	}

	for _, address := range addresses {
		for _, vote := range votes {
			if vote.Type == model.VoteTypeProposal {
				err = s.pool.ProposalVote(address.SpendingAddress, vote.Hash, voteOptions[vote.Choice])
			} else if vote.Type == model.VoteTypePaymentRequest {
				err = s.pool.PaymentRequestVote(address.SpendingAddress, vote.Hash, voteOptions[vote.Choice])
			}

			if err == nil {
				vote.Committed = true
				db.Save(&vote)
			} else {
				log.WithError(err).Error("Failed to update pool votes")
			}
		}
	}

	return nil
}
