package staking

import (
	"github.com/NavExplorer/navexplorer-sdk-go"
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/address"
	"github.com/getsentry/raven-go"
)

func GetStakingRewardsForUser(user account.User) (accountRewards AccountRewards, err error) {
	stakingAddresses, err := address.GetStakingAddressesForUser(user)
	if err != nil {
		err = ErrorUnableToGetStakingReport
		return
	}

	explorerApi, err := navexplorer.NewExplorerApi(config.Get().Explorer.Url, config.Get().SelectedNetwork)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		err = ErrorUnableToGetStakingReport
		return
	}

	rewards, err := explorerApi.GetStakingRewardsForAddresses(stakingAddresses)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		err = ErrorUnableToGetStakingReport
		return
	}

	for _, reward := range rewards {
		for _, period := range reward.Periods {
			switch period.Period {
			case "last24Hours":
				{
					accountRewards.Last24Hours.Stakes += period.Stakes
					accountRewards.Last24Hours.Balance += period.Balance
					break
				}
			case "last7Days":
				{
					accountRewards.Last7Days.Stakes += period.Stakes
					accountRewards.Last7Days.Balance += period.Balance
					break
				}
			case "last30Days":
				{
					accountRewards.Last30Days.Stakes += period.Stakes
					accountRewards.Last30Days.Balance += period.Balance
					break
				}
			case "lastYear":
				{
					accountRewards.LastYear.Stakes += period.Stakes
					accountRewards.LastYear.Balance += period.Balance
					break
				}
			case "all":
				{
					accountRewards.All.Stakes += period.Stakes
					accountRewards.All.Balance += period.Balance
					break
				}
			}
		}

		return

	}
}
