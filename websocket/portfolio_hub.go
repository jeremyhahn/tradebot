package websocket

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type PortfolioHub struct {
	clients    map[*PortfolioClient]bool
	broadcast  chan *common.Portfolio
	register   chan *PortfolioClient
	unregister chan *PortfolioClient
}

func NewPortfolioHub() *PortfolioHub {
	return &PortfolioHub{
		broadcast:  make(chan *common.Portfolio),
		register:   make(chan *PortfolioClient),
		unregister: make(chan *PortfolioClient),
		clients:    make(map[*PortfolioClient]bool),
	}
}

func (h *PortfolioHub) Run() {
	for {

		fmt.Println("Portfolio hub running...")

		select {
		case client := <-h.register:
			client.ctx.Logger.Debugf("[PortfolioHub.run] Registering new client: %s", client.ctx.User.Username)
			fmt.Printf("[PortfolioHub.run] Registering new client: %s", client.ctx.User.Username)
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.ctx.Logger.Debugf("[PortfolioHub.run] Unregistering client: %s", client.ctx.User.Username)
				fmt.Printf("[PortfolioHub.run] Unregistering client: %s", client.ctx.User.Username)
				delete(h.clients, client)
				close(client.send)
			}

		case portfolio := <-h.broadcast:
			for client := range h.clients {
				client.ctx.Logger.Debugf("[PortfolioHub.run] Broadcasting portfolio: %+v\n", portfolio)
				fmt.Printf("[PortfolioHub.run] Broadcasting portfolio: %+v\n", portfolio)
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
