package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
)

type PortfolioHandler struct {
	ctx *common.Context
	hub *PortfolioHub
}

func NewPortfolioHandler(ctx *common.Context, hub *PortfolioHub) *PortfolioHandler {
	return &PortfolioHandler{
		ctx: ctx,
		hub: hub}
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
		ph.ctx.Logger.Error(err)
	}
	if conn == nil {
		ph.ctx.Logger.Error("[PortfolioHandler.onConnect] Unable to establish websocket connection")
		return
	}

	var portfolio common.Portfolio
	err = conn.ReadJSON(&portfolio)
	if err != nil {
		ph.ctx.Logger.Errorf("[PortfolioHandler.onConnect] Websocket Read Error: %v", err)
		conn.Close()
		return
	}

	ph.ctx.Logger.Debug("[PortfolioHandler.onConnect] Accepting connection from ", conn.RemoteAddr())
	ph.ctx.User = portfolio.User
	client := &PortfolioClient{
		hub:  ph.hub,
		conn: conn,
		send: make(chan *common.Portfolio, common.BUFFERED_CHANNEL_SIZE),
		ctx:  ph.ctx}

	client.hub.register <- client
	go client.writePump()
	go client.readPump()
	go client.keepAlive()
}

/*
func (ph *PortfolioHandler) Broadcast(portfolio *common.Portfolio) {
	ph.ctx.Logger.Debugf("[PortfolioHandler.Broadcast] Portfolio: %+v\n", portfolio)
	for client := range ph.clients {
		err := client.WriteJSON(portfolio)
		if err != nil {
			ph.ctx.Logger.Error(err)
		}
	}
}
*/
