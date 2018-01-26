package test

import (
	"os"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT *common.Context

func NewUnitTestContext() *common.Context {
	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)
	return &common.Context{
		Logger: logger,
		User: &common.User{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}
}

func NewIntegrationTestContext() *common.Context {

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	db, err := gorm.Open("sqlite3", TEST_DBPATH)
	db.LogMode(true)
	if err != nil {
		panic(err)
	}

	TEST_CONTEXT = &common.Context{
		DB:     db,
		Logger: logger,
		User: &common.User{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}

	var wallets []dao.UserWallet
	wallets = append(wallets, dao.UserWallet{
		Currency: "BTC",
		Address:  BTC_ADDRESS})
	wallets = append(wallets, dao.UserWallet{
		Currency: "XRP",
		Address:  XRP_ADDRESS})

	var exchanges []dao.UserCoinExchange
	exchanges = append(exchanges, dao.UserCoinExchange{
		Name:   "gdax",
		Key:    GDAX_APIKEY,
		Secret: GDAX_SECRET,
		Extra:  GDAX_PASSPHRASE})
	exchanges = append(exchanges, dao.UserCoinExchange{
		Name:   "bittrex",
		Key:    BITTREX_APIKEY,
		Secret: BITTREX_SECRET,
		Extra:  BITTREX_EXTRA})
	exchanges = append(exchanges, dao.UserCoinExchange{
		Name:   "binance",
		Key:    BINANCE_APIKEY,
		Secret: BINANCE_SECRET,
		Extra:  BINANCE_EXTRA})
	/*exchanges = append(exchanges, dao.UserCoinExchange{
	Name:   "bithumb",
	Key:    BITHUMB_APIKEY,
	Secret: BITHUMB_SECRET})*/

	userDAO := dao.NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&dao.User{Username: TEST_USERNAME, LocalCurrency: "USD", Exchanges: exchanges, Wallets: wallets})

	/*exchangeDAO := exchange.NewExchangeDAO(TEST_CONTEXT)
	exchangeDAO.Create(&exchange.CoinExchange{
		Name:       "gdax",
		Key:        GDAX_APIKEY,
		Secret:     GDAX_SECRET,
		Passphrase: GDAX_PASSPHRASE})
	exchangeDAO.Create(&exchange.CoinExchange{
		Name:   "bittrex",
		Key:    BITTREX_APIKEY,
		Secret: BITTREX_SECRET})
	exchangeDAO.Create(&exchange.CoinExchange{
		Name:   "binance",
		Key:    BINANCE_APIKEY,
		Secret: BINANCE_SECRET})
	exchangeDAO.Create(&exchange.CoinExchange{
		Name:   "bithumb",
		Key:    BITHUMB_APIKEY,
		Secret: BITHUMB_SECRET})

	userDAO.Create(&dao.User{
		Id: TEST_CONTEXT.User
		Exchanges: . exchange.CoinExchange{
		Name:       "gdax",
		Key:        GDAX_APIKEY,
		Secret:     GDAX_SECRET,
		Passphrase: GDAX_PASSPHRASE})
	exchangeDAO.Create(&exchange.CoinExchange{
		Name:   "bittrex",
		Key:    BITTREX_APIKEY,
		Secret: BITTREX_SECRET})
	exchangeDAO.Create(&exchange.CoinExchange{
		Name:   "binance",
		Key:    BINANCE_APIKEY,
		Secret: BINANCE_SECRET})
	exchangeDAO.Create(&exchange.CoinExchange{
		Name:   "bithumb",
		Key:    BITHUMB_APIKEY,
		Secret: BITHUMB_SECRET})*/

	return TEST_CONTEXT
}

func CleanupMockContext() {
	if TEST_CONTEXT != nil {
		TEST_CONTEXT.DB.Close()
		os.Remove(TEST_DBPATH)
	}
}
