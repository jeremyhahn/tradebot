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
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT *common.Context
var TEST_LOCK sync.Mutex

func NewUnitTestContext() *common.Context {
	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)
	return &common.Context{
		Logger: logger,
		User: &dto.UserDTO{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}
}

func NewIntegrationTestContext() *common.Context {

	TEST_LOCK.Lock()

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
		User: &dto.UserDTO{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}

	var wallets []entity.UserWallet
	wallets = append(wallets, entity.UserWallet{
		Currency: "BTC",
		Address:  BTC_ADDRESS})
	wallets = append(wallets, entity.UserWallet{
		Currency: "XRP",
		Address:  XRP_ADDRESS})

	var exchanges []entity.UserCryptoExchange
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "gdax",
		Key:    GDAX_APIKEY,
		Secret: GDAX_SECRET,
		Extra:  GDAX_PASSPHRASE})
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "bittrex",
		Key:    BITTREX_APIKEY,
		Secret: BITTREX_SECRET,
		Extra:  BITTREX_EXTRA})
	exchanges = append(exchanges, entity.UserCryptoExchange{
		Name:   "binance",
		Key:    BINANCE_APIKEY,
		Secret: BINANCE_SECRET,
		Extra:  BINANCE_EXTRA})
	/*exchanges = append(exchanges, entity.UserCryptoExchange{
	Name:   "bithumb",
	Key:    BITHUMB_APIKEY,
	Secret: BITHUMB_SECRET})*/

	userDAO := dao.NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&entity.User{Username: TEST_USERNAME, LocalCurrency: "USD", Exchanges: exchanges, Wallets: wallets})

	/*exchangeDAO := exchange.NewExchangeDAO(TEST_CONTEXT)
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:       "gdax",
		Key:        GDAX_APIKEY,
		Secret:     GDAX_SECRET,
		Passphrase: GDAX_PASSPHRASE})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bittrex",
		Key:    BITTREX_APIKEY,
		Secret: BITTREX_SECRET})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "binance",
		Key:    BINANCE_APIKEY,
		Secret: BINANCE_SECRET})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bithumb",
		Key:    BITHUMB_APIKEY,
		Secret: BITHUMB_SECRET})

	userDAO.Create(&entity.User{
		Id: TEST_CONTEXT.User
		Exchanges: . exchange.CryptoExchange{
		Name:       "gdax",
		Key:        GDAX_APIKEY,
		Secret:     GDAX_SECRET,
		Passphrase: GDAX_PASSPHRASE})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bittrex",
		Key:    BITTREX_APIKEY,
		Secret: BITTREX_SECRET})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "binance",
		Key:    BINANCE_APIKEY,
		Secret: BINANCE_SECRET})
	exchangeDAO.Create(&exchange.CryptoExchange{
		Name:   "bithumb",
		Key:    BITHUMB_APIKEY,
		Secret: BITHUMB_SECRET})*/

	return TEST_CONTEXT
}

func CleanupIntegrationTest() {
	if TEST_CONTEXT != nil {
		TEST_CONTEXT.DB.Close()
		os.Remove(TEST_DBPATH)
		TEST_LOCK.Unlock()
	}
}
