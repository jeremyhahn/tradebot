package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
)

type PortfolioServiceImpl struct {
	ctx              *common.Context
	stopChans        map[uint]chan bool
	portfolios       map[uint]common.Portfolio
	marketcapService *MarketCapService
	userService      UserService
	PortfolioService
}

func NewPortfolioService(ctx *common.Context, marketcapService *MarketCapService,
	userService UserService) PortfolioService {
	return &PortfolioServiceImpl{
		ctx:              ctx,
		stopChans:        make(map[uint]chan bool),
		portfolios:       make(map[uint]common.Portfolio),
		marketcapService: marketcapService,
		userService:      userService}
}

func (ps *PortfolioServiceImpl) Build(user common.User, currencyPair *common.CurrencyPair) common.Portfolio {
	ps.ctx.Logger.Debugf("[PortfolioService.Build] Building portfolio for %s", user.GetUsername())
	var netWorth float64
	exchangeList := ps.userService.GetExchanges(ps.ctx.GetUser(), currencyPair)
	walletList := ps.userService.GetWallets(ps.ctx.GetUser())
	for _, ex := range exchangeList {
		netWorth += ex.GetTotal()
	}
	for _, w := range walletList {
		netWorth += w.GetNetWorth()
	}
	currentUser, err := ps.userService.GetCurrentUser()
	if err != nil {
		ps.ctx.Logger.Errorf("[PortfolioService.Build] Error getting current user: %s", err.Error())
	}
	portfolio := &dto.PortfolioDTO{
		User:      currentUser,
		NetWorth:  netWorth,
		Exchanges: exchangeList,
		Wallets:   walletList}
	ps.portfolios[user.GetId()] = portfolio
	ps.stopChans[user.GetId()] = make(chan bool, 1)
	return portfolio
}

func (ps *PortfolioServiceImpl) Queue(user common.User) <-chan common.Portfolio {
	ps.ctx.Logger.Debugf("[PortfolioService.Queue] Adding portfolio to queue on behalf of %s", user.GetUsername())
	currencyPair := &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"}
	portfolio := ps.Build(user, currencyPair)
	ps.ctx.Logger.Debugf("[PortfolioService.Queue] portfolio=%+v\n", portfolio)
	portChan := make(chan common.Portfolio, 1)
	portChan <- portfolio
	return portChan
}

func (ps *PortfolioServiceImpl) Stream(user common.User, currencyPair *common.CurrencyPair) <-chan common.Portfolio {
	portfolio := ps.Build(user, currencyPair)
	ps.ctx.Logger.Debugf("[PortfolioService.Stream] Starting stream for %s", portfolio.GetUser().GetUsername())
	portChan := make(chan common.Portfolio, 10)
	go func() {
		for {
			select {
			case <-ps.stopChans[user.GetId()]:
				ps.ctx.Logger.Debug("[PortfolioService.Stream] Stopping stream")
				delete(ps.stopChans, user.GetId())
			default:
				ps.ctx.Logger.Debugf("[PortfolioService.Stream] Broadcasting portfolio: %+v\n", portfolio)
				portChan <- portfolio
			}
		}
	}()
	return portChan
}

func (ps *PortfolioServiceImpl) Stop(user common.User) {
	ps.ctx.Logger.Debugf("[PortfolioService.Stop] Stopping stream for %s\n", user.GetUsername())
	if ps.IsStreaming(user) {
		ps.stopChans[user.GetId()] <- true
	}
}

func (ps *PortfolioServiceImpl) IsStreaming(user common.User) bool {
	_, ok := ps.portfolios[user.GetId()]
	return ok
}
