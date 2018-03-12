package service

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT common.Context
var TEST_LOCK sync.Mutex
var TEST_USERNAME = "test"

var database = common.CreateDatabase("/tmp", "service-", true)

func NewIntegrationTestContext() common.Context {

	TEST_LOCK.Lock()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading test environment from .env")
	}

	appRoot := "../"
	database.MigrateCoreDB()
	database.MigratePriceDB()

	TEST_CONTEXT = &common.Ctx{
		AppRoot: appRoot,
		CoreDB:  database.ConnectCoreDB(),
		PriceDB: database.ConnectPriceDB(),
		Logger:  logger,
		Debug:   true,
		SSL:     true,
		User: &dto.UserDTO{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"},
		IPC:      fmt.Sprintf("%stest/ethereum/blockchain/geth.ipc", appRoot),
		Keystore: fmt.Sprintf("%stest/ethereum/blockchain/keystore", appRoot)}

	var wallets []entity.UserWallet
	wallets = append(wallets, entity.UserWallet{
		Currency: "BTC",
		Address:  os.Getenv("BTC_ADDRESS")})
	wallets = append(wallets, entity.UserWallet{
		Currency: "XRP",
		Address:  os.Getenv("XRP_ADDRESS")})

	var exchanges []entity.UserCryptoExchange
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "gdax",
		Key:    os.Getenv("GDAX_APIKEY"),
		Secret: os.Getenv("GDAX_SECRET"),
		Extra:  os.Getenv("GDAX_PASSPHRASE")})
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "bittrex",
		Key:    os.Getenv("BITTREX_APIKEY"),
		Secret: os.Getenv("BITTREX_SECRET"),
		Extra:  os.Getenv("BITTREX_EXTRA")})
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "binance",
		Key:    os.Getenv("BINANCE_APIKEY"),
		Secret: os.Getenv("BINANCE_SECRET"),
		Extra:  os.Getenv("BINANCE_EXTRA")})

	userDAO := dao.NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&entity.User{Username: TEST_USERNAME, LocalCurrency: "USD", Exchanges: exchanges, Wallets: wallets})

	return TEST_CONTEXT
}

func CleanupIntegrationTest() {
	if TEST_CONTEXT != nil {
		database.Close(TEST_CONTEXT.GetCoreDB())
		database.Close(TEST_CONTEXT.GetPriceDB())
		database.DropCoreDB()
		database.DropPriceDB()
		TEST_LOCK.Unlock()
	}
}

func createIntegrationTestChart(ctx common.Context) entity.ChartEntity {
	userIndicators := []entity.ChartIndicator{
		entity.ChartIndicator{
			Name:       "RelativeStrengthIndex",
			Parameters: "14,70,30"},
		entity.ChartIndicator{
			Name:       "BollingerBands",
			Parameters: "20,2"},
		entity.ChartIndicator{
			Name:       "MovingAverageConvergenceDivergence",
			Parameters: "12,26,9"}}

	userStrategies := []entity.ChartStrategy{
		entity.ChartStrategy{
			ChartId:    1,
			Name:       "DefaultTradingStrategy",
			Parameters: "1,2,3"}}

	trades := []entity.Trade{
		entity.Trade{
			UserId:    ctx.GetUser().GetId(),
			Base:      "BTC",
			Quote:     "USD",
			Exchange:  "Test",
			Date:      time.Now(),
			Type:      "buy",
			Amount:    2,
			Price:     10000,
			ChartData: "test-trade-1"},
		entity.Trade{
			UserId:    ctx.GetUser().GetId(),
			Base:      "BTC",
			Quote:     "USD",
			Exchange:  "Test",
			Date:      time.Now(),
			Type:      "sell",
			Amount:    2,
			Price:     12000,
			ChartData: "test-trade-2"}}

	chart := &entity.Chart{
		UserId:     ctx.GetUser().GetId(),
		Base:       "BTC",
		Quote:      "USD",
		Exchange:   "gdax",
		Period:     900, // 15 minutes
		Indicators: userIndicators,
		Strategies: userStrategies,
		Trades:     trades,
		AutoTrade:  1}

	return chart
}

func createIntegrationTestCandles() []common.Candlestick {
	var candles []common.Candlestick
	candles = append(candles, common.Candlestick{Close: 100.00})
	candles = append(candles, common.Candlestick{Close: 200.00})
	candles = append(candles, common.Candlestick{Close: 300.00})
	candles = append(candles, common.Candlestick{Close: 400.00})
	candles = append(candles, common.Candlestick{Close: 500.00})
	candles = append(candles, common.Candlestick{Close: 600.00})
	candles = append(candles, common.Candlestick{Close: 700.00})
	candles = append(candles, common.Candlestick{Close: 800.00})
	candles = append(candles, common.Candlestick{Close: 900.00})
	candles = append(candles, common.Candlestick{Close: 1000.00})
	candles = append(candles, common.Candlestick{Close: 1100.00})
	candles = append(candles, common.Candlestick{Close: 1200.00})
	candles = append(candles, common.Candlestick{Close: 1300.00})
	candles = append(candles, common.Candlestick{Close: 1400.00})
	candles = append(candles, common.Candlestick{Close: 1500.00})
	candles = append(candles, common.Candlestick{Close: 1600.00})
	candles = append(candles, common.Candlestick{Close: 1700.00})
	candles = append(candles, common.Candlestick{Close: 1800.00})
	candles = append(candles, common.Candlestick{Close: 1900.00})
	candles = append(candles, common.Candlestick{Close: 2000.00})
	candles = append(candles, common.Candlestick{Close: 2100.00})
	candles = append(candles, common.Candlestick{Close: 2200.00})
	candles = append(candles, common.Candlestick{Close: 2300.00})
	candles = append(candles, common.Candlestick{Close: 2400.00})
	candles = append(candles, common.Candlestick{Close: 2500.00})
	candles = append(candles, common.Candlestick{Close: 2600.00})
	candles = append(candles, common.Candlestick{Close: 2700.00})
	candles = append(candles, common.Candlestick{Close: 2800.00})
	candles = append(candles, common.Candlestick{Close: 2900.00})
	candles = append(candles, common.Candlestick{Close: 3000.00})
	candles = append(candles, common.Candlestick{Close: 3200.00})
	candles = append(candles, common.Candlestick{Close: 3300.00})
	candles = append(candles, common.Candlestick{Close: 3400.00})
	candles = append(candles, common.Candlestick{Close: 3500.00})
	candles = append(candles, common.Candlestick{Close: 3600.00})
	candles = append(candles, common.Candlestick{Close: 3700.00})
	candles = append(candles, common.Candlestick{Close: 3800.00})
	candles = append(candles, common.Candlestick{Close: 3900.00})
	candles = append(candles, common.Candlestick{Close: 4000.00})
	return candles
}
