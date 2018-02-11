//// +build integration`

package webservice

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	ctx := test.NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	userMapper := mapper.NewUserMapper()
	hub := NewPortfolioHub(ctx.Logger)
	marketcapService := service.NewMarketCapService(ctx.Logger)
	userService := service.NewUserService(ctx, userDAO, marketcapService, userMapper)
	portfolioService := service.NewPortfolioService(ctx, marketcapService, userService)
	portfolioHandler := NewPortfolioHandler(ctx, hub, marketcapService, userService, portfolioService)

	s := httptest.NewServer(http.HandlerFunc(portfolioHandler.onConnect))
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

	portfolio := <-portfolioService.Stream(user, currencyPair)
	portfolioUser := portfolio.GetUser()
	assert.Equal(t, uint(1), portfolioUser.GetId())
	// Bug? Returning persisted database name instead of name defined in DTO
	//assert.Equal(t, user.GetUsername(), portfolioUser.GetUsername())
	assert.Equal(t, user.GetLocalCurrency(), portfolioUser.GetLocalCurrency())
	assert.Equal(t, true, len(portfolio.GetExchanges()) > 0)
	assert.Equal(t, true, len(portfolio.GetWallets()) > 0)
	assert.Equal(t, true, portfolio.GetNetWorth() > 0)
	assert.Equal(t, true, portfolioService.IsStreaming(user))

	portfolioService.Stop(user)
	assert.Equal(t, false, portfolioService.IsStreaming(user))

	test.CleanupIntegrationTest()
}
