package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	logging "github.com/op/go-logging"
)

type PortfolioHandler struct {
	logger            *logging.Logger
	hub               *PortfolioHub
	middlewareService service.Middleware
}

func NewPortfolioHandler(logger *logging.Logger, hub *PortfolioHub, middlewareService service.Middleware) *PortfolioHandler {
	return &PortfolioHandler{
		logger:            logger,
		hub:               hub,
		middlewareService: middlewareService}
}

func (ph *PortfolioHandler) OnConnect(w http.ResponseWriter, r *http.Request) {
	var user dto.UserContextDTO
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
		ph.logger.Error("[PortfolioHandler.onConnect] Unable to establish webservice connection")
		return
	}
	err = conn.ReadJSON(&user)
	if err != nil {
		ph.logger.Errorf("[PortfolioHandler.onConnect] webservice Read Error: %v", err)
		conn.Close()
		return
	}
	ph.logger.Debug("[PortfolioHandler.onConnect] Accepting connection from ", conn.RemoteAddr())

	ctx := ph.middlewareService.GetContext(user.GetId())
	if ctx == nil {
		ctx.GetLogger().Errorf("[PortfolioHandler.stream] Error: Unable to retrieve context from JsonWebTokenService")
		return
	}

	userDAO := dao.NewUserDAO(ctx)
	pluginDAO := dao.NewPluginDAO(ctx)

	userMapper := mapper.NewUserMapper()
	userExchangeMapper := mapper.NewUserExchangeMapper()

	marketcapService := service.NewMarketCapService(ctx)
	priceHistoryService := service.NewPriceHistoryService(ctx)
	ethereumService, _ := service.NewEthereumService(ctx, userDAO, userMapper, marketcapService)
	if err != nil {
		ctx.GetLogger().Errorf("[PortfolioHandler.stream] Error: %s", err.Error())
		return
	}

	userService := service.NewUserService(ctx, userDAO, pluginDAO, marketcapService, ethereumService, userMapper, userExchangeMapper, priceHistoryService)
	portfolioService := service.NewPortfolioService(ctx, marketcapService, userService, ethereumService)

	client := &PortfolioClient{
		hub:              ph.hub,
		conn:             conn,
		send:             make(chan common.Portfolio, common.BUFFERED_CHANNEL_SIZE),
		ctx:              ctx,
		marketcapService: marketcapService,
		userService:      userService,
		portfolioService: portfolioService}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
	//go client.keepAlive()
}
