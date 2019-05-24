package migrate

import (
	"github.com/NavPool/navpool-hq-api/database"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/address/model"
	model2 "github.com/NavPool/navpool-hq-api/service/communityFund/model"
)

func Migrate() {
	db, err := database.NewConnection()
	if err != nil {
		return
	}
	defer db.Close()

	db.AutoMigrate(&account.User{}, &account.TwoFactor{}, model.Address{}, model2.Vote{})
}
