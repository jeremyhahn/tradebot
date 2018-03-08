// +build integration

package dao

import (
	"sync"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT common.Context
var TEST_LOCK sync.Mutex
var TEST_USERNAME = "test"

var database = common.CreateDatabase("/tmp", "dao-", true)

func NewIntegrationTestContext() common.Context {

	TEST_LOCK.Lock()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	database.MigrateCoreDB()
	database.MigratePriceDB()

	TEST_CONTEXT = &common.Ctx{
		AppRoot: "../",
		CoreDB:  database.ConnectCoreDB(),
		PriceDB: database.ConnectPriceDB(),
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
