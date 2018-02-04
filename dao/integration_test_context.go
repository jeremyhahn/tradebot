//// +build integration

package dao

import (
	"os"
	"sync"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT *common.Context
var TEST_LOCK sync.Mutex
var TEST_USERNAME = "test"
var TEST_DBPATH = "/tmp/tradebot-integration-testing.db"

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
		User: &common.User{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}

	userDAO := NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&User{Username: TEST_USERNAME, LocalCurrency: "USD"})

	return TEST_CONTEXT
}

func CleanupIntegrationTest() {
	if TEST_CONTEXT != nil {
		TEST_CONTEXT.DB.Close()
		TEST_LOCK.Unlock()
		os.Remove(TEST_DBPATH)
	}
}

func createIntegrationTestChart(ctx *common.Context) (*Chart, []Indicator, []Trade) {
	indicators := []Indicator{
		Indicator{
			Name:       "RelativeStrengthIndex",
			Parameters: "14,70,30"},
		Indicator{
			Name:       "BollingerBands",
			Parameters: "20,2"}}
	trades := []Trade{
		Trade{
			UserId:    ctx.User.Id,
			Base:      "BTC",
			Quote:     "USD",
			Exchange:  "Test",
			Date:      time.Now(),
			Type:      "buy",
			Amount:    2,
			Price:     10000,
			ChartData: "test-trade-1"},
		Trade{
			UserId:    ctx.User.Id,
			Base:      "BTC",
			Quote:     "USD",
			Exchange:  "Test",
			Date:      time.Now(),
			Type:      "sell",
			Amount:    2,
			Price:     12000,
			ChartData: "test-trade-2"}}
	chart := &Chart{
		UserId:     ctx.User.Id,
		Base:       "BTC",
		Quote:      "USD",
		Exchange:   "gdax",
		Period:     900, // 15 minutes
		Indicators: indicators,
		Trades:     trades,
		AutoTrade:  1}
	return chart, indicators, trades
}
