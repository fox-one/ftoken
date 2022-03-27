package cmd

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/quorum"
	assetz "github.com/fox-one/ftoken/service/asset"
	walletz "github.com/fox-one/ftoken/service/wallet"
	"github.com/fox-one/ftoken/store/asset"
	"github.com/fox-one/ftoken/store/order"
	"github.com/fox-one/ftoken/store/transaction"
	"github.com/fox-one/ftoken/store/wallet"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/property"
	"github.com/fox-one/pkg/store/db"
	propertystore "github.com/fox-one/pkg/store/property"
)

func provideSystem(ctx context.Context, client *mixin.Client, factories []core.Factory) (core.System, error) {
	system := core.System{
		ClientID:     cfg.Dapp.ClientID,
		ClientSecret: cfg.Dapp.ClientSecret,
		Version:      rootCmd.Version,
		WhiteList:    map[string]bool{},
		Fees:         map[string]*core.Fee{},
	}

	for _, fee := range cfg.Fees {
		system.Fees[fee.Platform] = &core.Fee{
			FeeAssetID: fee.FeeAssetID,
			FeeAmount:  fee.FeeAmount,
		}
	}
	for _, id := range cfg.WhiteList {
		system.WhiteList[id] = true
	}
	return system, nil
}

func provideMixinClient() *mixin.Client {
	c, err := mixin.NewFromKeystore(&cfg.Dapp.Keystore)
	if err != nil {
		panic(err)
	}

	return c
}

func provideDatabase() (*db.DB, error) {
	database, err := db.Open(cfg.DB)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(database); err != nil {
		return nil, err
	}

	return database, nil
}

func providePropertyStore(db *db.DB) property.Store {
	return propertystore.New(db)
}

func provideWalletStore(db *db.DB) core.WalletStore {
	return wallet.New(db)
}

func provideOrderStore(db *db.DB) core.OrderStore {
	return order.New(db)
}

func provideTransactionStore(db *db.DB) core.TransactionStore {
	return transaction.New(db)
}

func provideAssetStore(db *db.DB) core.AssetStore {
	return asset.New(db)
}

func provideWalletService(c *mixin.Client) core.WalletService {
	return walletz.New(walletz.Config{Pin: cfg.Dapp.Pin}, c)
}

func provideAssetService(c *mixin.Client) core.AssetService {
	return assetz.New(c)
}

func provideAllFactories() []core.Factory {
	return []core.Factory{
		provideQuorumFactory(),
	}
}

func provideQuorumFactory() core.Factory {
	return quorum.New(cfg.Eth.Endpoint, cfg.Eth.PrivateKey, cfg.Eth.FactoryContract, *cfg.Eth.MaxGasPrice)
}
