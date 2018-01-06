package exchange

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	logging "github.com/op/go-logging"
)

var SupportedExchangeMap = map[string]func(*dao.UserCoinExchange, *logging.Logger, *common.CurrencyPair) common.Exchange{
	"gdax":    NewGDAX,
	"bittrex": NewBittrex,
	"binance": NewBinance,
	"bithumb": nil}

var CurrencyPairMap = map[string]*common.CurrencyPair{
	"gdax":    &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"},
	"bittrex": &common.CurrencyPair{Base: "USDT", Quote: "BTC", LocalCurrency: "USDT"},
	"binance": &common.CurrencyPair{Base: "BTC", Quote: "USDT", LocalCurrency: "USDT"},
	"bithumb": nil}
