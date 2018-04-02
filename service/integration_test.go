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
	"github.com/shopspring/decimal"
)

var TEST_CONTEXT common.Context
var TEST_LOCK sync.Mutex
var TEST_USERNAME = "test"

var database = common.CreateDatabase("/tmp", "service-", true)

func NewIntegrationTestContext() common.Context {
	return CreateIntegrationTestContext("../.env", "../")
}

func CreateIntegrationTestContext(dotEnvDir, appRoot string) common.Context {

	TEST_LOCK.Lock()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	err := godotenv.Load(dotEnvDir)
	if err != nil {
		panic("Error loading test environment from .env")
	}

	if address := os.Getenv("BTC_ADDRESS"); address == "" {
		panic("Unable to load BTC_ADDRESS environment variable")
	}

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

	var userWwallets []entity.UserWallet
	userWwallets = append(userWwallets, entity.UserWallet{
		Currency: "BTC",
		Address:  os.Getenv("BTC_ADDRESS")})
	userWwallets = append(userWwallets, entity.UserWallet{
		Currency: "XRP",
		Address:  os.Getenv("XRP_ADDRESS")})

	var exchanges []entity.UserCryptoExchange
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "Coinbase",
		Key:    os.Getenv("COINBASE_APIKEY"),
		Secret: os.Getenv("COINBASE_SECRET")})
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "GDAX",
		Key:    os.Getenv("GDAX_APIKEY"),
		Secret: os.Getenv("GDAX_SECRET"),
		Extra:  os.Getenv("GDAX_PASSPHRASE")})
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "Bittrex",
		Key:    os.Getenv("BITTREX_APIKEY"),
		Secret: os.Getenv("BITTREX_SECRET"),
		Extra:  os.Getenv("BITTREX_EXTRA")})
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "Binance",
		Key:    os.Getenv("BINANCE_APIKEY"),
		Secret: os.Getenv("BINANCE_SECRET"),
		Extra:  os.Getenv("BINANCE_EXTRA")})

	userDAO := dao.NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&entity.User{Username: TEST_USERNAME, LocalCurrency: "USD", Exchanges: exchanges, Wallets: userWwallets})

	pluginDAO := dao.NewPluginDAO(TEST_CONTEXT)
	pluginDAO.Create(&entity.Plugin{
		Name:     "GDAX",
		Filename: "gdax.so",
		Version:  "0.0.1a",
		Type:     common.EXCHANGE_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "Coinbase",
		Filename: "coinbase.so",
		Version:  "0.0.1a",
		Type:     common.EXCHANGE_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "Bittrex",
		Filename: "bittrex.so",
		Version:  "0.0.1a",
		Type:     common.EXCHANGE_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "Binance",
		Filename: "binance.so",
		Version:  "0.0.1a",
		Type:     common.EXCHANGE_PLUGIN_TYPE})

	pluginDAO.Create(&entity.Plugin{
		Name:     "BTC",
		Filename: "btc.so",
		Version:  "0.0.1a",
		Type:     common.WALLET_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "ETH",
		Filename: "eth.so",
		Version:  "0.0.1a",
		Type:     common.WALLET_PLUGIN_TYPE})
	pluginDAO.Create(&entity.Plugin{
		Name:     "XRP",
		Filename: "xrp.so",
		Version:  "0.0.1a",
		Type:     common.WALLET_PLUGIN_TYPE})

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
			Amount:    "2",
			Price:     "10000",
			ChartData: "test-trade-1"},
		entity.Trade{
			UserId:    ctx.GetUser().GetId(),
			Base:      "BTC",
			Quote:     "USD",
			Exchange:  "Test",
			Date:      time.Now(),
			Type:      "sell",
			Amount:    "2",
			Price:     "12000",
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
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(100.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1000.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1100.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2000.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2100.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3000.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(4000.00)})
	return candles
}
