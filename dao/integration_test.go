// +build integration

package dao

import (
	"os"
	"sync"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT *common.Context
var TEST_LOCK sync.Mutex
var TEST_USERNAME = "test"
var TEST_COREDB_PATH = "/tmp/tradebot-integration-dao-testing.db"
var TEST_PRICEDB_PATH = "/tmp/tradebot-pricehistory.db"

func NewIntegrationTestContext() *common.Context {

	TEST_LOCK.Lock()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	coredb, err := gorm.Open("sqlite3", TEST_COREDB_PATH)
	coredb.LogMode(true)
	if err != nil {
		panic(err)
	}
	pricedb, err := gorm.Open("sqlite3", TEST_COREDB_PATH)
	pricedb.LogMode(true)
	if err != nil {
		panic(err)
	}

	TEST_CONTEXT = &common.Context{
		CoreDB:  coredb,
		PriceDB: pricedb,
		Logger:  logger,
		User: &dto.UserDTO{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}

	userDAO := NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&entity.User{Username: TEST_USERNAME, LocalCurrency: "USD"})

	return TEST_CONTEXT
}

func CleanupIntegrationTest() {
	if TEST_CONTEXT != nil {
		TEST_CONTEXT.CoreDB.Close()
		TEST_CONTEXT.PriceDB.Close()
		os.Remove(TEST_COREDB_PATH)
		os.Remove(TEST_PRICEDB_PATH)
		TEST_LOCK.Unlock()
		TEST_CONTEXT = nil
	}
}

func createIntegrationTestChart(ctx *common.Context) entity.ChartEntity {
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
			UserId:    ctx.User.GetId(),
			Base:      "BTC",
			Quote:     "USD",
			Exchange:  "Test",
			Date:      time.Now(),
			Type:      "buy",
			Amount:    2,
			Price:     10000,
			ChartData: "test-trade-1"},
		entity.Trade{
			UserId:    ctx.User.GetId(),
			Base:      "BTC",
			Quote:     "USD",
			Exchange:  "Test",
			Date:      time.Now(),
			Type:      "sell",
			Amount:    2,
			Price:     12000,
			ChartData: "test-trade-2"}}

	chart := &entity.Chart{
		UserId:     ctx.User.GetId(),
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
