package network

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/address"
)

var (
	ErrorUnableToRetrieveStats = errors.New("Unable to retrieve stats")
)

func GetPoolStats() (poolStats PoolStats, err error) {
	accounts, err := account.GetUserCount()
	if err == nil {
		poolStats.Accounts = accounts
	}

	balance, err := address.GetPoolBalance()
	if err == nil {
		poolStats.Balance = balance
	}

	return
}
