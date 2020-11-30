package address

import (
	"errors"
	"github.com/NavExplorer/navexplorer-sdk-go"
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/database"
	"github.com/NavPool/navpool-hq-api/logger"
	"github.com/NavPool/navpool-hq-api/navpool"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/address/model"
	uuid "github.com/satori/go.uuid"
	"strings"
)

var (
	ErrorUnableToFindAddress          = errors.New("Unable to find the address on your account")
	ErrorUnableToDeleteAddress        = errors.New("Unable to delete the address")
	ErrorSpendingAddressAlreadyInUse  = errors.New("The spending address provided is already in use")
	ErrorUnableToSaveAddress          = errors.New("Unable to save the address")
	ErrorUnableToRetrieveTransactions = errors.New("Unable to retrieve transactions")
)

func CreateNewAddress(addressDto AddressDto, user account.User) (address *model.Address, err error) {
	poolAddress, err := getPoolAddress(addressDto.Hash, addressDto.Signature)
	if err != nil {
		logger.LogError(err)
		return
	}

	address = &model.Address{
		UserID:             user.ID,
		SpendingAddress:    addressDto.Hash,
		StakingAddress:     poolAddress.StakingAddress,
		ColdStakingAddress: poolAddress.ColdStakingAddress,
	}

	db, err := database.NewConnection()
	if err != nil {
		logger.LogError(err)
		return
	}
	defer db.Close()

	err = db.Create(address).Error
	if err != nil {
		logger.LogError(err)

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"addresses_spending_address_key\"") {
			err = ErrorSpendingAddressAlreadyInUse
		} else {
			err = ErrorUnableToSaveAddress
		}
	}

	return address, err
}

func GetAddresses(user account.User) (addresses []model.Address, err error) {
	db, err := database.NewConnection()
	if err != nil {
		return
	}
	defer database.Close(db)

	db.Where(&model.Address{UserID: user.ID}).Order("created_at asc").Find(&addresses)

	explorerApi, err := navexplorer.NewExplorerApi(config.Get().Explorer.Url, config.Get().SelectedNetwork)
	if err != nil {
		logger.LogError(err)
		return
	}

	var hashes = make([]string, len(addresses))
	for _, address := range addresses {
		hashes = append(hashes, address.StakingAddress)
	}

	balances, err := explorerApi.GetBalances(hashes)
	if err != nil {
		logger.LogError(err)
		return
	}

	tx := db.Begin()
	for i := range addresses {
		for _, balance := range balances {
			if addresses[i].StakingAddress == balance.Hash {
				if addresses[i].Balance != float64(balance.Stakable) {
					tx.Save(addresses[i])
				}
				addresses[i].Balance = float64(balance.Stakable)
			}
		}
	}
	tx.Commit()

	return
}

func GetAddress(id uuid.UUID, user account.User) (address model.Address, err error) {
	db, err := database.NewConnection()
	if err != nil {
		return
	}
	defer database.Close(db)

	db.Where(&model.Address{UserID: user.ID, ID: id}).First(&address)

	explorerApi, err := navexplorer.NewExplorerApi(config.Get().Explorer.Url, config.Get().SelectedNetwork)
	if err != nil {
		logger.LogError(err)
		return
	}

	balances, err := explorerApi.GetBalances([]string{address.StakingAddress})
	if err != nil {
		logger.LogError(err)
		return
	}
	if len(balances) == 1 {
		address.Balance = float64(balances[0].Stakable)
	}

	return
}

func DeleteAddress(id uuid.UUID, user account.User) (err error) {
	db, err := database.NewConnection()
	if err != nil {
		return
	}
	defer database.Close(db)

	err = db.Where(&model.Address{UserID: user.ID, ID: id}).Delete(&model.Address{}).Error
	return
}

func GetPoolBalance() (balance float64, err error) {
	db, err := database.NewConnection()
	if err != nil {
		return
	}
	defer database.Close(db)

	var address model.Address
	err = db.Model(&model.Address{}).Select("sum(balance) as balance").Scan(&address).Error
	if err != nil {
		logger.LogError(err)
		return
	}

	balance = address.Balance

	return
}

func GetStakingAddressesForUser(user account.User) (stakingAddresses []string, err error) {
	addresses, err := GetAddresses(user)
	if err != nil {
		return
	}

	for _, address := range addresses {
		stakingAddresses = append(stakingAddresses, address.StakingAddress)
	}

	return
}

func getPoolAddress(hash string, signature string) (poolAddress navpool.PoolAddress, err error) {
	poolApi, err := navpool.NewPoolApi(config.Get().Pool.Url, config.Get().SelectedNetwork)
	if err != nil {
		logger.LogError(err)
		return
	}

	return poolApi.GetPoolAddress(hash, signature)
}
