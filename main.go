package main

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/webservice"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

func main() {

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	sqlite := InitSQLite()
	defer sqlite.Close()

	ctx := &common.Context{
		DB:     sqlite,
		Logger: logger}

	userDAO := dao.NewUserDAO(ctx)
	ctx.User = userDAO.GetById(1)

	chartDAO := dao.NewChartDAO(ctx)
	indicatorDAO := dao.NewIndicatorDAO(ctx)
	chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
	tradeDAO := dao.NewTradeDAO(ctx)
	profitDAO := dao.NewProfitDAO(ctx)

	indicatorMapper := mapper.NewIndicatorMapper()

	marketcapService := service.NewMarketCapService(logger)
	exchangeService := service.NewExchangeService(ctx, dao.NewExchangeDAO(ctx))
	pluginService := service.NewPluginService(ctx)
	indicatorService := service.NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)
	chartService := service.NewChartService(ctx, chartDAO, exchangeService, indicatorService)
	tradeService := service.NewTradeService(ctx, tradeDAO)
	profitService := service.NewProfitService(ctx, profitDAO)

	autoTradeService := service.NewAutoTradeService(ctx, exchangeService, chartService,
		profitService, tradeService, pluginService)
	err := autoTradeService.EndWorldHunger()
	if err != nil {
		ctx.Logger.Error(err.Error())
	}

	ws := webservice.NewWebServer(ctx, 8080, marketcapService, exchangeService)
	go ws.Start()
	ws.Run()
}

func InitSQLite() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./db/tradebot.db")
	db.LogMode(true)
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
