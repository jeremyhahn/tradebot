package reporting

import (
	"fmt"
	"sort"
	"strconv"
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

func (report *Report) Run(start, end time.Time) *Form8949 {

	var longs, shorts []Form8949LineItem
	lots := make(map[string][]common.Transaction)
	sales := make(map[string][]common.Transaction)

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
		/*
			if currency != "BTC" {
				continue
			}*/
		report.sort(&txs)
		lots[currency] = txs
		for _, tx := range txs {
			fmt.Printf("%s = %+v\n", currency, tx)
		}
	}

	for currency, salesTxs := range sales {

		fmt.Printf("PROCESSING %s\n", currency)

		if buyTxs, ok := lots[currency]; ok {

			var sublot *Sublot
			for _, saleTx := range salesTxs {

				util.DUMP("-----------New Sale--------------")

				buyQty, _ := decimal.NewFromString(buyTxs[0].GetQuantity())
				buyPrice, _ := decimal.NewFromString(buyTxs[0].GetFiatPrice())
				buyTotal, _ := decimal.NewFromString(buyTxs[0].GetFiatTotal())
				saleQty, _ := decimal.NewFromString(saleTx.GetQuantity())

				fmt.Printf("buyTx=%+v\n", buyTxs[0])
				fmt.Printf("buyQty=%s\n", buyQty)
				fmt.Printf("buyPrice=%s\n", buyPrice)
				fmt.Printf("buyTotal=%s\n", buyTotal)
				fmt.Printf("saleTx=%+v\n", saleTx)
				fmt.Printf("saleQty=%s\n\n", saleQty)

				if sublot != nil {

					util.DUMP("Procesing SUBLOT...")
					util.DUMP(sublot)

					if saleQty.GreaterThan(sublot.Quantity) {
						var txs []common.Transaction
						sum := sublot.Quantity
						i := 1
						for sum.LessThan(saleQty) {
							qty, _ := decimal.NewFromString(buyTxs[0].GetQuantity())
							sum = sum.Add(qty)
							txs = append(txs, buyTxs[0])
							report.ctx.GetLogger().Debugf("[FIFOReport.Run] qty=%s, sum=%s, txs size: %d, i=%d", qty, sum, len(txs), i)
							buyTxs = buyTxs[1:]
							fmt.Printf("buyTx=%+v\n", buyTxs[0])
							i++
						}
						_sublot, lineItems := report.calculateLotsWithSublot(&txs, saleTx, sublot)
						sublot = _sublot
						util.DUMP(sublot)
						util.DUMP("")
						util.DUMP(lineItems)
						util.DUMP("")
						for _, item := range lineItems {
							if report.isShortSale(&item) {
								shorts = append(shorts, item)
							} else {
								longs = append(longs, item)
							}
						}
					} else {
						_sublot, lineItem := report.calculateSublot(sublot, saleTx)
						sublot = _sublot
						util.DUMP(sublot)
						util.DUMP(lineItem)
						util.DUMP("")
						if report.isShortSale(lineItem) {
							shorts = append(shorts, *lineItem)
						} else {
							longs = append(longs, *lineItem)
						}
					}
					continue
				}

				if saleQty.LessThanOrEqual(buyQty) {
					_sublot, lineItem := report.calculate(buyTxs[0], saleTx)
					sublot = _sublot
					util.DUMP(lineItem)
					util.DUMP("")
					if report.isShortSale(lineItem) {
						shorts = append(shorts, *lineItem)
					} else {
						longs = append(longs, *lineItem)
					}
					buyTxs = buyTxs[1:]
				} else {
					util.DUMP("saleQty > buyQty... getting multiple txs")
					var txs []common.Transaction
					var sum decimal.Decimal
					for sum.LessThan(saleQty) {
						qty, _ := decimal.NewFromString(buyTxs[0].GetQuantity())
						sum = sum.Add(qty)
						txs = append(txs, buyTxs[0])
						report.ctx.GetLogger().Debugf("[FIFOReport.Run] qty=%s, sum=%s, fiatTotal:%s, txs size: %d",
							qty, sum, buyTxs[0].GetFiatTotal(), len(txs))
						buyTxs = buyTxs[1:]
					}
					_sublot, lineItem := report.calculateLots(&txs, saleTx)
					sublot = _sublot
					util.DUMP(lineItem)
					util.DUMP(sublot)
					util.DUMP("")
					if report.isShortSale(lineItem) {
						shorts = append(shorts, *lineItem)
					} else {
						longs = append(longs, *lineItem)
					}
				}
			}
		}
	}

	form := &Form8949{
		LongHolds:  longs,
		ShortHolds: shorts}
	form.sort()
	return form
}

func (report *Report) isShortSale(lineItem *Form8949LineItem) bool {
	diff := lineItem.DateSold.Sub(lineItem.DateAcquired)
	return int(diff.Hours()/24) > 365
}

