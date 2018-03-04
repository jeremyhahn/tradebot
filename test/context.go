package test

import (
	"os"
	"sync"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT common.Context
var TEST_LOCK sync.Mutex
var TEST_USERNAME = "test"
var TEST_COREDBPATH = "/tmp/tradebot-coredb-testing.db"
var TEST_PRICEDBPATH = "/tmp/tradebot-pricedb-testing.db"

func NewUnitTestContext() common.Context {
	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)
	return &common.Ctx{
		Logger: logger,
		User: &dto.UserDTO{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}
}

func NewIntegrationTestContext() common.Context {

	TEST_LOCK.Lock()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	coreDB, err := gorm.Open("sqlite3", TEST_COREDBPATH)
	coreDB.LogMode(true)
	if err != nil {
		panic(err)
	}

	priceDB, err := gorm.Open("sqlite3", TEST_PRICEDBPATH)
	priceDB.LogMode(true)
	if err != nil {
		panic(err)
	}

	err = godotenv.Load("../.env")
	if err != nil {
		panic("Error loading test environment from .env")
	}

	if address := os.Getenv("BTC_ADDRESS"); address == "" {
		panic("Unable to load BTC_ADDRESS environment variable")
	}

	TEST_CONTEXT = &common.Ctx{
		Logger: logger,
		User: &dto.UserDTO{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}

	dao.NewMigrator(TEST_CONTEXT).MigrateCoreDB()
	dao.NewMigrator(TEST_CONTEXT).MigratePriceDB()

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
	/*exchanges = append(exchanges, entity.UserCryptoExchange{
	Name:   "bithumb",
	Key:     os.Getenv("BITHUMB_APIKEY"),
	Secret: os.Getenv("BINANCE_SECRET")})*/

	userDAO := dao.NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&entity.User{Username: TEST_USERNAME, LocalCurrency: "USD", Exchanges: exchanges, Wallets: wallets})

	/*exchangeDAO := exchange.NewExchangeDAO(TEST_CONTEXT)
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:       "gdax",
		Key:        os.Getenv("GDAX_APIKEY"),
		Secret:     os.Getenv("GDAX_SECRET"),
		Passphrase: os.Getenv("GDAX_PASSPHRASE")})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bittrex",
		Key:    os.Getenv("BITTREX_APIKEY"),
		Secret: os.Getenv("BITTREX_SECRET")})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "binance",
		Key:    os.Getenv("BINANCE_APIKEY"),
		Secret: os.Getenv("BINANCE_SECRET")})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bithumb",
		Key:    os.Getenv("BITHUMB_APIKEY"),
		Secret: os.Getenv("BINANCE_SECRET")})

	userDAO.Create(&entity.User{
		Id: TEST_CONTEXT.User
		Exchanges: . exchange.CryptoExchange{
		Name:       "gdax",
		Key:        os.Getenv("GDAX_APIKEY"),
		Secret:     os.Getenv("GDAX_SECRET"),
		Passphrase: os.Getenv("GDAX_PASSPHRASE")})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bittrex",
		Key:    os.Getenv("BITTREX_APIKEY"),
		Secret: os.Getenv("BITTREX_SECRET")})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "binance",
		Key:    os.Getenv("BINANCE_APIKEY"),
		Secret: os.Getenv("BINANCE_SECRET")})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bithumb",
		Key:    os.Getenv("BITHUMB_APIKEY"),
		Secret: os.Getenv("BITHUMB_SECRET")})*/

	return TEST_CONTEXT
}

func CleanupIntegrationTest() {
	if TEST_CONTEXT != nil {
		TEST_CONTEXT.GetCoreDB().Close()
		TEST_CONTEXT.GetPriceDB().Close()
		os.Remove(TEST_COREDBPATH)
		os.Remove(TEST_PRICEDBPATH)
		TEST_LOCK.Unlock()
	}
}
