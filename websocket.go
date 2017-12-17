package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type WebsocketServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan *common.ChartData
	logger    *logging.Logger
	port      int
	chart     *Chart
	common.PriceListener
}

type WebsocketRequest struct {
	Message string  `json:"message"`
	Price   float64 `json:"price"`
}

func NewWebsocketServer(port int, chart *Chart, logger *logging.Logger) *WebsocketServer {
	ws := &WebsocketServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan *common.ChartData),
		logger:    logger,
		port:      port,
		chart:     chart}

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
	ws.clients[conn] = true
	for {
		var msg WebsocketRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			ws.logger.Errorf("[WebSocket] Websocket Read Error: %v", err)
			delete(ws.clients, conn)
			break
		}
	}
}

func (ws *WebsocketServer) run() {
	for {
		select {
		case msg := <-ws.broadcast:
			ws.logger.Debug("[WebSocket] Broadcasting: ", msg)
			data := ws.chart.GetChartData()
			for client := range ws.clients {
				err := client.WriteJSON(data)
				if err != nil {
					ws.logger.Error(err)
				}
			}
		}
	}
}

func (ws *WebsocketServer) Broadcast(price float64) {
	ws.logger.Debugf("[Websocket] OnPriceChange: %f", price)
	ws.broadcast <- ws.chart.GetChartData()
}