func (report *Report) isTaxable(tx common.Transaction) bool {
	if tx.GetCategory() == common.TX_CATEGORY_LOST ||
		tx.GetCategory() == common.TX_CATEGORY_GIFT ||
		tx.GetCategory() == common.TX_CATEGORY_DONATION {
		return false
	}
	return true
}

func (report *Report) calculateSublot(lot *Sublot, saleTx common.Transaction) (*Sublot, *Form8949LineItem) {

	report.ctx.GetLogger().Debugf("[FIFOReport.calculateSublot] Calculate sublot (%s) for saleTx: %s\n", lot.Quantity, saleTx.GetQuantity())

	var sublot *Sublot
	saleQty, _ := decimal.NewFromString(saleTx.GetQuantity())

	if lot.Quantity.GreaterThan(saleQty) {
		subqty := lot.Quantity.Sub(saleQty)
		sublot = &Sublot{
			Date:      lot.Date,
			Quantity:  subqty,
			Price:     lot.Price,
			CostBasis: lot.Price.Mul(subqty)}
	}

	proceeds, _ := decimal.NewFromString(saleTx.GetFiatQuantity())
	costBasis := lot.Price.Mul(saleQty)

	return sublot, &Form8949LineItem{
		Currency:         saleTx.GetQuantityCurrency(),
		Description:      fmt.Sprintf("%s %s", saleTx.GetQuantity(), saleTx.GetQuantityCurrency()),
		DateAcquired:     lot.Date,
		DateSold:         saleTx.GetDate(),
		Proceeds:         proceeds.StringFixed(2),
		CostBasis:        costBasis.StringFixed(2),
		AdjustmentCode:   "",
		AdjustmentAmount: "",
		GainOrLoss:       proceeds.Sub(costBasis).StringFixed(2)}
}

func (report *Report) calculateLots(lots *[]common.Transaction, saleTx common.Transaction) (*Sublot, *Form8949LineItem) {
	report.ctx.GetLogger().Debugf("[FIFOReport.calculateLots] Calculating %d lots for saleTx: %s\n", len(*lots), saleTx.GetQuantity())
	var sublot *Sublot
	_lots := *lots
	lotLen, _ := decimal.NewFromString(strconv.Itoa(len(_lots)))
	dateAcquired := _lots[0].GetDate()
	saleQty, _ := decimal.NewFromString(saleTx.GetQuantity())
	proceeds, _ := decimal.NewFromString(saleTx.GetFiatQuantity())
	deductedQty := decimal.NewFromFloat(0)
	costBasis := decimal.NewFromFloat(0)
	for _, lot := range *lots {
		lotprice, _ := decimal.NewFromString(lot.GetFiatPrice())
		lotQty, _ := decimal.NewFromString(lot.GetQuantity())
		deductedQty = deductedQty.Add(lotQty)
		costBasis = costBasis.Add(lotprice.Mul(saleQty))
		if deductedQty.GreaterThan(saleQty) {
			subqty := deductedQty.Sub(saleQty)
			sublot = &Sublot{
				Date:      lot.GetDate(),
				Quantity:  subqty,
				Price:     lotprice,
				CostBasis: lotprice.Mul(subqty)}
			break
		}
	}
	costBasis = costBasis.Div(lotLen)
	return sublot, &Form8949LineItem{
		Currency:         saleTx.GetQuantityCurrency(),
		Description:      fmt.Sprintf("%s %s", saleTx.GetQuantity(), saleTx.GetQuantityCurrency()), // deductedQty
		DateAcquired:     dateAcquired,
		DateSold:         saleTx.GetDate(),
		Proceeds:         proceeds.StringFixed(2),
		CostBasis:        costBasis.StringFixed(2),
		AdjustmentCode:   "",
		AdjustmentAmount: "",
		GainOrLoss:       proceeds.Sub(costBasis).StringFixed(2)}
}

