package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/service"
)

type PortfolioHandler struct {
	ctx              *common.Context
	hub              *PortfolioHub
	marketcapService *service.MarketCapService
	userService      service.UserService
	portfolioService service.PortfolioService
}

func NewPortfolioHandler(ctx *common.Context, hub *PortfolioHub,
	marketcapService *service.MarketCapService, userService service.UserService,
	portfolioService service.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{
		ctx:              ctx,
		hub:              hub,
		marketcapService: marketcapService,
		portfolioService: portfolioService}
}

func (ph *PortfolioHandler) OnConnect(w http.ResponseWriter, r *http.Request) {
	var user dto.UserDTO
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
		ph.ctx.Logger.Error("[PortfolioHandler.onConnect] Unable to establish webservice connection")
		return
	}
	err = conn.ReadJSON(&user)
	if err != nil {
		ph.ctx.Logger.Errorf("[PortfolioHandler.onConnect] webservice Read Error: %v", err)
		conn.Close()
		return
	}
	ph.ctx.Logger.Debug("[PortfolioHandler.onConnect] Accepting connection from ", conn.RemoteAddr())
	ph.stream(conn, &user)
}

func (ph *PortfolioHandler) stream(conn *websocket.Conn, user common.User) {
	ph.ctx.SetUser(user) // TODO: REPLACE THIS WITH "ETHEREUM SESSION USER"
	client := &PortfolioClient{
		hub:              ph.hub,
		conn:             conn,
		send:             make(chan common.Portfolio, common.BUFFERED_CHANNEL_SIZE),
		ctx:              ph.ctx,
		marketcapService: ph.marketcapService,
		userService:      ph.userService,
		portfolioService: ph.portfolioService}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
	//go client.keepAlive()
}
