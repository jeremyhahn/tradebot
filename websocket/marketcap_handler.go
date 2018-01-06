package websocket

import (
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type MarketCapHandler struct {
	clients map[*websocket.Conn]net.Addr
	channel chan common.MarketCap
	logger  *logging.Logger
}

func NewMarketCapHandler(logger *logging.Logger) *MarketCapHandler {
	return &MarketCapHandler{
		clients: make(map[*websocket.Conn]net.Addr),
		logger:  logger}
}

func (mc *MarketCapHandler) onConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		mc.logger.Error(err)
	}
	if conn == nil {
		mc.logger.Error("[MarketCapHandler.onConnect] Unable to establish websocket connection")
		return
	}
	defer conn.Close()
	mc.logger.Debug("[MarketCapHandler.onConnect] Accepting connection from: ", conn.RemoteAddr())
	for {
		var msg WebsocketRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			mc.logger.Errorf("[MarketCapHandler.onConnect] Websocket Read Error: %v", err)
			delete(mc.clients, conn)
			break
		}
		mc.clients[conn] = conn.RemoteAddr()
	}
}

func (mc *MarketCapHandler) Broadcast(marketcap *common.MarketCap) {
	mc.logger.Debugf("[MarketCapHandler.Broadcast] MarketCap: %+v\n", marketcap)
	for client := range mc.clients {
		err := client.WriteJSON(marketcap)
		if err != nil {
			mc.logger.Error(err)
		}
	}
}
