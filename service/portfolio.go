package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
)

type PortfolioService struct {
	ctx              *common.Context
	stopChan         chan bool
	marketcapService *MarketCapService
	userService      UserService
}

func NewPortfolioService(ctx *common.Context, marketcapService *MarketCapService,
	userService UserService) *PortfolioService {
	return &PortfolioService{
		ctx:              ctx,
		stopChan:         make(chan bool),
		marketcapService: marketcapService,
		userService:      userService}
}

func (ps *PortfolioService) Build(user common.User, currencyPair *common.CurrencyPair) common.Portfolio {
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
	return &dto.PortfolioDTO{
		User:      ps.ctx.GetUser(),
		NetWorth:  netWorth,
		Exchanges: exchangeList,
		Wallets:   walletList}
}

func (ps *PortfolioService) Queue(user common.User) <-chan common.Portfolio {
	ps.ctx.Logger.Debugf("[PortfolioService.Queue] Adding portfolio to queue on behalf of %s", user.GetUsername())
	currencyPair := &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"}
	portfolio := ps.Build(user, currencyPair)
	ps.ctx.Logger.Debugf("[PortfolioService.Queue] portfolio=%+v\n", portfolio)
	portChan := make(chan common.Portfolio, 1)
	portChan <- portfolio
	return portChan
}

func (ps *PortfolioService) Stream(user common.User, currencyPair *common.CurrencyPair) <-chan common.Portfolio {
	portfolio := ps.Build(user, currencyPair)
	ps.ctx.Logger.Debugf("[PortfolioService.Stream] Starting stream for %s", portfolio.GetUser().GetUsername())
	portChan := make(chan common.Portfolio)
	go func() {
		for {
			select {
			case stop := <-ps.stopChan:
				ps.ctx.Logger.Debug("[PortfolioService.Stream] Stopping stream")
				if stop {
					return
				}
			default:
				ps.ctx.Logger.Debugf("[PortfolioService.Stream] Broadcasting portfolio: %+v\n", portfolio)
				portChan <- portfolio
			}
			time.Sleep(10 * time.Second)
		}
	}()
	return portChan
}

func (ps *PortfolioService) Stop() {
	ps.stopChan <- true
}
