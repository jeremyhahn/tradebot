package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
	logging "github.com/op/go-logging"
)

type ChartHandler struct {
	clients map[*websocket.Conn]common.ChartData
	channel chan common.ChartData
	charts  []service.Chart
	logger  *logging.Logger
}

func NewChartHandler(logger *logging.Logger) *ChartHandler {
	return &ChartHandler{
		logger: logger}
}

func (ch *ChartHandler) onConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ch.logger.Error(err)
		return
	}
	if conn == nil {
		ch.logger.Error("[ChartHandler.onConnect] Unable to establish websocket connection")
		return
	}
	defer conn.Close()
	ch.logger.Debug("[ChartHandler.onConnect] Accepting connection from: ", conn.RemoteAddr())
	for {
		var msg WebsocketRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			ch.logger.Errorf("[ChartHandler.onConnect] Websocket Read Error: %v", err)
			delete(ch.clients, conn)
			break
		}
		ch.clients[conn] = common.ChartData{
			Exchange:     msg.Exchange,
			CurrencyPair: *msg.CurrencyPair}
	}
}

func (ch *ChartHandler) Broadcast(data *common.ChartData) {
	ch.logger.Debugf("[ChartHandler.Broadcast] ChartData: %+v\n", data)
	for client := range ch.clients {
		clientChart := ch.clients[client]
		if clientChart.Exchange == data.Exchange &&
			clientChart.CurrencyPair.Base == data.CurrencyPair.Base &&
			clientChart.CurrencyPair.Quote == data.CurrencyPair.Quote {
			err := client.WriteJSON(data)
			if err != nil {
				ch.logger.Error(err)
			}
		}
	}
}
