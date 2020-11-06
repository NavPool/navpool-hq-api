package service

import (
	"github.com/NavExplorer/navexplorer-sdk-go"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/NavPool/navpool-hq-api/internal/service/staking"
)

type StakingService struct {
	addresses *AddressService
	explorer  *navexplorer.ExplorerApi
}

func NewStakingService(addresses *AddressService, explorer *navexplorer.ExplorerApi) *StakingService {
	return &StakingService{addresses, explorer}
}

func (s *StakingService) GetStakingRewardsForUser(user account.User) (*staking.AccountRewards, error) {
	stakingAddresses, err := s.addresses.GetStakingAddressesForUser(user)
	if err != nil {
		return nil, err
	}

	rewards, err := s.explorer.GetStakingRewardsForAddresses(stakingAddresses)
	if err != nil {
		return nil, err
	}

	accountRewards := new(staking.AccountRewards)
	for _, r := range rewards {
		for _, p := range r.Periods {
			switch p.Period {
			case "last24Hours":
				{
					accountRewards.Last24Hours.Stakes += p.Stakes
					accountRewards.Last24Hours.Balance += p.Balance
					break
				}
			case "last7Days":
				{
					accountRewards.Last7Days.Stakes += p.Stakes
					accountRewards.Last7Days.Balance += p.Balance
					break
				}
			case "last30Days":
				{
					accountRewards.Last30Days.Stakes += p.Stakes
					accountRewards.Last30Days.Balance += p.Balance
					break
				}
			case "lastYear":
				{
					accountRewards.LastYear.Stakes += p.Stakes
					accountRewards.LastYear.Balance += p.Balance
					break
				}
			case "all":
				{
					accountRewards.All.Stakes += p.Stakes
					accountRewards.All.Balance += p.Balance
					break
				}
			}
		}
	}

	return accountRewards, nil
}
