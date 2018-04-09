package reporting

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/shopspring/decimal"
)

type Report struct {
	ctx         common.Context
	Deposits    []common.Transaction
	Withdrawals []common.Transaction
	Trades      []common.Transaction
	Income      []common.Transaction
	Spends      []common.Transaction
}

type Form8949 struct {
	LongHolds  []Form8949LineItem
	ShortHolds []Form8949LineItem
}

type Form8949LineItem struct {
	Description      string
	DateAcquired     string
	DateSold         string
	Proceeds         string
	CostBasis        string
	AdjustmentCode   string
	AdjustmentAmount string
	GainOrLoss       string
}

type ReportParams struct {
	ctx common.Context
}

func NewFifoReport(ctx common.Context, transactions []common.Transaction) *Report {

	var trades []common.Transaction
	var deposits []common.Transaction
	var withdrawals []common.Transaction
	var income []common.Transaction
	var spends []common.Transaction

	for _, tx := range transactions {
		switch {
		case tx.GetType() == common.DEPOSIT_ORDER_TYPE:
			deposits = append(deposits, tx)
		case tx.GetType() == common.WITHDRAWAL_ORDER_TYPE:
			withdrawals = append(withdrawals, tx)
		case tx.GetType() == common.BUY_ORDER_TYPE || tx.GetType() == common.SELL_ORDER_TYPE:
			trades = append(trades, tx)
		case tx.GetCategory() == common.TX_CATEGORY_MINING || tx.GetCategory() == common.TX_CATEGORY_INCOME:
			income = append(income, tx)
		case tx.GetCategory() == common.TX_CATEGORY_SPEND || tx.GetCategory() == common.TX_CATEGORY_DONATION ||
			tx.GetCategory() == common.TX_CATEGORY_GIFT || tx.GetCategory() == common.TX_CATEGORY_LOST:
			spends = append(spends, tx)
		}
	}

	report := &Report{
		ctx:         ctx,
		Deposits:    deposits,
		Withdrawals: withdrawals,
		Trades:      trades,
		Income:      income,
		Spends:      spends}

	report.sort(&report.Deposits)
	report.sort(&report.Withdrawals)
	report.sort(&report.Trades)
	report.sort(&report.Income)
	report.sort(&report.Spends)

	return report
}

func (report *Report) isTaxable(tx common.Transaction) bool {
	if tx.GetCategory() == common.TX_CATEGORY_LOST ||
		tx.GetCategory() == common.TX_CATEGORY_GIFT ||
		tx.GetCategory() == common.TX_CATEGORY_DONATION {
		return false
	}
	return true
}

