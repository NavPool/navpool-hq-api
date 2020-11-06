package service

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/internal/navpool"
	"github.com/NavPool/navpool-hq-api/internal/service/network"
	log "github.com/sirupsen/logrus"
)

type NetworkService struct {
	accounts *AccountService
	pool     *navpool.PoolApi
}

func NewNetworkService(accountService *AccountService, poolApi *navpool.PoolApi) *NetworkService {
	return &NetworkService{accountService, poolApi}
}

var (
	ErrorUnableToRetrieveStats = errors.New("Unable to network retrieve stats")
)

func (s *NetworkService) GetNetworkStats() (*network.PoolStats, error) {
	accounts, err := s.accounts.GetUserCount()
	if err != nil {
		log.WithError(err).Error("Failed to get user count")
		return nil, err
	}

	stakingInfo, err := s.pool.GetStakingInfo()
	if err != nil {
		log.WithError(err).Error("Failed to get stkaing info")
		return nil, err
	}

	return &network.PoolStats{
		Enabled:  stakingInfo.Enabled,
		Staking:  stakingInfo.Staking,
		Accounts: accounts,
		Weight:   stakingInfo.Weight,
	}, nil
}
