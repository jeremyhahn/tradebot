package webservice

import (
	"fmt"
	"net/http"

	"github.com-backup/gorilla/mux"
	"github.com/codegangsta/negroni"
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
	orderService     service.OrderService
	jwt              *JsonWebToken
}

func NewWebServer(ctx *common.Context, port int, marketcapService *service.MarketCapService,
	exchangeService service.ExchangeService, authService service.AuthService,
	userService service.UserService, portfolioService service.PortfolioService,
	orderService service.OrderService, jwt *JsonWebToken) *WebServer {
	return &WebServer{
		ctx:              ctx,
		port:             port,
		closeChan:        make(chan bool, 1),
		marketcapService: marketcapService,
		exchangeService:  exchangeService,
		authService:      authService,
		userService:      userService,
		portfolioService: portfolioService,
		orderService:     orderService,
		jwt:              jwt}
}

func (ws *WebServer) Start() {

	router := mux.NewRouter()
	jsonWriter := NewJsonWriter()
	fs := http.FileServer(http.Dir("webui/public"))

	// REST Handlers - Public Access
	reg := rest.NewRegisterRestService(ws.ctx, ws.authService, jsonWriter)
	router.HandleFunc("/api/v1/register", reg.Register)
	router.HandleFunc("/api/v1/login", ws.jwt.GenerateToken)

	// REST Handlers - Authentication Required
	ohrs := rest.NewOrderHistoryRestService(ws.ctx, ws.orderService, jsonWriter)
	ers := rest.NewExchangeRestService(ws.ctx, ws.exchangeService, ws.userService, jsonWriter)
	router.Handle("/api/v1/orderhistory", negroni.New(
		negroni.HandlerFunc(ws.jwt.MiddlewareValidator),
		negroni.Wrap(http.HandlerFunc(ohrs.GetOrderHistory)),
	))
	router.Handle("/api/v1/import", negroni.New(
		negroni.HandlerFunc(ws.jwt.MiddlewareValidator),
		negroni.Wrap(http.HandlerFunc(ohrs.Import)),
	))
	router.Handle("/api/v1/exchanges", negroni.New(
		negroni.HandlerFunc(ws.jwt.MiddlewareValidator),
		negroni.Wrap(http.HandlerFunc(ers.GetExchanges)),
	))

	// Websocket Handlers
	router.HandleFunc("/ws/portfolio", func(w http.ResponseWriter, r *http.Request) {
		portfolioHub := websocket.NewPortfolioHub(ws.ctx.Logger)
		go portfolioHub.Run()
		ph := websocket.NewPortfolioHandler(ws.ctx, portfolioHub, ws.marketcapService, ws.userService, ws.portfolioService)
		ph.OnConnect(w, r)
	})

	// React Routes
	routes := []string{"login", "register", "portfolio", "trades", "orders",
		"exchanges", "settings", "logout", "scripts"}
	for _, r := range routes {
		route := fmt.Sprintf("/%s", r)
		router.Handle(route, http.StripPrefix(route, fs))
	}

	// Default route / static content
	router.PathPrefix("/").Handler(fs)
	http.Handle("/", router)

	sPort := fmt.Sprintf(":%d", ws.port)
	if ws.ctx.SSL {

		ws.ctx.Logger.Debugf("Starting web services on TLS port %d", ws.port)

		// Redirect HTTP -> HTTPS
		go http.ListenAndServe(sPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))

		// Serve HTTPS Requests
		err := http.ListenAndServeTLS(fmt.Sprintf(":%d", ws.port), "keys/cert.pem", "keys/key.pem", router)
		if err != nil {
			ws.ctx.Logger.Fatalf("[WebServer] Unable to start TLS web server: %s", err.Error())
		}
	} else {

		ws.ctx.Logger.Debugf("Starting web services on port %d", ws.port)

		// Serve HTTP Requests
		err := http.ListenAndServe(sPort, router)
		if err != nil {
			ws.ctx.Logger.Fatalf("[WebServer] Unable to start web server: %s", err.Error())
		}
	}
}

func (ws *WebServer) Run() {
	for {
		select {
		case <-ws.closeChan:
			ws.ctx.Logger.Debug("[WebServer.run] Stopping Web server")
			break
		}
	}
}

func (ws *WebServer) Stop() {
	ws.ctx.Logger.Info("Stopping web server")
	ws.closeChan <- true
}