func (report *Report) Run(start, end time.Time) {

	zero := decimal.NewFromFloat(0)
	unidentified := make(map[string][]common.Transaction)
	lots := make(map[string][]common.Transaction)
	sales := make(map[string][]common.Transaction)
	var shorts []Form8949LineItem
	var longs []Form8949LineItem

	for _, trade := range report.Trades {
		if trade.GetDate().After(end) || !report.isTaxable(trade) {
			continue
		}
		if trade.GetType() == common.BUY_ORDER_TYPE {
			currency := trade.GetCurrencyPair().Base
			lots[currency] = append(lots[currency], trade)
		}
		if trade.GetType() == common.SELL_ORDER_TYPE {
			currency := trade.GetCurrencyPair().Base
			sales[currency] = append(sales[currency], trade)
		}
	}
	for _, deposit := range report.Deposits {
		if deposit.GetDate().After(end) || !report.isTaxable(deposit) {
			continue
		}
		currency := deposit.GetCurrencyPair().Base
		lots[currency] = append(lots[currency], deposit)
	}

	for currency, txs := range lots {
		if currency != "BTC" {
			continue
		}
		report.sort(&txs)
		lots[currency] = txs
		for _, tx := range txs {
			fmt.Printf("%s = %+v\n", currency, tx)
		}
	}

	for currency, salesTxs := range sales {

		if currency != "BTC" {
			continue
		}

		if buyTxs, ok := lots[currency]; ok {

			var deductedQty, costBasis, buyFees, saleTotal, gainOrLoss decimal.Decimal
			remainder := make(map[string]*Remainder)
			needsMore := make(map[string]*NeedsMore)

			for i, saleTx := range salesTxs {

				processed := false
				buyTx := buyTxs[0]
				buyQty, _ := decimal.NewFromString(buyTx.GetQuantity())
				buyTotal, _ := decimal.NewFromString(buyTx.GetFiatTotal())
				buyFee, _ := decimal.NewFromString(buyTx.GetFiatFee())
				buyFees = buyFees.Add(buyFee)
				dateAcquired := buyTx.GetDate()
				saleQty, _ := decimal.NewFromString(saleTx.GetQuantity())
				/*
					salePrice, _ := decimal.NewFromString(saleTx.GetFiatTotal())
					saleFee, _ := decimal.NewFromString(saleTx.GetFiatFee())
					saleTotal := saleTx.GetFiatTotal()
				*/
				fmt.Printf("buyTx=%+v\n", buyTx)
				fmt.Printf("buyQty=%+v\n", buyQty)
				fmt.Printf("buyTotal=%s\n", buyTotal)
				fmt.Printf("buyFee=%s\n", buyFee)
				fmt.Printf("buyFees=%s\n", buyFees)
				fmt.Printf("saleTx=%+v\n", saleTx)
				fmt.Printf("saleQty=%s\n\n", saleQty)
				/*
					fmt.Printf("salePrice=%s\n", salePrice)
					fmt.Printf("saleFee=%s\n", saleFee)
					fmt.Printf("saleTotal=%s\n", saleTotal)
					fmt.Printf("costBasis=%s\n", costBasis)
				*/

				if r, ok := remainder[currency]; ok {
					util.DUMP("remainder!")

					dateAcquired = r.Tx.GetDate()
					deductedQty = deductedQty.Add(r.Quantity)
					costBasis = r.CostBasis
					saleTotal, _ = decimal.NewFromString(saleTx.GetFiatTotal())
					gainOrLoss = saleTotal.Sub(costBasis)
					fmt.Printf("deducted from remainder: %s\n", deductedQty)

					if deductedQty.LessThan(saleQty) {
						fmt.Printf("remainder - deducted buyQty: %s\n", buyQty)

						needs := saleQty.Sub(deductedQty)

						fmt.Printf("remainder- needs: %s\n", needs)

						deductedQty = deductedQty.Add(buyQty)

						fmt.Printf("remainder- deducted: %s\n", deductedQty)
					}

					r.Quantity = saleQty.Sub(deductedQty)

					fmt.Printf("remainder (new value): %+v\n", remainder)

					if r.Quantity.LessThanOrEqual(zero) {
						delete(remainder, currency)
						fmt.Printf("remainder (after delete): %+v\n", remainder)
					}

					processed = true
				}

				if !processed {
					if nm, ok := needsMore[currency]; ok {
						util.DUMP("needs more!")
						if saleQty.GreaterThan(nm.Needed) || saleQty.Equals(nm.Needed) {

							dateAcquired = nm.SalesTx.GetDate()
							saleQty = nm.SaleQuantity
							saleTotal, _ = decimal.NewFromString(nm.SalesTx.GetFiatTotal())
							//	saleFee, _ := decimal.NewFromString(nm.SalesTx.GetFiatFee())
							costBasis = costBasis.Add(saleTotal) //.Add(saleFee).Add(nm.PurchaseFee)
							gainOrLoss = saleTotal.Sub(costBasis)
							deductedQty = deductedQty.Add(nm.Needed)

							fmt.Printf("needsMore (needed): %s\n", nm.Needed)
							fmt.Printf("needsMore deducted (total): %s\n", deductedQty)

							needsMore[currency] = &NeedsMore{
								SaleQuantity:  nm.SaleQuantity,
								Needed:        nm.SaleQuantity.Sub(deductedQty),
								SalesTx:       saleTx,
								PurchasePrice: buyTotal}
							fmt.Printf("needsMore (new value): %+v\n", needsMore)

							if needsMore[currency].Needed.LessThanOrEqual(zero) {
								remainder[currency] = &Remainder{
									Quantity:  buyQty.Sub(nm.Needed),
									CostBasis: costBasis,
									Tx:        buyTx}
								fmt.Printf("needsMore remainder: %+v\n", remainder[currency])
								delete(needsMore, currency)
							} else {
								continue
							}

							processed = true
						} else {
							remainder[currency] = &Remainder{
								Quantity:  needsMore[currency].SaleQuantity.Sub(saleQty),
								CostBasis: costBasis,
								Tx:        buyTx}
							fmt.Printf("needsMore remainder: %+v\n", remainder)
							deductedQty = deductedQty.Add(saleQty)
							fmt.Printf("needsMore (new value 2): %s\n", deductedQty.String())

							buyTxs = buyTxs[1:]
							continue
						}
					}
				}

				if !processed {
					if saleQty.GreaterThan(buyQty) || saleQty.Equals(buyQty) {

						deductedQty = deductedQty.Add(buyQty)
						fmt.Printf("deducted: %s\n", deductedQty.String())

						needsMore[currency] = &NeedsMore{
							Needed:        saleQty.Sub(deductedQty),
							SaleQuantity:  saleQty,
							SalesTx:       saleTx,
							PurchasePrice: buyTotal}
						fmt.Printf("needsMore: %+v\n", needsMore)

						buyTxs = buyTxs[1:]
						continue
					} else {

						//saleFee, _ := decimal.NewFromString(saleTx.GetFiatFee())
						saleTotal, _ = decimal.NewFromString(saleTx.GetFiatTotal())
						costBasis = saleQty.Div(buyQty).Mul(buyTotal) //.Add(buyFee).Add(saleFee)
						gainOrLoss = saleTotal.Sub(costBasis)

						remainder[currency] = &Remainder{
							Quantity:  buyQty.Sub(saleQty),
							CostBasis: costBasis,
							Tx:        buyTx}

						fmt.Printf("remainder: %+v\n", remainder)

						deductedQty = deductedQty.Add(saleQty)
						fmt.Printf("deducted: %s\n", deductedQty.String())

						buyTxs = buyTxs[1:]
					}
				}

				lineItem := Form8949LineItem{
					Description:      saleQty.StringFixed(8),
					DateAcquired:     dateAcquired.String(),
					DateSold:         saleTx.GetDate().String(),
					Proceeds:         saleTotal.StringFixed(2),
					CostBasis:        costBasis.StringFixed(2),
					AdjustmentCode:   "",
					AdjustmentAmount: "",
					GainOrLoss:       gainOrLoss.StringFixed(2)}

				util.DUMP(lineItem)

				deductedQty = decimal.NewFromFloat(0)
				costBasis = decimal.NewFromFloat(0)
				buyFees = decimal.NewFromFloat(0)
				saleTotal = decimal.NewFromFloat(0)
				gainOrLoss = decimal.NewFromFloat(0)

				if i == 4 {
					os.Exit(0)
				}

				if dateAcquired.Before(time.Now().AddDate(-1, 0, 0)) {
					longs = append(longs, lineItem)
				} else {
					shorts = append(shorts, lineItem)
				}

			}

		} else {
			for _, tx := range buyTxs {
				unidentified[currency] = append(unidentified[currency], tx)
			}
		}
	}

}

