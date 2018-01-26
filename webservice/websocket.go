package webservice

import (
	"fmt"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/restapi"
	"github.com/jeremyhahn/tradebot/service"
)

type WebRequest struct {
	Exchange     string `json:"exchange"`
	CurrencyPair *common.CurrencyPair
}

type WebServer struct {
	ctx              *common.Context
	port             int
	closeChan        chan bool
	portfolioHandler *PortfolioHandler
	marketcapService *service.MarketCapService
	exchangeService  service.ExchangeService
}

func NewWebServer(ctx *common.Context, port int, marketcapService *service.MarketCapService,
	exchangeService service.ExchangeService) *WebServer {
	return &WebServer{
		ctx:              ctx,
		port:             port,
		marketcapService: marketcapService,
		exchangeService:  exchangeService}
}

func (ws *WebServer) Start() {

	ws.ctx.Logger.Debug("[Web] Starting web services on port: ", ws.port)
	http.Handle("/", http.FileServer(http.Dir("webui/public")))

	// RestAPI
	ohrs := restapi.NewOrderHistoryRestService(ws.ctx, ws.exchangeService)
	http.HandleFunc("/orderhistory", ohrs.GetOrderHistory)

	// Websocket
	http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		portfolioHub := NewPortfolioHub(ws.ctx.Logger)
		go portfolioHub.Run()
		ph := NewPortfolioHandler(ws.ctx, portfolioHub, ws.marketcapService)
		ph.onConnect(w, r)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
	if err != nil {
		ws.ctx.Logger.Error("[Web] Unable to start server: ", err)
	}
}

func (ws *WebServer) Run() {
	ws.ctx.Logger.Debug("[WebServer.run] Starting loop")
	for {
		ws.ctx.Logger.Debug("[WebServer.run] Main loop...")
		select {
		case close := <-ws.closeChan:
			if close {
				ws.ctx.Logger.Debug("[WebServer.run] Stopping Web server")
				break
			}
		}
	}
}

func (ws *WebServer) Stop() {
	ws.closeChan <- true
}
