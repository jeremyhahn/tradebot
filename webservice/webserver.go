package webservice

import (
	"fmt"
	"net/http"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/webservice/rest"
	"github.com/jeremyhahn/tradebot/webservice/websocket"
)

type WebRequest struct {
	Exchange     string `json:"exchange"`
	CurrencyPair *common.CurrencyPair
}

type WebServer struct {
	ctx              *common.Context
	port             int
	closeChan        chan bool
	portfolioHandler *websocket.PortfolioHandler
	marketcapService *service.MarketCapService
	exchangeService  service.ExchangeService
	authService      service.AuthService
	userService      service.UserService
	portfolioService service.PortfolioService
}

func NewWebServer(ctx *common.Context, port int, marketcapService *service.MarketCapService,
	exchangeService service.ExchangeService, authService service.AuthService,
	userService service.UserService, portfolioService service.PortfolioService) *WebServer {
	return &WebServer{
		ctx:              ctx,
		port:             port,
		marketcapService: marketcapService,
		exchangeService:  exchangeService,
		authService:      authService,
		userService:      userService,
		portfolioService: portfolioService}
}

func (ws *WebServer) Start() {

	ws.ctx.Logger.Debug("[WebServer] Starting on port: ", ws.port)

	// Static content
	http.Handle("/", http.FileServer(http.Dir("webui/public")))

	// REST Handlers
	ohrs := rest.NewOrderHistoryRestService(ws.ctx, ws.exchangeService)
	as := rest.NewLoginRestService(ws.ctx, ws.authService)
	reg := rest.NewRegisterRestService(ws.ctx, ws.authService)
	http.HandleFunc("/orderhistory", ohrs.GetOrderHistory)
	http.HandleFunc("/login", as.Login)
	http.HandleFunc("/register", reg.Register)

	// Websocket Handlers
	http.HandleFunc("/portfolio", func(w http.ResponseWriter, r *http.Request) {
		portfolioHub := websocket.NewPortfolioHub(ws.ctx.Logger)
		go portfolioHub.Run()
		ph := websocket.NewPortfolioHandler(ws.ctx, portfolioHub, ws.marketcapService, ws.userService, ws.portfolioService)
		ph.OnConnect(w, r)
	})

	sPort := fmt.Sprintf(":%d", ws.port)
	if ws.ctx.SSL {

		// Redirect HTTP -> HTTPS
		go http.ListenAndServe(sPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))

		// Serve HTTPS Requests
		err := http.ListenAndServeTLS(fmt.Sprintf(":%d", ws.port), "ssl/cert.pem", "ssl/key.pem", nil)
		if err != nil {
			ws.ctx.Logger.Fatalf("[WebServer] Unable to start TLS web server: %s", err.Error())
		}
	} else {

		// Serve HTTP Requests
		err := http.ListenAndServe(sPort, nil)
		if err != nil {
			ws.ctx.Logger.Fatalf("[WebServer] Unable to start web server: %s", err.Error())
		}
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
