package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

type PortfolioClient struct {
	ctx              common.Context
	service          service.PortfolioService
	hub              *PortfolioHub
	conn             *websocket.Conn
	send             chan common.Portfolio
	marketcapService common.MarketCapService
	userService      service.UserService
	portfolioService service.PortfolioService
}

func (c *PortfolioClient) disconnect() {
	c.portfolioService.Stop(c.ctx.GetUser())
	c.conn.Close()
}

func (c *PortfolioClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var portfolio common.Portfolio
		err := c.conn.ReadJSON(&portfolio)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- portfolio
	}
}

func (c *PortfolioClient) writePump() {
	defer func() {
		c.portfolioService.Stop(c.ctx.GetUser())
		c.conn.Close()
	}()
	for {
		queue, err := c.portfolioService.Queue(c.ctx.GetUser())
		if err != nil {
			c.ctx.GetLogger().Errorf("[PortfolioClient.writePump] Error: %s", err.Error())
			continue
		}
		select {
		case message, _ := <-c.send:
			err := c.conn.WriteJSON(message)
			if err != nil {
				c.ctx.GetLogger().Errorf("[PortfolioClient.writePump] Error: %s", err.Error())
				return
			}
			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.WriteJSON(<-c.send)
			}
		case portfolio := <-queue:
			if err := c.conn.WriteJSON(portfolio); err != nil {
				c.ctx.GetLogger().Errorf("[PortfolioClient.writePump] Error: %s", err.Error())
				return
			}
			time.Sleep(10 * time.Second)
		}
	}
}

/*
func (c *PortfolioClient) keepAlive() {
	lastResponse := time.Now()
	c.conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})
	go func() {
		for {
			c.ctx.GetLogger().Debug("[PortfolioClient.keepAlive] Sending keepalive")
			err := c.conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				c.ctx.GetLogger().Debugf("[PortfolioClient.keepAlive] Error: %s", err.Error())
				return
			}
			time.Sleep(common.WEBSOCKET_KEEPALIVE / 2)
			if time.Now().Sub(lastResponse) > common.WEBSOCKET_KEEPALIVE {
				c.ctx.GetLogger().Debug("[PortfolioClient.keepAlive] Closing orphaned connection")
				c.conn.Close()
				return
			}
		}
	}()
}
*/
