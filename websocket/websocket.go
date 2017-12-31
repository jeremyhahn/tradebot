package websocket

import (
	"fmt"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type WebsocketRequest struct {
	Exchange     string `json:"exchange"`
	CurrencyPair *common.CurrencyPair
}

type WebsocketServer struct {
	logger           *logging.Logger
	port             int
	running          bool
	chartHandler     *ChartHandler
	portfolioHandler *PortfolioHandler
	PriceChangeChan  chan common.PriceChange
	CandlestickChan  chan common.Candlestick
	ChartChan        chan *common.ChartData
	CloseChan        chan bool
	PortfolioChan    chan *common.Portfolio
}

func NewWebsocketServer(port int, logger *logging.Logger) *WebsocketServer {
	return &WebsocketServer{
		logger:           logger,
		port:             port,
		running:          false,
		chartHandler:     NewChartHandler(logger),
		portfolioHandler: NewPortfolioHandler(logger),
		PriceChangeChan:  make(chan common.PriceChange),
		CandlestickChan:  make(chan common.Candlestick),
		ChartChan:        make(chan *common.ChartData),
		PortfolioChan:    make(chan *common.Portfolio)}
}

func (ws *WebsocketServer) Start() {

	ws.logger.Debug("[WebSocket] Starting file service on port: ", ws.port)
	http.Handle("/", http.FileServer(http.Dir("webui/public")))

	ws.logger.Debug("[WebSocket] Starting websocket service on port: ", ws.port)

	http.HandleFunc("/chart", ws.chartHandler.onConnect)
	http.HandleFunc("/portfolio", ws.portfolioHandler.onConnect)

	err := http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
	if err != nil {
		ws.logger.Error("[WebSocket] Unable to start server: ", err)
	}
}

func (ws *WebsocketServer) Run() {
	ws.logger.Debug("[WebsocketServer.run] Starting loop")
	for {
		ws.logger.Debug("[WebsocketServer.run] Main loop...")
		select {
		case price := <-ws.PriceChangeChan:
			ws.logger.Debugf("[WebsocketServer.run] Broadcasting price: %+v\n: ", price)
		case candle := <-ws.CandlestickChan:
			ws.logger.Debugf("[WebsocketServer.run] Broadcasting candlestick: %+v\n", candle)
		case chart := <-ws.ChartChan:
			ws.logger.Debugf("[WebsocketServer.run] Broadcasting chart: %+v\n", chart)
		case portfolio := <-ws.PortfolioChan:
			ws.logger.Debugf("[WebsocketServer.run] Broadcasting portfolio: %+v\n", portfolio)
			ws.portfolioHandler.Broadcast(portfolio.Exchanges)
		case close := <-ws.CloseChan:
			if close {
				ws.logger.Debug("[WebsocketServer.run] Stopping websocket server")
				break
			}
		}
	}
}

func (ws *WebsocketServer) Stop() {
	ws.running = false
}
