package network

type PoolStats struct {
	Staking  bool `json:"staking"`
	Accounts int  `json:"accounts"`
	Weight   int  `json:"weight"`
}
