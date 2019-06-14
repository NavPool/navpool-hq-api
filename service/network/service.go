package network

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/logger"
	"github.com/NavPool/navpool-hq-api/navpool"
	"github.com/NavPool/navpool-hq-api/service/account"
)

var (
	ErrorUnableToRetrieveStats = errors.New("Unable to network retrieve stats")
)

func GetNetworkStats() (poolStats PoolStats, err error) {
	accounts, err := account.GetUserCount()
	if err != nil {
		logger.LogError(err)
		return
	}

	poolApi, err := navpool.NewPoolApi(config.Get().Pool.Url, config.Get().SelectedNetwork)
	if err != nil {
		logger.LogError(err)
		return
	}

	stakingInfo, err := poolApi.GetStakingInfo()
	if err != nil {
		logger.LogError(err)
		return
	}

	poolStats.Accounts = accounts
	poolStats.Weight = stakingInfo.Weight
	poolStats.Enabled = stakingInfo.Enabled
	poolStats.Staking = stakingInfo.Staking

	return
}
