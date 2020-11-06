package network

type PoolStats struct {
	Enabled  bool `json:"enabled"`
	Staking  bool `json:"staking"`
	Accounts int  `json:"accounts"`
	Weight   int  `json:"weight"`
}
