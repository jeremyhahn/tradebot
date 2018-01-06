package websocket

import (
	"fmt"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
)

type WebsocketRequest struct {
	Exchange     string `json:"exchange"`
	CurrencyPair *common.CurrencyPair
}

type WebsocketServer struct {
	ctx              *common.Context
	port             int
	running          bool
	chartHandler     *ChartHandler
	portfolioHandler *PortfolioHandler
	marketcapHandler *MarketCapHandler
	PriceChangeChan  chan common.PriceChange
	CandlestickChan  chan common.Candlestick
	ChartChan        chan *common.ChartData
	CloseChan        chan bool
	PortfolioChan    chan *common.Portfolio
	MarketCapChan    chan *common.MarketCap
}

func NewWebsocketServer(ctx *common.Context, port int) *WebsocketServer {
	portfolioChan := make(chan *common.Portfolio)
	return &WebsocketServer{
		ctx:              ctx,
		port:             port,
		chartHandler:     NewChartHandler(ctx.Logger),
		marketcapHandler: NewMarketCapHandler(ctx.Logger),
		PriceChangeChan:  make(chan common.PriceChange),
		CandlestickChan:  make(chan common.Candlestick),
		ChartChan:        make(chan *common.ChartData),
		PortfolioChan:    portfolioChan,
		MarketCapChan:    make(chan *common.MarketCap)}
}

func (ws *WebsocketServer) Start(portfolioHub *PortfolioHub) {

	ws.ctx.Logger.Debug("[WebSocket] Starting file service on port: ", ws.port)
	http.Handle("/", http.FileServer(http.Dir("webui/public")))

	ws.ctx.Logger.Debug("[WebSocket] Starting websocket service on port: ", ws.port)

	http.HandleFunc("/chart", ws.chartHandler.onConnect)

	http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		ph := NewPortfolioHandler(ws.ctx, portfolioHub)
		ph.onConnect(w, r)
	})

	http.HandleFunc("/marketcap", ws.marketcapHandler.onConnect)

	err := http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
	if err != nil {
		ws.ctx.Logger.Error("[WebSocket] Unable to start server: ", err)
	}
}

func (ws *WebsocketServer) Run() {
	ws.ctx.Logger.Debug("[WebsocketServer.run] Starting loop")
	for {

		ws.ctx.Logger.Debug("[WebsocketServer.run] Main loop...")
		select {

		case price := <-ws.PriceChangeChan:
			ws.ctx.Logger.Debugf("[WebsocketServer.run] Broadcasting price: %+v\n: ", price)

		case candle := <-ws.CandlestickChan:
			ws.ctx.Logger.Debugf("[WebsocketServer.run] Broadcasting candlestick: %+v\n", candle)

		case chart := <-ws.ChartChan:
			ws.ctx.Logger.Debugf("[WebsocketServer.run] Broadcasting chart: %+v\n", chart)

		case marketcap := <-ws.MarketCapChan:
			ws.ctx.Logger.Debugf("[WebsocketServer.run] Broadcasting market cap: %+v\n", marketcap)
			ws.marketcapHandler.Broadcast(marketcap)

		case close := <-ws.CloseChan:
			if close {
				ws.ctx.Logger.Debug("[WebsocketServer.run] Stopping websocket server")
				break
			}
		}
	}
}
