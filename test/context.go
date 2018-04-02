package test

import (
	"fmt"
	"os"
	"sync"

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

var database = common.CreateDatabase("/tmp", "test-", true)

func NewUnitTestContext() common.Context {
	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)
	return &common.Ctx{
		Logger: logger,
		User: &dto.UserDTO{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD",
			FiatExchange:  "GDAX"}}
}

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
	/*exchanges = append(exchanges, entity.UserCryptoExchange{
	Name:   "bithumb",
	Key:     os.Getenv("BITHUMB_APIKEY"),
	Secret: os.Getenv("BINANCE_SECRET")})*/

	userDAO := dao.NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&entity.User{Username: TEST_USERNAME, LocalCurrency: "USD", Exchanges: exchanges, Wallets: wallets})

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
