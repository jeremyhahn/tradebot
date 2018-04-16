// +build integration

package service

import (
	"os"
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEthereumWebClient_GetWallet(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)

	localAuthService := NewLocalAuthService(ctx, userDAO, userMapper)
	service, err := NewEthereumWebClient(ctx, userDAO, localAuthService, NewMarketCapService(ctx), fiatPriceService)
	assert.Nil(t, err)

	wallet, err := service.GetWallet(os.Getenv("ETH_ADDRESS"))
	assert.Nil(t, err)

	zero := decimal.NewFromFloat(0)
	assert.Equal(t, true, len(wallet.GetAddress()) > 0)
	assert.Equal(t, true, wallet.GetBalance().GreaterThan(zero))
	assert.Equal(t, true, wallet.GetValue().GreaterThan(zero))

	test.CleanupIntegrationTest()
}

func TestEthereumWebClientService_GetTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	walletAddress := os.Getenv("ETH_ADDRESS")

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := NewMarketCapService(ctx)
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)
	ethereumService, err := NewEthereumService(ctx, userDAO, userMapper, marketcapService, exchangeService)
	assert.Nil(t, err)

	walletService := NewWalletService(ctx, pluginService, fiatPriceService)
	userService := NewUserService(ctx, userDAO, userMapper, userExchangeMapper, marketcapService,
		ethereumService, exchangeService, walletService)

	userService.CreateWallet(&dto.UserCryptoWalletDTO{
		Address:  os.Getenv("ETH_ADDRESS"),
		Currency: "ETH"})

	userService.CreateToken(&dto.EthereumTokenDTO{
		Symbol:          os.Getenv("TOKEN_SYMBOL"),
		ContractAddress: os.Getenv("TOKEN_ADDRESS"),
		WalletAddress:   os.Getenv("ETH_ADDRESS")})

	localAuthService := NewLocalAuthService(ctx, userDAO, userMapper)
	webclient, err := NewEthereumWebClient(ctx, userDAO, localAuthService, NewMarketCapService(ctx), fiatPriceService)
	assert.Nil(t, err)

	transactions, err := webclient.GetTransactions()
	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	var totalDeposit, totalWithdrawl, totalFees decimal.Decimal

	for _, tx := range transactions {
		qty, err := decimal.NewFromString(tx.GetQuantity())
		if err != nil {
			panic(err)
		}
		if tx.GetType() == common.DEPOSIT_ORDER_TYPE {
			totalDeposit = totalDeposit.Add(qty)
		} else if tx.GetType() == common.WITHDRAWAL_ORDER_TYPE {
			fee, err := decimal.NewFromString(tx.GetFee())
			if err != nil {
				panic(err)
			}
			totalWithdrawl = totalWithdrawl.Add(qty)
			totalFees = totalFees.Add(fee)
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}

	assert.Equal(t, true, totalDeposit.GreaterThan(decimal.NewFromFloat(0)))
	assert.Equal(t, true, totalWithdrawl.GreaterThan(decimal.NewFromFloat(0)))

	txSum := totalDeposit.Sub(totalWithdrawl).Sub(totalFees)
	wallet, err := webclient.GetWallet(walletAddress)
	assert.Nil(t, err)

	actual := wallet.GetBalance().StringFixed(8)
	expected := txSum.StringFixed(8)
	assert.Equal(t, actual, expected)

	test.CleanupIntegrationTest()
}

func TestEthereumWebClient_GetToken(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := NewFiatPriceService(ctx, exchangeService)

	localAuthService := NewLocalAuthService(ctx, userDAO, userMapper)
	service, err := NewEthereumWebClient(ctx, userDAO, localAuthService, NewMarketCapService(ctx), fiatPriceService)
	assert.Nil(t, err)

	token, err := service.GetToken(os.Getenv("ETH_ADDRESS"), os.Getenv("TOKEN_ADDRESS"))
	assert.Nil(t, err)

	assert.Equal(t, true, token.GetBalance().GreaterThan(decimal.NewFromFloat(0)))
	//fmt.Printf("balance: %.8f\n", balance)

	test.CleanupIntegrationTest()
}

func TestEthereumWebClientService_GetTokenTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()
	pluginService := NewPluginService(ctx, pluginDAO, pluginMapper)
	exchangeService := NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	fiatPriceService, err := NewFiatPriceService(ctx, exchangeService)
	assert.Nil(t, err)

	localAuthService := NewLocalAuthService(ctx, userDAO, userMapper)
	service, err := NewEthereumWebClient(ctx, userDAO, localAuthService, NewMarketCapService(ctx), fiatPriceService)
	assert.Nil(t, err)

	transactions, err := service.GetTokenTransactions(os.Getenv("TOKEN_INTERNAL_ADDRESS"))

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	for _, tx := range transactions {
		util.DUMP(tx)
	}

	test.CleanupIntegrationTest()
}
