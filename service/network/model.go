package network

type PoolStats struct {
	Enabled  bool    `json:"enabled"`
	Accounts int     `json:"accounts"`
	Balance  float64 `json:"balance"`
}
