package address

import (
	"errors"
	"github.com/NavExplorer/navexplorer-sdk-go"
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/database"
	"github.com/NavPool/navpool-hq-api/navpool"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/address/model"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
)

var (
	ErrorInvalidAddress              = errors.New("The address provided is invalid")
	ErrorColdStakingDetected         = errors.New("The address provided is a cold stakign address")
	ErrorSpendingAddressAlreadyInUse = errors.New("The spending address provided is already in use")
	ErrorUnableToSaveAddress         = errors.New("Unable to save the address")
)

func CreateNewAddress(addressDto AddressDto, user account.User) (address *model.Address, err error) {
	poolAddress, err := getPoolAddress(addressDto.Hash, addressDto.Signature)
	if err != nil {
		log.Print(err)
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
		return
	}
	defer db.Close()

	err = db.Create(address).Error
	if err != nil {
		log.Print(err)
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
		log.Print(err)
		return
	}

	var hashes = make([]string, len(addresses))
	for _, address := range addresses {
		hashes = append(hashes, address.StakingAddress)
	}
	balances, err := explorerApi.GetBalances(hashes)
	if err != nil {
		log.Print(err)
		return
	}
	for i := range addresses {
		for _, balance := range balances {
			if addresses[i].StakingAddress == balance.Address {
				addresses[i].Balance = balance.ColdStakedBalance
			}
		}
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

func getPoolAddress(hash string, signature string) (poolAddress navpool.PoolAddress, err error) {
	poolApi, err := navpool.NewPoolApi(config.Get().Pool.Url, config.Get().SelectedNetwork)
	if err != nil {
		log.Print(err)
		return
	}

	return poolApi.GetPoolAddress(hash, signature)
}
