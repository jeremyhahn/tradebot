package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type WebsocketServer struct {
	chartClients     map[*websocket.Conn]string
	portfolioClients map[*websocket.Conn]string
	broadcast        chan *common.ChartData
	logger           *logging.Logger
	port             int
	charts           []*Chart
	running          bool
	common.PriceListener
}

type WebsocketRequest struct {
	Currency string `json:"currency"`
}

func NewWebsocketServer(port int, charts []*Chart, logger *logging.Logger) *WebsocketServer {
	return &WebsocketServer{
		chartClients:     make(map[*websocket.Conn]string),
		portfolioClients: make(map[*websocket.Conn]string),
		broadcast:        make(chan *common.ChartData),
		logger:           logger,
		port:             port,
		charts:           charts}
}

func (ws *WebsocketServer) Start() {

	ws.logger.Debug("[WebSocket] Starting file service on port: ", ws.port)
	http.Handle("/", http.FileServer(http.Dir("frontend/public")))

	ws.logger.Debug("[WebSocket] Starting websocket service on port: ", ws.port)

	http.HandleFunc("/chart", ws.onChartConnect)
	http.HandleFunc("/portfolio", ws.onPortfolioConnect)

	err := http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
	if err != nil {
		ws.logger.Error("[WebSocket] Unable to start server: ", err)
	}

	ws.running = true

	go ws.run()
}

func (ws *WebsocketServer) Stop() {
	ws.running = false
}

func (ws *WebsocketServer) onChartConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error(err)
	}
	if conn == nil {
		return
	}
	defer conn.Close()
	ws.logger.Debug("[WebSocket.onChartConnect] Accepting connection from: ", conn.RemoteAddr())
	for {
		var msg WebsocketRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			ws.logger.Errorf("[WebSocket.onChartConnect] Websocket Read Error: %v", err)
			delete(ws.chartClients, conn)
			break
		}
		ws.chartClients[conn] = msg.Currency
	}
}

func (ws *WebsocketServer) onPortfolioConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error(err)
	}
	if conn == nil {
		return
	}
	defer conn.Close()
	ws.logger.Debug("[WebSocket.OnPriceConnect] Accepting connection from: ", conn.RemoteAddr())
	for {
		var msg WebsocketRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			ws.logger.Errorf("[WebSocket.OnPriceConnect] Websocket Read Error: %v", err)
			delete(ws.portfolioClients, conn)
			break
		}
		ws.portfolioClients[conn] = msg.Currency
	}
}

func (ws *WebsocketServer) run() {
	for ws.running {
		select {
		case msg := <-ws.broadcast:
			logMsg := fmt.Sprintf("[Candlestick] Close: %.8f", msg.Price)
			ws.logger.Debug(logMsg)
			ws.logger.Debug("[WebSocket] Broadcasting: ", msg)
			for client := range ws.chartClients {
				if ws.chartClients[client] == msg.Currency {
					err := client.WriteJSON(msg)
					if err != nil {
						ws.logger.Error(err)
					}
				}
			}
		}
	}
	ws.logger.Debug("[WebSocket] Shutting down")
}

func (ws *WebsocketServer) BroadcastChart(message *common.ChartData) {
	for _, chart := range ws.charts {
		if chart.GetChartData().Currency == message.Currency {
			ws.broadcast <- chart.GetChartData()
		}
	}
}

func (ws *WebsocketServer) BroadcastPortfolio(exchangeList []common.CoinExchange) {
	fmt.Printf("[WebsocketServer.BroadcastPortfolio] ExchangeList: %+v\n", exchangeList)
	for client := range ws.portfolioClients {
		err := client.WriteJSON(exchangeList)
		if err != nil {
			ws.logger.Error(err)
		}
	}
}
