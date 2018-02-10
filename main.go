package main

import (
	"flag"
	"fmt"
	"os"

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

	wd, _ := os.Getwd()
	defaultIpc := fmt.Sprintf("%s/%s", wd, "test/ethereum/blockchain/geth.ipc")
	defaultKeystore := fmt.Sprintf("%s/%s", wd, "test/ethereum/blockchain/keystore/test")

	portFlag := flag.Int("port", 8080, "Web server listen port")
	sslFlag := flag.Bool("ssl", true, "Enable debug level logging")
	ipcFlag := flag.String("ipc", defaultIpc, "Path to geth IPC socket")
	keystoreFlag := flag.String("keystore", defaultKeystore, "Path to Ethereum keystore")
	debugFlag := flag.Bool("debug", false, "Enable debug level logging")
	flag.Parse()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logger := logging.MustGetLogger(common.APPNAME)
	logging.SetBackend(backend)
	if *debugFlag == false {
		logger.Debug("Starting in debug mode...")
		logging.SetLevel(logging.ERROR, "")
	}

	sqlite := InitSQLite(*debugFlag)
	defer sqlite.Close()

	ctx := &common.Context{
		DB:     sqlite,
		Logger: logger,
		Debug:  *debugFlag,
		SSL:    *sslFlag}

	userDAO := dao.NewUserDAO(ctx)
	//users := userDAO.Find()

	userMapper := mapper.NewUserMapper()
	marketcapService := service.NewMarketCapService(logger)
	userService := service.NewUserService(ctx, userDAO, marketcapService, userMapper)

	user, err := userService.GetUserById(1)
	if err != nil {
		logger.Fatal(err.Error())
	}
	ctx.SetUser(user)

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

	exchangeService := service.NewExchangeService(ctx, dao.NewExchangeDAO(ctx))
	pluginService := service.NewPluginService(ctx)
	indicatorService := service.NewIndicatorService(ctx, indicatorDAO, chartIndicatorDAO, pluginService, indicatorMapper)
	chartService := service.NewChartService(ctx, chartDAO, exchangeService, indicatorService)
	profitService := service.NewProfitService(ctx, profitDAO)
	tradeService := service.NewTradeService(ctx, tradeDAO, tradeMapper)
	strategyService := service.NewStrategyService(ctx, strategyDAO, chartStrategyDAO, pluginService, indicatorService, chartMapper, strategyMapper)
	autoTradeService := service.NewAutoTradeService(ctx, exchangeService, chartService, profitService, tradeService, strategyService)

	err = autoTradeService.EndWorldHunger()
	if err != nil {
		ctx.Logger.Fatalf(fmt.Sprintf("Error: %s", err.Error()))
	}

	ethereumService, err := service.NewEthereumService(ctx, *ipcFlag, *keystoreFlag, userDAO)
	if err != nil {
		ctx.Logger.Fatalf(fmt.Sprintf("Error: %s", err.Error()))
	}

	ws := webservice.NewWebServer(ctx, *portFlag, marketcapService, exchangeService, ethereumService, userService)
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
