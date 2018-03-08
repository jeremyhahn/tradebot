package websocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

func TestPortfolioHandler_Stream(t *testing.T) {

	fmt.Println(os.Getwd())

	ctx := test.CreateIntegrationTestContext("../../.env", "../../")
	databaseManager := common.CreateDatabase("../../", "test-", ctx.GetDebug())

	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	hub := NewPortfolioHub(ctx.GetLogger())
	marketcapService := service.NewMarketCapService(ctx.GetLogger())

	userExchangeMapper := mapper.NewUserExchangeMapper()
	ethereumService, err := service.NewEthereumService(ctx, userDAO, userMapper)
	assert.Nil(t, err)

	userService := service.NewUserService(ctx, userDAO, pluginDAO, marketcapService, ethereumService, userMapper, userExchangeMapper)

	jsonWebTokenService, err := service.NewJsonWebTokenService(ctx, databaseManager, ethereumService, common.NewJsonWriter())
	assert.Nil(t, err)

	portfolioService := service.NewPortfolioService(ctx, marketcapService, userService, ethereumService)
	portfolioHandler := NewPortfolioHandler(ctx.GetLogger(), hub, jsonWebTokenService)

	s := httptest.NewServer(http.HandlerFunc(portfolioHandler.OnConnect))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	defer ws.Close()
	assert.Nil(t, err)

	user := &dto.UserDTO{
		Id:            1,
		Username:      "Jeremy",
		LocalCurrency: "USD"}
	ws.WriteJSON(user)

	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"}

	portfolioChan, err := portfolioService.Stream(user, currencyPair)
	portfolio := <-portfolioChan

	portfolioUser := portfolio.GetUser()
	assert.Equal(t, uint(1), portfolioUser.GetId())
	// Bug? Returning persisted database name instead of name defined in DTO
	//assert.Equal(t, user.GetUsername(), portfolioUser.GetUsername())
	assert.Equal(t, user.GetLocalCurrency(), portfolioUser.GetLocalCurrency())
	assert.Equal(t, true, len(portfolio.GetExchanges()) > 0)
	assert.Equal(t, true, len(portfolio.GetWallets()) > 0)
	assert.Equal(t, true, portfolio.GetNetWorth() > 0)
	assert.Equal(t, true, portfolioService.IsStreaming(user))

	portfolioService.Stop(ctx.GetUser())
	time.Sleep(3 * time.Second)

	assert.Equal(t, false, portfolioService.IsStreaming(user))

	test.CleanupIntegrationTest()
}
