package webservice

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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
	ctx                 common.Context
	port                int
	closeChan           chan bool
	portfolioHandler    *websocket.PortfolioHandler
	authService         service.AuthService
	jsonWebTokenService service.JsonWebTokenService
}

func NewWebServer(ctx common.Context, port int, authService service.AuthService,
	jsonWebTokenService service.JsonWebTokenService) *WebServer {
	return &WebServer{
		ctx:                 ctx,
		port:                port,
		closeChan:           make(chan bool, 1),
		authService:         authService,
		jsonWebTokenService: jsonWebTokenService}
}

func (ws *WebServer) Start() {

	router := mux.NewRouter()
	jsonWriter := common.NewJsonWriter()
	fs := http.FileServer(http.Dir("webapp/public"))

	// REST Handlers - Public Access
	registrationService := rest.NewRegisterRestService(ws.ctx, ws.authService, jsonWriter)
	router.HandleFunc("/api/v1/register", registrationService.Register)
	router.HandleFunc("/api/v1/login", ws.jsonWebTokenService.GenerateToken)

	// REST Handlers - Authentication Required
	exchangeRestService := rest.NewExchangeRestService(ws.jsonWebTokenService, jsonWriter)
	userRestService := rest.NewUserRestService(ws.jsonWebTokenService, jsonWriter)
	transactionRestService := rest.NewTransactionRestService(ws.jsonWebTokenService, jsonWriter)
	router.Handle("/api/v1/transactions", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.GetHistory)),
	))
	router.Handle("/api/v1/transactions/orderhistory", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.GetOrderHistory)),
	))
	router.Handle("/api/v1/transactions/desposits", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.GetDepositHistory)),
	))
	router.Handle("/api/v1/transactions/withdrawals", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.GetWithdrawalHistory)),
	))
	router.Handle("/api/v1/transactions/imported", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.GetImportedTransactions)),
	))
	router.Handle("/api/v1/transactions/import", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.Import)),
	))
	router.Handle("/api/v1/transactions/sync", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.Synchronize)),
	))
	router.Handle("/api/v1/transactions/{id}", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(transactionRestService.UpdateCategory)),
	)).Methods("PUT")
	router.Handle("/api/v1/exchanges/names", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(exchangeRestService.GetDisplayNames)),
	))
	router.Handle("/api/v1/user/exchanges", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(userRestService.GetExchanges)),
	))
	router.Handle("/api/v1/user/exchange", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(userRestService.CreateExchange)),
	)).Methods("POST")
	router.Handle("/api/v1/user/exchange/{id}", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(userRestService.DeleteExchanges)),
	)).Methods("POST")

	// Websocket Handlers
	router.Handle("/ws/portfolio", negroni.New(
		negroni.HandlerFunc(ws.jsonWebTokenService.Validate),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			portfolioHub := websocket.NewPortfolioHub(ws.ctx.GetLogger())
			go portfolioHub.Run()
			ph := websocket.NewPortfolioHandler(ws.ctx.GetLogger(), portfolioHub, ws.jsonWebTokenService)
			ph.OnConnect(w, r)
		})),
	))

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
	if ws.ctx.GetSSL() {

		ws.ctx.GetLogger().Debugf("Starting web services on TLS port %d", ws.port)

		// Redirect HTTP -> HTTPS
		go http.ListenAndServe(sPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))

		// Serve HTTPS Requests
		err := http.ListenAndServeTLS(fmt.Sprintf(":%d", ws.port), "keys/cert.pem", "keys/key.pem", router)
		if err != nil {
			ws.ctx.GetLogger().Fatalf("[WebServer] Unable to start TLS web server: %s", err.Error())
		}
	} else {

		ws.ctx.GetLogger().Debugf("Starting web services on port %d", ws.port)

		// Serve HTTP Requests
		err := http.ListenAndServe(sPort, router)
		if err != nil {
			ws.ctx.GetLogger().Fatalf("[WebServer] Unable to start web server: %s", err.Error())
		}
	}
}

func (ws *WebServer) Run() {
	for {
		select {
		case <-ws.closeChan:
			ws.ctx.GetLogger().Debug("[WebServer.run] Stopping Web server")
			break
		}
	}
}

func (ws *WebServer) Stop() {
	ws.ctx.GetLogger().Info("Stopping web server")
	ws.closeChan <- true
}
