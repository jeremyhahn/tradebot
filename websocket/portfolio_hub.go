package websocket

import (
	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type PortfolioHub struct {
	logger     *logging.Logger
	clients    map[*PortfolioClient]bool
	broadcast  chan *common.Portfolio
	register   chan *PortfolioClient
	unregister chan *PortfolioClient
}

func NewPortfolioHub(logger *logging.Logger) *PortfolioHub {
	return &PortfolioHub{
		broadcast:  make(chan *common.Portfolio),
		register:   make(chan *PortfolioClient),
		unregister: make(chan *PortfolioClient),
		clients:    make(map[*PortfolioClient]bool),
		logger:     logger}
}

func (h *PortfolioHub) Run() {
	for {

		h.logger.Debug("Portfolio hub running...")

		select {
		case client := <-h.register:
			client.ctx.Logger.Debugf("[PortfolioHub.run] Registering new client: %s", client.ctx.User.Username)
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.logger.Debugf("[PortfolioHub.run] Unregistering client: %s", client.ctx.User.Username)
				client.disconnect()
				delete(h.clients, client)
				close(client.send)
			}

		case portfolio := <-h.broadcast:
			for client := range h.clients {
				h.logger.Debugf("[PortfolioHub.run] Broadcasting portfolio: %+v\n", portfolio)
				select {
				case client.send <- portfolio:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
