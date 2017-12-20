package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type WebsocketServer struct {
	clients   map[*websocket.Conn]string
	broadcast chan *common.ChartData
	logger    *logging.Logger
	port      int
	charts    []*Chart
	common.PriceListener
}

type WebsocketRequest struct {
	Currency string `json:"currency"`
}

func NewWebsocketServer(port int, charts []*Chart, logger *logging.Logger) *WebsocketServer {
	ws := &WebsocketServer{
		clients:   make(map[*websocket.Conn]string),
		broadcast: make(chan *common.ChartData),
		logger:    logger,
		port:      port,
		charts:    charts}

	return ws
}

func (ws *WebsocketServer) Start() {

	ws.logger.Debug("[WebSocket] Starting file server on port: ", ws.port)
	http.Handle("/", http.FileServer(http.Dir("frontend/public")))

	http.HandleFunc("/chart", ws.onConnect)
	ws.logger.Debug("[WebSocket] Starting server on port: ", ws.port)

	go ws.run()

	err := http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
	if err != nil {
		ws.logger.Error("[WebSocket] Unable to start server: ", err)
	}
}

func (ws *WebsocketServer) onConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error(err)
	}
	if conn == nil {
		return
	}
	defer conn.Close()
	ws.logger.Debug("[WebSocket] Accepting connection from: ", conn.RemoteAddr())
	for {
		var msg WebsocketRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			ws.logger.Errorf("[WebSocket] Websocket Read Error: %v", err)
			delete(ws.clients, conn)
			break
		}
		ws.clients[conn] = msg.Currency
	}
}

func (ws *WebsocketServer) run() {
	for {
		select {
		case msg := <-ws.broadcast:
			logMsg := fmt.Sprintf("[Candlestick] Close: %.8f", msg.Price)
			ws.logger.Debug(logMsg)
			ws.logger.Debug("[WebSocket] Broadcasting: ", msg)
			for client := range ws.clients {
				if ws.clients[client] == msg.Currency {
					err := client.WriteJSON(msg)
					if err != nil {
						ws.logger.Error(err)
					}
				}
			}
		}
	}
}

func (ws *WebsocketServer) Broadcast(message *common.ChartData) {
	for _, chart := range ws.charts {
		if chart.GetChartData().Currency == message.Currency {
			ws.broadcast <- chart.GetChartData()
		}
	}
}
