package migrate

import (
	"github.com/NavPool/navpool-hq-api/internal/di"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/NavPool/navpool-hq-api/internal/service/address/model"
	model2 "github.com/NavPool/navpool-hq-api/internal/service/communityFund/model"
)

func Migrate() {
	db, err := di.Get().GetDb().Connect()
	if err != nil {
		return
	}
	defer db.Close()

	db.AutoMigrate(&account.User{}, &account.TwoFactor{}, model.Address{}, model2.Vote{})
}
