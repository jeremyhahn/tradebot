package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type FiatPriceServiceImpl struct {
	ctx        common.Context
	datasource common.FiatPriceService
	common.FiatPriceService
}

func NewFiatPriceService(ctx common.Context, exchangeService ExchangeService) (common.FiatPriceService, error) {
	var datasource common.FiatPriceService
	user := ctx.GetUser()
	if user == nil {
		ctx.GetLogger().Debug("[FiatPriceService] Using SlickCharts data source for historical prices")
		datasource = NewSlickChartsService(ctx)
	} else {
		fiatExchange := ctx.GetUser().GetFiatExchange()
		if fiatExchange == "" {
			datasource = NewSlickChartsService(ctx)
		} else {
			ctx.GetLogger().Debugf("[FiatPriceService] Using %s as data source for historical prices", fiatExchange)
			ex, err := exchangeService.GetExchange(fiatExchange)
			datasource = ex.(common.FiatPriceService)
			if err != nil {
				return nil, err
			}
		}
	}
	return &FiatPriceServiceImpl{
		ctx:        ctx,
		datasource: datasource}, nil
}

func (service *FiatPriceServiceImpl) GetPriceAt(currency string, date time.Time) (*common.Candlestick, error) {
	return service.datasource.GetPriceAt(currency, date)
}