func (report *Report) calculateLotsWithSublot(lots *[]common.Transaction, saleTx common.Transaction, slot *Sublot) (*Sublot, []Form8949LineItem) {
	report.ctx.GetLogger().Debugf("[FIFOReport.calculateLotsWithSublot] Calculating %d lots with sublot %s for saleTx: %s\n",
		len(*lots), slot, saleTx.GetQuantity())

	util.DUMP("calculateLotsWithSublot")
	util.DUMP(slot)

	lineItems := make([]Form8949LineItem, 2)

	saleTotal, _ := decimal.NewFromString(saleTx.GetFiatPrice())
	slotProceeds := saleTotal.Mul(slot.Quantity)
	lineItems[0] = Form8949LineItem{
		Currency:         saleTx.GetQuantityCurrency(),
		Description:      fmt.Sprintf("%s %s", slot.Quantity, saleTx.GetQuantityCurrency()),
		DateAcquired:     slot.Date,
		DateSold:         saleTx.GetDate(),
		Proceeds:         slotProceeds.StringFixed(2),
		CostBasis:        slot.CostBasis.StringFixed(2),
		AdjustmentCode:   "",
		AdjustmentAmount: "",
		GainOrLoss:       slotProceeds.Sub(slot.CostBasis).StringFixed(2)}

	util.DUMP(lineItems[0])

	var sublot *Sublot
	_lots := *lots
	lotLen, _ := decimal.NewFromString(strconv.Itoa(len(_lots)))
	dateAcquired := _lots[0].GetDate()
	proceeds, _ := decimal.NewFromString(saleTx.GetFiatQuantity())
	proceeds = proceeds.Sub(slotProceeds)
	preSaleQty, _ := decimal.NewFromString(saleTx.GetQuantity())
	saleQty := preSaleQty.Sub(slot.Quantity)
	deductedQty := decimal.NewFromFloat(0)
	costBasis := decimal.NewFromFloat(0)
	for _, lot := range *lots {
		lotprice, _ := decimal.NewFromString(lot.GetFiatPrice())
		lotQty, _ := decimal.NewFromString(lot.GetQuantity())
		deductedQty = deductedQty.Add(lotQty)
		costBasis = costBasis.Add(lotprice.Mul(saleQty))
		if deductedQty.GreaterThan(saleQty) {
			subqty := deductedQty.Sub(saleQty)
			sublot = &Sublot{
				Date:      lot.GetDate(),
				Quantity:  subqty,
				Price:     lotprice,
				CostBasis: lotprice.Mul(subqty)}
			break
		}
	}
	costBasis = costBasis.Div(lotLen)
	lineItems[1] = Form8949LineItem{
		Currency:         saleTx.GetQuantityCurrency(),
		Description:      fmt.Sprintf("%s %s", deductedQty, saleTx.GetQuantityCurrency()),
		DateAcquired:     dateAcquired,
		DateSold:         slot.Date,
		Proceeds:         proceeds.StringFixed(2),
		CostBasis:        costBasis.StringFixed(2),
		AdjustmentCode:   "",
		AdjustmentAmount: "",
		GainOrLoss:       proceeds.Sub(costBasis).StringFixed(2)}

	util.DUMP(lineItems[1])

	return sublot, lineItems
}

func (report *Report) calculate(lot common.Transaction, saleTx common.Transaction) (*Sublot, *Form8949LineItem) {
	report.ctx.GetLogger().Debugf("[FIFOReport.calculate] Calculate single lot (%s) for saleTx: %s\n", lot.GetQuantity(), saleTx.GetQuantity())
	var sublot *Sublot
	lotQty, _ := decimal.NewFromString(lot.GetQuantity())
	lotPrice, _ := decimal.NewFromString(lot.GetFiatPrice())
	saleQty, _ := decimal.NewFromString(saleTx.GetQuantity())
	proceeds, _ := decimal.NewFromString(saleTx.GetFiatQuantity())
	costBasis := lotPrice.Mul(saleQty)
	if lotQty.GreaterThan(saleQty) {
		subqty := lotQty.Sub(saleQty)
		sublot = &Sublot{
			Date:      lot.GetDate(),
			Quantity:  subqty,
			Price:     lotPrice,
			CostBasis: lotPrice.Mul(subqty)}
	}
	return sublot, &Form8949LineItem{
		Currency:         saleTx.GetQuantityCurrency(),
		Description:      fmt.Sprintf("%s %s", saleTx.GetQuantity(), saleTx.GetQuantityCurrency()),
		DateAcquired:     lot.GetDate(),
		DateSold:         saleTx.GetDate(),
		Proceeds:         proceeds.StringFixed(2),
		CostBasis:        costBasis.StringFixed(2),
		AdjustmentCode:   "",
		AdjustmentAmount: "",
		GainOrLoss:       proceeds.Sub(costBasis).StringFixed(2)}
}

func (report *Report) sort(txs *[]common.Transaction) {
	report.ctx.GetLogger().Debugf("[FIFOReport.Sort] Sorting %d transactions", len(*txs))
	sort.Slice(*txs, func(i, j int) bool {
		return (*txs)[i].GetDate().Before((*txs)[j].GetDate()) ||
			(*txs)[i].GetDate().Equal((*txs)[j].GetDate())
	})
}

/*func (report *Report) sortLots(lots map[string][]common.Transaction) map[string][]common.Transaction {
	exists := make(map[string]bool, len(lots))
	keys := make([]string, len(lots))
	for k, _ := range lots {
		if _, ok := exists[k]; !ok {
			keys = append(keys, k)
			exists[k] = true
		}
	}
	sort.Strings(keys)
	sorted := make(map[string][]common.Transaction)
	for _, k := range keys {
		sorted[k] = lots[k]
	}
	return sorted
}*/

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
