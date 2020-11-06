package navpool

import (
	"encoding/json"
)

type Info struct {
	Version            int           `json:"version"`
	ProtocolVersion    int           `json:"protocolversion"`
	WalletVersion      int           `json:"walletversion"`
	Balance            float64       `json:"balance"`
	ColdStakingBalance float64       `json:"coldstaking_balance"`
	NewMint            float64       `json:"newmint"`
	Stake              float64       `json:"stake"`
	Blocks             int           `json:"blocks"`
	CommunityFund      CommunityFund `json:"communityfund"`
	TimeOffset         int           `json:"timeoffset"`
	NtpTimeOffset      int           `json:"ntptimeoffset"`
	Connections        int           `json:"connections"`
	Proxy              string        `json:"proxy"`
	TestNet            bool          `json:"testnet"`
	KeyPoolOldest      int           `json:"keypoololdest"`
	KeyPoolSize        int           `json:"keypoolsize"`
	UnlockedUntil      int           `json:"unlocked_until"`
	PayTxFee           float64       `json:"paytxfee"`
	RelayFee           float64       `json:"relayfee"`
	Errors             string        `json:"errors"`
}

type CommunityFund struct {
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
}

type StakingInfo struct {
	Enabled          bool    `json:"enabled"`
	Staking          bool    `json:"staking"`
	Errors           string  `json:"errors"`
	CurrentBlockSize int     `json:"currentblocksize"`
	CurrentBlockTx   int     `json:"currentblocktx"`
	Difficulty       float64 `json:"difficulty"`
	SearchInterval   int     `json:"search-interval"`
	Weight           int     `json:"weight"`
	NetStakeWeight   int     `json:"netstakeweight"`
	ExpectedTime     int     `json:"expectedtime"`
}

func (e *PoolApi) GetInfo() (info Info, err error) {
	method := "/info"

	response, err := e.client.call(method, "GET", nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &info)
	return
}

func (e *PoolApi) GetStakingInfo() (stakingInfo StakingInfo, err error) {
	method := "/info/staking"

	response, err := e.client.call(method, "GET", nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &stakingInfo)
	return
}
