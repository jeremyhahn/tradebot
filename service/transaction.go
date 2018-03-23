package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
)

type TransactionServiceImpl struct {
	ctx              common.Context
	dao              dao.TransactionDAO
	mapper           mapper.TransactionMapper
	exchangeService  ExchangeService
	ethereumService  EthereumService
	fiatPriceService common.FiatPriceService
	TransactionService
}

func NewTransactionService(ctx common.Context, transactionDAO dao.TransactionDAO, transactionMapper mapper.TransactionMapper,
	exchangeService ExchangeService, ethereumService EthereumService, fiatPriceService common.FiatPriceService) TransactionService {
	return &TransactionServiceImpl{
		ctx:              ctx,
		dao:              transactionDAO,
		mapper:           transactionMapper,
		exchangeService:  exchangeService,
		ethereumService:  ethereumService,
		fiatPriceService: fiatPriceService}
}

func (service *TransactionServiceImpl) GetMapper() mapper.TransactionMapper {
	return service.mapper
}

func (service *TransactionServiceImpl) GetTransactions() ([]common.Transaction, error) {
	transactions, err := service.ethereumService.GetTransactions()
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, service.GetOrderHistory()...)
	transactions = append(transactions, service.GetImportedTransactions()...)
	transactions = append(transactions, service.GetDepositHistory()...)
	transactions = append(transactions, service.GetWithdrawalHistory()...)
	service.sort(&transactions)
	return transactions, nil
}

func (service *TransactionServiceImpl) GetOrderHistory() []common.Transaction {
	service.ctx.GetLogger().Debugf("[TransactionService.GetOrderHistory] Retrieving order history for: %s",
		service.ctx.GetUser().GetUsername())
	var txs []common.Transaction
	exchanges, err := service.exchangeService.GetExchanges()
	if err != nil {
		service.ctx.GetLogger().Errorf("[TransactionService.GetOrderHistory] Error loading user exchanges: %s", err.Error())
	}
	for _, ex := range exchanges {
		if ex.GetName() == "gdax" || ex.GetName() == "coinbase" { // TODO standardize exchange interface
			balances, _ := ex.GetBalances()
			for _, coin := range balances {
				currencyPair := &common.CurrencyPair{
					Base:          coin.GetCurrency(),
					Quote:         service.ctx.GetUser().GetLocalCurrency(),
					LocalCurrency: service.ctx.GetUser().GetLocalCurrency()}
				txs = append(txs, ex.GetOrderHistory(currencyPair)...)
			}
			continue
		}
		currencyPairs, err := service.exchangeService.GetCurrencyPairs(ex.GetName())
		if err != nil {
			service.ctx.GetLogger().Errorf("[TransactionService.GetTransactionHistory] Error: %s", err.Error())
			return txs
		}
		for _, currencyPair := range currencyPairs {
			history := ex.GetOrderHistory(&common.CurrencyPair{
				Base:          currencyPair.Base,
				Quote:         currencyPair.Quote,
				LocalCurrency: service.ctx.GetUser().GetLocalCurrency()})
			txs = append(txs, history...)
		}
	}
	service.sort(&txs)
	return txs
}

func (service *TransactionServiceImpl) GetDepositHistory() []common.Transaction {
	service.ctx.GetLogger().Debugf("[TransactionService.GetDepositHistory] Retrieving deposit history for: %s",
		service.ctx.GetUser().GetUsername())
	var txs []common.Transaction
	exchanges, err := service.exchangeService.GetExchanges()
	if err != nil {
		service.ctx.GetLogger().Errorf("[TransactionService.GetDepositHistory] Error loading user exchanges: %s", err.Error())
	}
	for _, ex := range exchanges {
		deposits, err := ex.GetDepositHistory()
		if err != nil {
			service.ctx.GetLogger().Errorf("[TransactionService.GetDepositHistory] Error loading %s deposits: %s", ex.GetName(), err.Error())
			continue
		}
		txs = append(txs, deposits...)
	}
	service.sort(&txs)
	return txs
}

