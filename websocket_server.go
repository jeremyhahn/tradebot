package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type WebsocketServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan *common.Chart
	logger    *logging.Logger
	traders   []*Trader
}

type WebsocketRequest struct {
	Message string
}

func NewWebsocketServer(port int, traders []*Trader, logger *logging.Logger) *WebsocketServer {
	ws := &WebsocketServer{
		clients: make(map[*websocket.Conn]bool),
		logger:  logger,
		traders: traders}
	fs := http.FileServer(http.Dir("public"))
	logger.Debug("Starting file server on port: ", port)
	http.Handle("/", fs)
	http.HandleFunc("/trader", ws.onConnect)
	logger.Debug("Starting WebSocket server on port: ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logger.Error("Unable to start WebSocket server: ", err)
	}
	go ws.run()
	return ws
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
	ws.logger.Debug("Accepting WebSocket connection from: ", conn.RemoteAddr())
	defer conn.Close()
	ws.clients[conn] = true
	for {
		var request WebsocketRequest
		err := conn.ReadJSON(request)
		if err != nil {
			ws.logger.Error(fmt.Sprintf("WebSocket Error: %#v", err))
			delete(ws.clients, conn)
			break
		}
		ws.broadcast <- ws.traders[0].GetChart()
	}
}

func (ws *WebsocketServer) run() {
	for {
		msg := <-ws.broadcast
		ws.logger.Debug("Received message: ", msg)
		for client := range ws.clients {
			jsonData, err := json.Marshal(&common.Chart{
				Price: 100})
			if err != nil {
				ws.logger.Error(err)
				return
			}
			//err = client.WriteJSON(jsonData)
			err = client.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(ws.clients, client)
			}
		}
	}
}
