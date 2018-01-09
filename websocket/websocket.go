package websocket

import (
	"fmt"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

type WebsocketRequest struct {
	Exchange     string `json:"exchange"`
	CurrencyPair *common.CurrencyPair
}

type WebsocketServer struct {
	ctx              *common.Context
	port             int
	closeChan        chan bool
	chartHandler     *ChartHandler
	portfolioHandler *PortfolioHandler
	marketcapHandler *MarketCapHandler
	marketcapService *service.MarketCapService
}

func NewWebsocketServer(ctx *common.Context, port int, marketcapService *service.MarketCapService) *WebsocketServer {
	return &WebsocketServer{
		ctx:              ctx,
		port:             port,
		marketcapService: marketcapService}
}

func (ws *WebsocketServer) Start() {

	ws.ctx.Logger.Debug("[WebSocket] Starting file service on port: ", ws.port)
	http.Handle("/", http.FileServer(http.Dir("webui/public")))

	ws.ctx.Logger.Debug("[WebSocket] Starting websocket service on port: ", ws.port)

	http.HandleFunc("/chart", ws.chartHandler.onConnect)

	http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		portfolioHub := NewPortfolioHub(ws.ctx.Logger)
		go portfolioHub.Run()
		ph := NewPortfolioHandler(ws.ctx, portfolioHub, ws.marketcapService)
		ph.onConnect(w, r)
	})

	http.HandleFunc("/marketcap", ws.marketcapHandler.onConnect)

	go ws.run()

	err := http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
	if err != nil {
		ws.ctx.Logger.Error("[WebSocket] Unable to start server: ", err)
	}
}

func (ws *WebsocketServer) run() {
	ws.ctx.Logger.Debug("[WebsocketServer.run] Starting loop")
	for {
		ws.ctx.Logger.Debug("[WebsocketServer.run] Main loop...")
		select {
		case close := <-ws.closeChan:
			if close {
				ws.ctx.Logger.Debug("[WebsocketServer.run] Stopping websocket server")
				break
			}
		}
	}
}

func (ws *WebsocketServer) Stop() {
	ws.closeChan <- true
}
