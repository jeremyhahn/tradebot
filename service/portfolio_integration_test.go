// +build integration

package service

import (
	"testing"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestPortfolioService_Build(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	marketcapService := NewMarketCapService(ctx.Logger)
	userService := NewUserService(ctx,
		dao.NewUserDAO(ctx),
		marketcapService,
		mapper.NewUserMapper())
	service := NewPortfolioService(ctx, marketcapService, userService)
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"}
	portfolio := service.Build(ctx.GetUser(), currencyPair)
	assert.Equal(t, uint(1), portfolio.GetUser().GetId())
	assert.Equal(t, true, len(portfolio.GetExchanges()) > 0)
	assert.Equal(t, true, len(portfolio.GetWallets()) > 0)
	assert.Equal(t, true, portfolio.GetNetWorth() > 0)
	test.CleanupIntegrationTest()
}

func TestPortfolioService_Stream(t *testing.T) {
	ctx := test.NewIntegrationTestContext()
	marketcapService := NewMarketCapService(ctx.Logger)
	userService := NewUserService(ctx,
		dao.NewUserDAO(ctx),
		marketcapService,
		mapper.NewUserMapper())
	service := NewPortfolioService(ctx, marketcapService, userService)
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"}
	portfolio := <-service.Stream(ctx.GetUser(), currencyPair)
	service.Stop(ctx.GetUser())
	assert.Equal(t, uint(1), portfolio.GetUser().GetId())
	assert.Equal(t, true, len(portfolio.GetExchanges()) > 0)
	assert.Equal(t, true, len(portfolio.GetWallets()) > 0)
	assert.Equal(t, true, portfolio.GetNetWorth() > 0)
	test.CleanupIntegrationTest()
}
