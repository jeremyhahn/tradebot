package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type PortfolioClient struct {
	hub       *PortfolioHub
	conn      *websocket.Conn
	send      chan *common.Portfolio
	portfolio *common.Portfolio
	ctx       *common.Context
}

func (c *PortfolioClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var portfolio common.Portfolio
		err := c.conn.ReadJSON(&portfolio)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- &portfolio
	}
}

func (c *PortfolioClient) writePump() {
	portfolio := service.NewPortfolioService(c.ctx)
	defer func() {
		portfolio.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				c.ctx.Logger.Errorf("[PortfolioClient.writePump] Error: %s", err.Error())
				return
			}

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.WriteJSON(<-c.send)
			}

			if err := c.conn.Close(); err != nil {
				return
			}

		case portfolio := <-portfolio.Queue(c.ctx.User):
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteJSON(portfolio); err != nil {
				c.ctx.Logger.Errorf("[PortfolioClient.writePump] Error: %s", err.Error())
				return
			}

			time.Sleep(5 * time.Second)
		}
	}
}