func (report *Report) sort(txs *[]common.Transaction) {
	report.ctx.GetLogger().Debugf("[FIFOReport.Sort] Sorting %d transactions", len(*txs))
	sort.Slice(*txs, func(i, j int) bool {
		if (*txs)[i].GetDate().Equal((*txs)[j].GetDate()) {
			return true
		}
		return (*txs)[i].GetDate().Before((*txs)[j].GetDate())
	})
}

func (report *Report) sum(txs []common.Transaction) decimal.Decimal {
	var sum decimal.Decimal
	for _, tx := range txs {
		quantity, err := decimal.NewFromString(tx.GetQuantity())
		if err != nil {
			report.ctx.GetLogger().Errorf("[FIFOReport.sum] Error: %s", err.Error())
		}
		sum = sum.Add(quantity)
	}
	return sum
}

/*
func WriteCSV() {
	  filename := fmt.Sprintf("/tmp/%s-%s", ctx.GetUser().GetUsername(), "fifo.csv")
		file, err := os.Create(filename)
		if err != nil {
			return "", err
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

			ctx.ResponseWriter.Header().Set("Content-Type", "text/csv")
			ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
			ctx.ResponseWriter.Write(b.Bytes())
}
*/

func createTransactionService(ctx common.Context) service.TransactionService {
	pluginDAO := dao.NewPluginDAO(ctx)
	userDAO := dao.NewUserDAO(ctx)
	transactionDAO := dao.NewTransactionDAO(ctx)
	userMapper := mapper.NewUserMapper()
	pluginMapper := mapper.NewPluginMapper()
	transactionMapper := mapper.NewTransactionMapper(ctx)
	userExchangeMapper := mapper.NewUserExchangeMapper()
	marketcapService := service.NewMarketCapService(ctx)
	pluginService := service.CreatePluginService(ctx, "../plugins/", pluginDAO, pluginMapper)
	exchangeService := service.NewExchangeService(ctx, userDAO, userMapper, userExchangeMapper, pluginService)
	ethereumService, _ := service.NewEthereumService(ctx, userDAO, userMapper, marketcapService, exchangeService)
	fiatPriceService, _ := service.NewFiatPriceService(ctx, exchangeService)
	walletService := service.NewWalletService(ctx, pluginService, fiatPriceService)
	userService := service.NewUserService(ctx, userDAO, userMapper, userExchangeMapper, marketcapService, ethereumService, exchangeService, walletService)
	return service.NewTransactionService(ctx, transactionDAO, transactionMapper,
		exchangeService, userService, ethereumService, fiatPriceService)
}
