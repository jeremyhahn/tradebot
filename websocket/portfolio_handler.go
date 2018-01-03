package websocket

import (
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
	logging "github.com/op/go-logging"
)

type PortfolioHandler struct {
	clients map[*websocket.Conn]net.Addr
	channel chan common.ChartData
	charts  []service.Chart
	logger  *logging.Logger
}

func NewPortfolioHandler(logger *logging.Logger) *PortfolioHandler {
	return &PortfolioHandler{
		clients: make(map[*websocket.Conn]net.Addr),
		logger:  logger}
}

func (ph *PortfolioHandler) onConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ph.logger.Error(err)
	}
	if conn == nil {
		ph.logger.Error("[PortfolioHandler.onConnect] Unable to establish websocket connection")
		return
	}
	defer conn.Close()
	ph.logger.Debug("[PortfolioHandler.onConnect] Accepting connection from: ", conn.RemoteAddr())
	for {
		var msg WebsocketRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			ph.logger.Errorf("[PortfolioHandler.onConnect] Websocket Read Error: %v", err)
			delete(ph.clients, conn)
			break
		}
		ph.clients[conn] = conn.RemoteAddr()
	}
}

func (ph *PortfolioHandler) Broadcast(portfolio *common.Portfolio) {
	ph.logger.Debugf("[PortfolioHandler.Broadcast] Portfolio: %+v\n", portfolio)
	for client := range ph.clients {
		err := client.WriteJSON(portfolio)
		if err != nil {
			ph.logger.Error(err)
		}
	}
}
