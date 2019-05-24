package staking

type AccountRewards struct {
	Last24Hours RewardPeriod `json:"last24Hours"`
	Last7Days   RewardPeriod `json:"last7Days"`
	Last30Days  RewardPeriod `json:"last30Days"`
	LastYear    RewardPeriod `json:"lastYear"`
	All         RewardPeriod `json:"all"`
}

type RewardPeriod struct {
	Stakes  int64 `json:"stakes"`
	Balance int64 `json:"balance"`
}
