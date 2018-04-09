package service

import (
	"sort"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/mapper"
)

type TransactionServiceImpl struct {
	ctx              common.Context
	dao              dao.TransactionDAO
	mapper           mapper.TransactionMapper
	exchangeService  ExchangeService
	userService      UserService
	ethereumService  EthereumService
	fiatPriceService common.FiatPriceService
	TransactionService
}

func NewTransactionService(ctx common.Context, transactionDAO dao.TransactionDAO, transactionMapper mapper.TransactionMapper,
	exchangeService ExchangeService, userService UserService, ethereumService EthereumService,
	fiatPriceService common.FiatPriceService) TransactionService {
	return &TransactionServiceImpl{
		ctx:              ctx,
		dao:              transactionDAO,
		mapper:           transactionMapper,
		exchangeService:  exchangeService,
		userService:      userService,
		ethereumService:  ethereumService,
		fiatPriceService: fiatPriceService}
}

func (service *TransactionServiceImpl) GetMapper() mapper.TransactionMapper {
	return service.mapper
}

func (service *TransactionServiceImpl) isUnique(needle common.Transaction, haystack *[]common.Transaction,
	persisted *[]entity.Transaction) (bool, error) {
	for _, persistedTx := range *persisted {
		for _, tx := range *haystack {

			service.ctx.GetLogger().Debugf("Comparing network %s to %s", tx.GetNetwork(), persistedTx.GetNetwork())
			service.ctx.GetLogger().Debugf("Comparing date %s to %s", tx.GetDate(), persistedTx.GetDate())
			service.ctx.GetLogger().Debugf("Comparing quantity %s to %s", tx.GetQuantity(), persistedTx.GetQuantity())
			service.ctx.GetLogger().Debugf("Comparing total %s to %s", tx.GetTotal(), persistedTx.GetTotal())

			if tx.GetNetwork() == persistedTx.GetNetwork() &&
				tx.GetDate() == persistedTx.GetDate() &&
				tx.GetQuantity() == persistedTx.GetQuantity() &&
				tx.GetTotal() == persistedTx.GetTotal() {
				return false, nil
			}
		}
	}
	return true, nil
}

func (service *TransactionServiceImpl) UpdateCategory(id, category string) error {
	entity, err := service.dao.Get(id)
	if err != nil {
		service.ctx.GetLogger().Errorf("[TransactionService.UpdateCategory] Error updating %s's transaction id %s to category %s: %s",
			service.ctx.GetUser().GetUsername(), id, category, err.Error())
		return err
	}
	if entity.GetCategory() != category {
		return service.dao.Update(entity, "category", category)
	}
	return nil
}

func (service *TransactionServiceImpl) Synchronize() ([]common.Transaction, error) {
	service.ctx.GetLogger().Debugf("[TransactionService.Synchronize] Synchronizing transaction history for: %s", service.ctx.GetUser().GetUsername())
	var synchronized []common.Transaction
	persisted, err := service.dao.Find()
	if err != nil {
		return synchronized, err
	}
	/*
		transactions, err := service.ethereumService.GetTransactions()
		if err != nil {
			return nil, err
		}
	*/
	transactions, err := service.GetWalletHistory()
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, service.GetOrderHistory()...)
	transactions = append(transactions, service.GetDepositHistory()...)
	transactions = append(transactions, service.GetWithdrawalHistory()...)
	transactions = append(transactions, service.GetImportedTransactions()...)
	service.Sort(&transactions)
	for _, tx := range transactions {
		if unique, err := service.isUnique(tx, &transactions, &persisted); err != nil {
			return synchronized, err
		} else {
			if unique {
				err := service.dao.Create(service.mapper.MapTransactionDtoToEntity(tx))
				if err != nil {
					service.ctx.GetLogger().Errorf("Error adding transaction to database: %s. Error: %s", tx, err.Error())
					//return nil, err
					continue
				}
				synchronized = append(synchronized, tx)
			}
		}
	}
	return synchronized, nil
}

func (service *TransactionServiceImpl) GetHistory() ([]common.Transaction, error) {
	service.ctx.GetLogger().Debugf("[TransactionService.GetHistory] Retrieving transaction history for %s.",
		service.ctx.GetUser().GetUsername())
	entities, err := service.dao.Find()
	if err != nil {
		return nil, err
	}
	size := len(entities)
	if size > 0 {
		transactions := make([]common.Transaction, size)
		for i, entity := range entities {
			transactions[i] = service.mapper.MapTransactionEntityToDto(&entity)
		}
		service.Sort(&transactions)
		return transactions, nil
	}
	return service.Synchronize()
}

func (service *TransactionServiceImpl) GetWalletHistory() ([]common.Transaction, error) {
	var transactions []common.Transaction
	walletPlugins, err := service.userService.GetWalletPlugins()
	if err != nil {
		return transactions, err
	}
	for _, wallet := range walletPlugins {
		txs, err := wallet.GetTransactions()
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, txs...)
	}
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
		if ex.GetName() == "GDAX" || ex.GetName() == "Coinbase" { // TODO standardize exchange interface
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

func (service *TransactionServiceImpl) Sort(txs *[]common.Transaction) {
	service.ctx.GetLogger().Debugf("[TransactionService.sort] Sorting %d transactions", len(*txs))
	sort.Slice(*txs, func(i, j int) bool {
		return (*txs)[i].GetDate().After((*txs)[j].GetDate())
	})
}
