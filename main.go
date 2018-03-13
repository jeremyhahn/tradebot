package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/webservice"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

func main() {

	wd, _ := os.Getwd()
	defaultIpc := fmt.Sprintf("%s/%s", wd, "test/ethereum/blockchain/geth.ipc")
	defaultKeystore := fmt.Sprintf("%s/%s", wd, "test/ethereum/blockchain/keystore")
	defaultMode := "etherscan"

	initDbFlag := flag.Bool("initdb", false, "Create / migrate database schema and exit")
	portFlag := flag.Int("port", 8080, "Web server listen port")
	sslFlag := flag.Bool("ssl", true, "Enable HTTPS / WSS over TLS")
	ipcFlag := flag.String("ipc", defaultIpc, "Path to geth IPC socket")
	keystoreFlag := flag.String("keystore", defaultKeystore, "Path to default Ethereum keystore")
	debugFlag := flag.Bool("debug", false, "Enable debug level logging")
	ethereumModeFlag := flag.String("mode", defaultMode, "Ethereum mode [native | etherscan] (default: etherscan)")
	flag.Parse()

	f, err := os.OpenFile("tradebot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic("Unable to open log file")
	}
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	logfile := logging.NewLogBackend(f, "", log.Lshortfile)
	//syslog, _ := logging.NewSyslogBackend(common.APPNAME)
	backends := logging.MultiLogger(stdout, logfile)
	logger := logging.MustGetLogger(common.APPNAME)
	logging.SetBackend(backends)
	if *debugFlag == false {
		logging.SetLevel(logging.ERROR, "")
	} else {
		logger.Debug("Starting in debug mode...")
	}

	databaseManager := common.CreateDatabase("./db", "", *debugFlag)

	if *initDbFlag {
		databaseManager.MigrateCoreDB()
		databaseManager.MigratePriceDB()
		os.Exit(0)
	}

	ctx := &common.Ctx{
		AppRoot:      wd,
		CoreDB:       databaseManager.ConnectCoreDB(),
		PriceDB:      databaseManager.ConnectPriceDB(),
		Logger:       logger,
		Debug:        *debugFlag,
		SSL:          *sslFlag,
		IPC:          *ipcFlag,
		Keystore:     *keystoreFlag,
		EthereumMode: *ethereumModeFlag}
	defer ctx.Close()

	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()

	ethereumService, err := service.NewEthereumService(ctx, userDAO, userMapper, service.NewMarketCapService(ctx))
	if err != nil {
		ctx.Logger.Fatalf(fmt.Sprintf("Error: %s", err.Error()))
	}

	jsonWebTokenService, err := service.NewJsonWebTokenService(ctx, databaseManager, ethereumService, common.NewJsonWriter())
	if err != nil {
		ctx.Logger.Fatalf(fmt.Sprintf("Error: %s", err.Error()))
	}

	ws := webservice.NewWebServer(ctx, *portFlag, ethereumService, jsonWebTokenService)

	go ws.Start()
	ws.Run()

	/*
		pluginDAO := dao.NewPluginDAO(ctx)
		chartDAO := dao.NewChartDAO(ctx)
		chartIndicatorDAO := dao.NewChartIndicatorDAO(ctx)
		profitDAO := dao.NewProfitDAO(ctx)
		tradeDAO := dao.NewTradeDAO(ctx)
		chartStrategyDAO := dao.NewChartStrategyDAO(ctx)

		chartMapper := mapper.NewChartMapper(ctx)
		tradeMapper := mapper.NewTradeMapper()
		pluginMapper := mapper.NewPluginMapper()
		userExchangeMapper := mapper.NewUserExchangeMapper()

		exchangeService := service.NewExchangeService(ctx, pluginDAO, userDAO, userMapper, userExchangeMapper)
		pluginService := service.NewPluginService(ctx, pluginDAO, pluginMapper)
		indicatorService := service.NewIndicatorService(ctx, chartIndicatorDAO, pluginService)
		chartService := service.NewChartService(ctx, userDAO, chartDAO, exchangeService, indicatorService)
		profitService := service.NewProfitService(ctx, profitDAO)
		tradeService := service.NewTradeService(ctx, tradeDAO, tradeMapper)
		strategyService := service.NewStrategyService(ctx, chartStrategyDAO, pluginService, indicatorService, chartMapper)
		autoTradeService := service.NewAutoTradeService(ctx, exchangeService, chartService, profitService, tradeService, strategyService, userMapper)

		err = autoTradeService.EndWorldHunger()
		if err != nil {
			ctx.Logger.Fatalf(fmt.Sprintf("Error: %s", err.Error()))
		}
	*/
}
