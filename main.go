package main

import (
	"flag"
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/webservice"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

func main() {

	debugFlag := flag.Bool("debug", false, "Enable debug level logging")
	flag.Parse()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	if *debugFlag == false {
		logging.SetLevel(logging.ERROR, "")
	}
	logger := logging.MustGetLogger(common.APPNAME)
	if *debugFlag == true {
		logger.Debug("Starting in debug mode...")
	}

	sqlite := InitSQLite(*debugFlag)
	defer sqlite.Close()

	ctx := &common.Context{
		DB:        sqlite,
		Logger:    logger,
		DebugMode: *debugFlag}

	userDAO := dao.NewUserDAO(ctx)
	ctx.User = userDAO.GetById(1)

	chartDAO := dao.NewChartDAO(ctx)
	indicatorDAO := dao.NewIndicatorDAO(ctx)
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	profitDAO := dao.NewProfitDAO(ctx)
	tradeDAO := dao.NewTradeDAO(ctx)
	strategyDAO := dao.NewStrategyDAO(ctx)
	chartStrategyDAO := dao.NewChartStrategyDAO(ctx)

	chartMapper := mapper.NewChartMapper(ctx)
	indicatorMapper := mapper.NewIndicatorMapper()
	strategyMapper := mapper.NewStrategyMapper()
	tradeMapper := mapper.NewTradeMapper()

	marketcapService := service.NewMarketCapService(logger)
	exchangeService := service.NewExchangeService(ctx, dao.NewExchangeDAO(ctx))
	pluginService := service.NewPluginService(ctx)
	indicatorService := service.NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO,
		pluginService, indicatorMapper)
	chartService := service.NewChartService(ctx, chartDAO, exchangeService, indicatorService)
	profitService := service.NewProfitService(ctx, profitDAO)
	tradeService := service.NewTradeService(ctx, tradeDAO, tradeMapper)
	strategyService := service.NewStrategyService(ctx, strategyDAO, chartStrategyDAO,
		pluginService, indicatorService, chartMapper, strategyMapper)
	autoTradeService := service.NewAutoTradeService(ctx, exchangeService, chartService,
		profitService, tradeService, strategyService)

	err := autoTradeService.EndWorldHunger()
	if err != nil {
		ctx.Logger.Errorf(fmt.Sprintf("Error: %s", err.Error()))
	}

	ws := webservice.NewWebServer(ctx, 8080, marketcapService, exchangeService)
	go ws.Start()
	ws.Run()
}

func InitSQLite(logMode bool) *gorm.DB {
	db, err := gorm.Open("sqlite3", "./db/tradebot.db")
	db.LogMode(logMode)
	if err != nil {
		panic(err)
	}
	return db
}

/*
func InitMySQL() *gorm.DB {
	db, err := gorm.Open("mysql", "user:pass@tcp(ip:3306)/mydb?charset=utf8&parseTime=True")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}
*/
