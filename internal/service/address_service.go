package service

import (
	"errors"
	"github.com/NavExplorer/navexplorer-sdk-go"
	"github.com/NavPool/navpool-hq-api/internal/database"
	"github.com/NavPool/navpool-hq-api/internal/navpool"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/NavPool/navpool-hq-api/internal/service/address/model"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type AddressService struct {
	db       *database.Database
	explorer *navexplorer.ExplorerApi
	pool     *navpool.PoolApi
}

func NewAddressService(db *database.Database, explorer *navexplorer.ExplorerApi, pool *navpool.PoolApi) *AddressService {
	return &AddressService{db, explorer, pool}
}

var (
	ErrSpendingAddressAlreadyInUse = errors.New("The spending address provided is already in use")
	ErrUnableToSaveAddress         = errors.New("Unable to save the address")
)

func (s *AddressService) CreateNewAddress(hash string, signature string, user account.User) (*model.Address, error) {
	poolAddress, err := s.pool.GetPoolAddress(hash, signature)
	if err != nil {
		return nil, err
	}

	address := &model.Address{
		UserID:             user.ID,
		SpendingAddress:    hash,
		StakingAddress:     poolAddress.StakingAddress,
		ColdStakingAddress: poolAddress.ColdStakingAddress,
	}

	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Create(address).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"addresses_spending_address_key\"") {
			return nil, ErrSpendingAddressAlreadyInUse
		}

		return nil, ErrUnableToSaveAddress
	}

	return address, nil
}

func (s *AddressService) GetAddresses(user account.User) ([]*model.Address, error) {
	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer database.Close(db)

	addresses := make([]*model.Address, 0)
	db.Where(&model.Address{UserID: user.ID}).Order("created_at asc").Find(&addresses)

	var hashes = make([]string, len(addresses))
	for _, address := range addresses {
		hashes = append(hashes, address.StakingAddress)
	}

	balances, err := s.explorer.GetBalances(hashes)
	if err != nil {
		logrus.WithError(err).Error("Failed to get balances")
		return nil, err
	}

	for i := range addresses {
		for _, balance := range balances {
			if addresses[i].StakingAddress == balance.Address {
				addresses[i].Balance = balance.ColdStakedBalance
			}
		}
	}

	return addresses, nil
}

func (s *AddressService) GetAddress(id uuid.UUID, user account.User) (*model.Address, error) {
	db, err := s.db.Connect()
	if err != nil {
		return nil, err
	}
	defer database.Close(db)

	address := new(model.Address)
	db.Where(&model.Address{UserID: user.ID, ID: id}).First(&address)

	balances, err := s.explorer.GetBalances([]string{address.StakingAddress})
	if err != nil {
		logrus.WithError(err).Error("Failed to get balance")
		return nil, err
	}
	if len(balances) == 1 {
		address.Balance = balances[0].ColdStakedBalance
	}

	return address, nil
}

func (s *AddressService) DeleteAddress(id uuid.UUID, user account.User) error {
	db, err := s.db.Connect()
	if err != nil {
		return err
	}
	defer database.Close(db)

	return db.Where(&model.Address{UserID: user.ID, ID: id}).Delete(&model.Address{}).Error
}

func (s *AddressService) GetPoolBalance() (float64, error) {
	db, err := s.db.Connect()
	if err != nil {
		return 0, err
	}
	defer database.Close(db)

	var address model.Address
	err = db.Model(&model.Address{}).Select("sum(balance) as balance").Scan(&address).Error
	if err != nil {
		return 0, err
	}

	return address.Balance, nil
}

func (s *AddressService) GetStakingAddressesForUser(user account.User) ([]string, error) {
	addresses, err := s.GetAddresses(user)
	if err != nil {
		return nil, err
	}

	stakingAddresses := make([]string, 0)
	for _, address := range addresses {
		stakingAddresses = append(stakingAddresses, address.StakingAddress)
	}

	return stakingAddresses, nil
}
