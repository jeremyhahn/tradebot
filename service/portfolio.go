package service

import (
	"math/big"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
)

type PortfolioServiceImpl struct {
	ctx              common.Context
	stopChans        map[uint]chan bool
	portfolios       map[uint]common.Portfolio
	marketcapService *MarketCapService
	userService      UserService
	ethereumService  EthereumService
	PortfolioService
}

func NewPortfolioService(ctx common.Context, marketcapService *MarketCapService,
	userService UserService, ethereumService EthereumService) PortfolioService {
	return &PortfolioServiceImpl{
		ctx:              ctx,
		stopChans:        make(map[uint]chan bool),
		portfolios:       make(map[uint]common.Portfolio),
		marketcapService: marketcapService,
		userService:      userService,
		ethereumService:  ethereumService}
}

func (ps *PortfolioServiceImpl) Build(user common.UserContext, currencyPair *common.CurrencyPair) (common.Portfolio, error) {
	ps.ctx.GetLogger().Debugf("[PortfolioService.Build] Building portfolio for %s", user.GetUsername())
	var netWorth float64
	exchangeList, err := ps.userService.GetExchangeSummary(currencyPair)
	if err != nil {
		ps.ctx.GetLogger().Errorf("[PortfolioService.Build] Error: %s", err.Error())
		return nil, err
	}
	walletList := ps.userService.GetWallets(ps.ctx.GetUser())
	var tokenList []common.EthereumToken
	for _, ex := range exchangeList {
		netWorth += ex.GetTotal()
	}
	for _, w := range walletList {
		netWorth += w.GetValue()
	}
	accounts, err := ps.ethereumService.GetAccounts()
	if err != nil {
		ps.ctx.GetLogger().Errorf("[PortfolioService.Build] Error getting local Ethereum accounts: %s", err.Error())
	}
	for _, acct := range accounts {
		sAcct := acct.Address.String()
		balance, err := ps.ethereumService.GetBalance(sAcct)
		if err != nil {
			ps.ctx.GetLogger().Errorf("[PortfolioService.Build] Error getting Ethereum account balance for address %s: %s", sAcct, err.Error())
		}
		floatBalance, _ := new(big.Float).SetInt(balance).Float64()
		priceUSD, err := strconv.ParseFloat(ps.marketcapService.GetMarket("ETH").PriceUSD, 64)
		if err != nil {
			ps.ctx.GetLogger().Errorf("[PortfolioService.Build] Error parsing MarketCap ETH response to float for address %s: %s", sAcct, err.Error())
		}
		total := floatBalance * priceUSD
		walletList = append(walletList, &dto.UserCryptoWalletDTO{
			Address:  sAcct,
			Balance:  floatBalance,
			Currency: "ETH",
			Value:    total})

		netWorth += total

		tokens, err := ps.userService.GetTokens(ps.ctx.GetUser(), sAcct)
		if err != nil {
			ps.ctx.GetLogger().Errorf("[PortfolioService.Build] Error getting current user: %s", err.Error())
		}
		for _, token := range tokens {
			f, err := strconv.ParseFloat(token.GetEthBalance(), 64)
			if err != nil {
				ps.ctx.GetLogger().Errorf("[PortfolioService.Build] Error parsing token balance to float: %s", err.Error())
			}
			tokenList = append(tokenList, token)
			netWorth += f * priceUSD
		}
	}

	currentUser, err := ps.userService.GetCurrentUser()
	if err != nil {
		ps.ctx.GetLogger().Errorf("[PortfolioService.Build] Error getting current user: %s", err.Error())
	}
	portfolio := &dto.PortfolioDTO{
		User:      currentUser,
		NetWorth:  netWorth,
		Exchanges: exchangeList,
		Wallets:   walletList}
	ps.portfolios[user.GetId()] = portfolio
	ps.stopChans[user.GetId()] = make(chan bool, 1)
	return portfolio, nil
}

func (ps *PortfolioServiceImpl) Queue(user common.UserContext) (<-chan common.Portfolio, error) {
	ps.ctx.GetLogger().Debugf("[PortfolioService.Queue] Adding portfolio to queue on behalf of %s", user.GetUsername())
	currencyPair := &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"}
	portfolio, err := ps.Build(user, currencyPair)
	if err != nil {
		ps.ctx.GetLogger().Debugf("[PortfolioService.Queue] Error: %s", err.Error())
		return nil, err
	}
	ps.ctx.GetLogger().Debugf("[PortfolioService.Queue] portfolio=%+v\n", portfolio)
	portChan := make(chan common.Portfolio, 1)
	portChan <- portfolio
	return portChan, nil
}

func (ps *PortfolioServiceImpl) Stream(user common.UserContext, currencyPair *common.CurrencyPair) (<-chan common.Portfolio, error) {
	portfolio, err := ps.Build(user, currencyPair)
	if err != nil {
		ps.ctx.GetLogger().Errorf("[PortfolioService.Stream] Error: %s", err.Error())
		return nil, err
	}
	ps.ctx.GetLogger().Debugf("[PortfolioService.Stream] Starting stream for %s", portfolio.GetUser().GetUsername())
	portChan := make(chan common.Portfolio, 10)
	go func() {
		for {
			select {
			case <-ps.stopChans[user.GetId()]:
				ps.ctx.GetLogger().Debug("[PortfolioService.Stream] Stopping stream")
				delete(ps.stopChans, user.GetId())
			default:
				ps.ctx.GetLogger().Debugf("[PortfolioService.Stream] Broadcasting portfolio: %+v\n", portfolio)
				portChan <- portfolio
			}
		}
	}()
	return portChan, nil
}

func (ps *PortfolioServiceImpl) Stop(user common.UserContext) {
	ps.ctx.GetLogger().Debugf("[PortfolioService.Stop] Stopping stream for %s\n", user.GetUsername())
	if ps.IsStreaming(user) {
		ps.stopChans[user.GetId()] <- true
	}
}

func (ps *PortfolioServiceImpl) IsStreaming(user common.UserContext) bool {
	_, ok := ps.portfolios[user.GetId()]
	return ok
}
