package definitions

import (
	"github.com/NavExplorer/navexplorer-sdk-go"
	"github.com/NavPool/navpool-hq-api/internal/config"
	"github.com/NavPool/navpool-hq-api/internal/database"
	"github.com/NavPool/navpool-hq-api/internal/navpool"
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/sarulabs/dingo/v3"
)

var Definitions = []dingo.Def{
	{
		Name: "db",
		Build: func() (*database.Database, error) {
			return database.NewDatabase(
				config.Get().DB.Dialect,
				config.Get().DB.Host,
				config.Get().DB.Port,
				config.Get().DB.DbName,
				config.Get().DB.Username,
				config.Get().DB.Password,
				config.Get().DB.SSLMode,
				config.Get().Debug), nil
		},
	},
	{
		Name: "explorer.api",
		Build: func() (*navexplorer.ExplorerApi, error) {
			return navexplorer.NewExplorerApi(config.Get().Explorer.Url, config.Get().Network)
		},
	},
	{
		Name: "pool.api",
		Build: func() (*navpool.PoolApi, error) {
			return navpool.NewPoolApi(config.Get().Pool.Url, config.Get().Network)
		},
	},
	{
		Name: "account.service",
		Build: func(db *database.Database) (*service.AccountService, error) {
			return service.NewAccountService(db), nil
		},
	},
	{
		Name: "address.service",
		Build: func(db *database.Database, explorerApi *navexplorer.ExplorerApi, poolApi *navpool.PoolApi) (*service.AddressService, error) {
			return service.NewAddressService(db, explorerApi, poolApi), nil
		},
	},
	{
		Name: "dao.service",
		Build: func(db *database.Database, addressService *service.AddressService, poolApi *navpool.PoolApi) (*service.DaoService, error) {
			return service.NewDaoService(db, addressService, poolApi), nil
		},
	},
	{
		Name: "network.service",
		Build: func(accountService *service.AccountService, poolApi *navpool.PoolApi) (*service.NetworkService, error) {
			return service.NewNetworkService(accountService, poolApi), nil
		},
	},
	{
		Name: "staking.service",
		Build: func(addressService *service.AddressService, explorerApi *navexplorer.ExplorerApi) (*service.StakingService, error) {
			return service.NewStakingService(addressService, explorerApi), nil
		},
	},
	{
		Name: "twofactor.service",
		Build: func(accounts *service.AccountService) (*service.TwoFactorService, error) {
			return service.NewTwoFactorService(accounts), nil
		},
	},
}
