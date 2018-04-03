package service

import (
	"github.com/jeremyhahn/tradebot/common"
)

type WalletServiceImpl struct {
	ctx              common.Context
	pluginService    PluginService
	fiatPriceService common.FiatPriceService
	currency         string
}

func NewWalletService(ctx common.Context, pluginService PluginService, fiatPriceService common.FiatPriceService) WalletService {
	return &WalletServiceImpl{
		ctx:              ctx,
		pluginService:    pluginService,
		fiatPriceService: fiatPriceService}
}

func (service *WalletServiceImpl) CreateWallet(currency, address string) (common.Wallet, error) {
	constructor, err := service.pluginService.CreateWallet(currency)
	if err != nil {
		service.ctx.GetLogger().Errorf("[WalletService.GetWallet] Failed to load wallet plugin: %s", err.Error())
		return nil, err
	}
	wallet := constructor(&common.WalletParams{
		Context:          service.ctx,
		Address:          address,
		MarketCapService: NewMarketCapService(service.ctx),
		FiatPriceService: service.fiatPriceService})
	return wallet, nil
}