func (service *TransactionServiceImpl) GetWithdrawalHistory() []common.Transaction {
	service.ctx.GetLogger().Debugf("[TransactionService.GetWithdrawalHistory] Retrieving withdrawal history for: %s",
		service.ctx.GetUser().GetUsername())
	var txs []common.Transaction
	exchanges, err := service.exchangeService.GetExchanges()
	if err != nil {
		service.ctx.GetLogger().Errorf("[TransactionService.GetWithdrawalHistory] Error loading user exchanges: %s", err.Error())
	}
	for _, ex := range exchanges {
		withdrawals, err := ex.GetWithdrawalHistory()
		if err != nil {
			service.ctx.GetLogger().Errorf("[TransactionService.GetWithdrawalHistory] Error loading %s withdrawals: %s", ex.GetName(), err.Error())
			continue
		}
		txs = append(txs, withdrawals...)
	}
	service.sort(&txs)
	return txs
}

func (service *TransactionServiceImpl) GetImportedTransactions() []common.Transaction {
	var txs []common.Transaction
	txEntities, err := service.dao.Find()
	if err != nil {
		service.ctx.GetLogger().Errorf("[TransactionService.GetTransactionHistory] %s", err.Error())
	} else {
		for _, entity := range txEntities {
			txs = append(txs, service.mapper.MapTransactionEntityToDto(&entity))
		}
	}
	service.sort(&txs)
	return txs
}

func (service *TransactionServiceImpl) ImportCSV(file, exchangeName string) ([]common.Transaction, error) {
	service.ctx.GetLogger().Debugf("[TransactionService.ImportCSV] Creating %s exchange service", exchangeName)
	exchange, err := service.exchangeService.GetExchange(exchangeName)
	if err != nil {
		return nil, err
	}
	txDTOs, err := exchange.ParseImport(file)
	if err != nil {
		return nil, err
	}
	for _, dto := range txDTOs {
		entity := service.mapper.MapTransactionDtoToEntity(dto)
		service.dao.Create(entity)
	}
	return txDTOs, nil
}

func (service *TransactionServiceImpl) sort(txs *[]common.Transaction) {
	service.ctx.GetLogger().Debug("[TransactionService.sort] Sorting transactions")
	sort.Slice(*txs, func(i, j int) bool {
		return (*txs)[i].GetDate().After((*txs)[j].GetDate())
	})
}

/*
func (service *TransactionServiceImpl) GetSourceTransaction(targetTx common.Transaction,
	transactions *[]common.Transaction) (common.Transaction, error) {
	var candidates []common.Transaction
	for _, sourceTx := range *transactions {
		if targetTx.GetNetwork() == sourceTx.GetNetwork() || sourceTx.GetType() != common.BUY_ORDER_TYPE {
			continue
		}
		if sourceTx.GetCurrencyPair().Base != targetTx.GetCurrencyPair().Quote {
			continue
		}
		if sourceTx.GetDate().After(targetTx.GetDate()) {
			continue
		}
		if sourceTx.GetNetwork() == "GDAX" {
			service.ctx.GetLogger().Debugf(
				"[TransactionService].GetSourceTransaction] Comparing requested transaction %+v against historical transaction %s",
				targetTx, sourceTx)
		}
		if sourceTx.GetDate().Before(targetTx.GetDate()) || sourceTx.GetDate().Equal(targetTx.GetDate()) {
			if sourceTx.GetQuantity() == targetTx.GetTotal() && sourceTx.GetCurrencyPair().Base == targetTx.GetCurrencyPair().Quote {
				candidates = append(candidates, targetTx)
			}
		}
	}
	if len(candidates) == 0 {
		service.ctx.GetLogger().Errorf("[TransactionService].GetSourceTransaction] No source transaction found for: %+v", targetTx)
		return nil, errors.New("Source transaction not found")
	}
	if len(candidates) > 1 { // TODO: Create configurable conflict resolution strategy
		service.ctx.GetLogger().Warningf("[TransactionService].GetSourceTransaction] Found multiple source price candidates for transaction: %+v", targetTx)
	}
	service.ctx.GetLogger().Debugf("[TransactionService].GetSourceTransaction] Returning source transaction %+v", candidates[0])
	return candidates[0], nil
}*/
