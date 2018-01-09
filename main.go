package main

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/websocket"
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

	//mysql := InitMySQL()
	//defer mysql.Close()

	ctx := &common.Context{
		DB:     sqlite,
		Logger: logger}

	userDAO := dao.NewUserDAO(ctx)
	ctx.User = userDAO.GetById(1)
	/*if user.Username == "" {
		userDAO.Create(&dao.User{
			Username: "test"})
	}*/

	ws := websocket.NewWebsocketServer(ctx, 8080, service.NewMarketCapService(logger))
	go ws.Start()

	exchangeDAO := dao.NewExchangeDAO(ctx)
	autotradeDAO := dao.NewAutoTradeDAO(ctx)
	fmt.Println(autotradeDAO.Find(ctx.User))

	for _, trade := range autotradeDAO.Find(ctx.User) {
		currencyPair := &common.CurrencyPair{
			Base:          trade.Base,
			Quote:         trade.Quote,
			LocalCurrency: ctx.User.LocalCurrency}

		exchangeService := service.NewExchangeService(ctx, exchangeDAO)
		exchange := exchangeService.NewExchange(ctx.User, trade.Exchange, currencyPair)
		chart := service.NewChart(ctx, exchange, trade.Period)

		fmt.Printf("Loading autotrade pair: %s-%s\n", trade.Base, trade.Quote)
		fmt.Printf("Chart: %+v\n", chart)

		chart.Stream()
	}

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
