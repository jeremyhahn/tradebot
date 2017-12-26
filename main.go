package main

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

const (
	APPNAME    = "tradebot"
	APPVERSION = "0.0.1"
)

func getExchangeList(exchanges *CoinExchanges, coinbase *Coinbase, bittrex *Bittrex) []common.CoinExchange {
	var exchangeList []common.CoinExchange
	for _, ex := range exchanges.Exchanges {
		if ex.Name == "coinbase" {
			total := 0.0
			balances := coinbase.GetBalances()
			for _, c := range balances {
				total += c.Total
			}
			exchangeList = append(exchangeList, common.CoinExchange{
				Name:  ex.Name,
				URL:   ex.URL,
				Total: total,
				Coins: balances})
		} else if ex.Name == "bittrex" {
			total := 0.0
			balances := bittrex.GetBalances()
			for _, c := range balances {
				total += c.Total
			}
			exchangeList = append(exchangeList, common.CoinExchange{
				Name:  ex.Name,
				URL:   ex.URL,
				Total: total,
				Coins: balances})
		}
	}
	return exchangeList
}

func main() {

	backend, _ := logging.NewSyslogBackend(APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(APPNAME)

	sqlite := InitSQLite()
	defer sqlite.Close()

	//mysql := InitMySQL()
	//defer mysql.Close()

	//config := NewConfiguration(sqlite, logger)

	period := 900 // seconds; 15 minutes
	priceStream := NewPriceStream(period)

	exchanges := NewCoinExchanges(sqlite, logger)

	coinbase := NewCoinbase(exchanges.Get("coinbase"), logger, priceStream)
	bittrex := NewBittrex(exchanges.Get("bittrex"), logger, priceStream)

	//btcChart := NewChart(sqlite, coinbase, logger, priceStream)
	charts := make([]*Chart, 0)
	//charts = append(charts, btcChart)
	ws := NewWebsocketServer(8080, charts, logger)
	go ws.Start()

	for {
		logger.Info("Looping...")
		exchangeList := getExchangeList(exchanges, coinbase, bittrex)
		ws.BroadcastPortfolio(exchangeList)
		time.Sleep(5 * time.Second)
	}

	//btcChart.Stream(ws)

	/*
		bittrex := exchanges.Get("bittrex")

		ada := NewBittrex(&bittrex, logger, "BTC-ADA", btcTicker)

		ada.GetBalances()

		os.Exit(1)
	*/

	/*
			markets, err := ada.client.GetCurrencies()
			if err != nil {
				ada.logger.Error(err)
			}
			data, _ := json.Marshal(markets)

		marketSummary, _ := bittrex.client.GetMarketSummary("BTC-ADA")
		data, _ := json.Marshal(marketSummary)
		fmt.Print(string(data))
	*/

	/*
		adaChart := NewChart(sqlite, ada, logger)
		charts := make([]*Chart, 0)
		charts = append(charts, adaChart)
	*/

	/*
		btc := NewCoinbase(config, logger, "BTC-USD")
		eth := NewCoinbase(config, logger, "ETH-USD")
		ltc := NewCoinbase(config, logger, "LTC-USD")

		//btcChart := NewChart(mysql, btc, logger)
		//ethChart := NewChart(mysql, eth, logger)
		//ltcChart := NewChart(mysql, ltc, logger)

		btcChart := NewChart(sqlite, btc, logger)
		ethChart := NewChart(sqlite, eth, logger)
		ltcChart := NewChart(sqlite, ltc, logger)

		charts := make([]*Chart, 0)
		charts = append(charts, btcChart)
		charts = append(charts, ethChart)
		charts = append(charts, ltcChart)
	*/
	//ws := NewWebsocketServer(8080, charts, logger)
	//priceStream.SubscribeToPrice(ws)

	//go ws.Start()

	//	go btcChart.Stream(ws)
	//	go ethChart.Stream(ws)
	//	go ltcChart.Stream(ws)
	//btcChart.Stream(ws)
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
