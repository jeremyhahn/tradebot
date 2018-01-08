package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type PortfolioService struct {
	ctx              *common.Context
	stopChan         chan bool
	marketcapService *MarketCapService
}

func NewPortfolioService(ctx *common.Context, marketcapService *MarketCapService) *PortfolioService {
	return &PortfolioService{
		ctx:              ctx,
		stopChan:         make(chan bool),
		marketcapService: marketcapService}
}

func (ps *PortfolioService) Build(user *common.User) *common.Portfolio {
	ps.ctx.Logger.Debugf("[PortfolioService.Build] Building portfolio for %s", user.Username)
	var netWorth float64
	userDAO := dao.CreateUserDAO(ps.ctx, user)
	userService := NewUserService(ps.ctx, userDAO, ps.marketcapService)
	exchangeList := userService.GetExchanges(ps.ctx.User)
	walletList := userService.GetWallets(ps.ctx.User)
	for _, ex := range exchangeList {
		netWorth += ex.Total
	}
	for _, w := range walletList {
		netWorth += w.NetWorth
	}
	return &common.Portfolio{
		User:      ps.ctx.User,
		NetWorth:  netWorth,
		Exchanges: exchangeList,
		Wallets:   walletList}
}

func (ps *PortfolioService) Queue(user *common.User) <-chan *common.Portfolio {
	ps.ctx.Logger.Debugf("[PortfolioService.Queue] Adding portfolio to queue on behalf of %s", user.Username)
	portfolio := ps.Build(user)
	ps.ctx.Logger.Debugf("[PortfolioService.Queue] portfolio=%+v\n", portfolio)
	portChan := make(chan *common.Portfolio, 1)
	portChan <- portfolio
	return portChan
}

func (ps *PortfolioService) Stream(user *common.User) <-chan *common.Portfolio {
	portfolio := ps.Build(user)
	ps.ctx.Logger.Debugf("[PortfolioService.Stream] Starting stream for %s", portfolio.User.Username)
	portChan := make(chan *common.Portfolio)
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
