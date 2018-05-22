package accounting

import (
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/mapper"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/shopspring/decimal"
)

type Report struct {
	ctx          common.Context
	Transactions []common.Transaction
	Deposits     []common.Transaction
	Withdrawals  []common.Transaction
	Trades       []common.Transaction
	Income       []common.Transaction
	Spends       []common.Transaction
}

func NewFifoReport(ctx common.Context, transactions []common.Transaction) *Report {
	return &Report{
		ctx:          ctx,
		Transactions: transactions}
}

func (report *Report) Run(start, end time.Time) *Form8949 {

	buyLots := make(map[string][]Coinlot)
	saleLots := make(map[string][]Coinlot)

	for i := len(report.Transactions) - 1; i >= 0; i-- {
		trade := report.Transactions[i]

		//if trade.GetDate().After(end) || !report.isTaxable(trade) || !strings.Contains(trade.GetCurrencyPair().String(), "BTC") {
		if trade.GetDate().After(end) || !report.isTaxable(trade) {
			continue
		}

		util.DUMP(trade)

		txType := trade.GetType()
		quantity, _ := decimal.NewFromString(trade.GetQuantity())
		total, _ := decimal.NewFromString(trade.GetTotal())
		fiatQuantity, _ := decimal.NewFromString(trade.GetFiatQuantity())
		fiatPrice, _ := decimal.NewFromString(trade.GetFiatPrice())
		quoteFiatPrice, _ := decimal.NewFromString(trade.GetQuoteFiatPrice())
		fiatTotal, _ := decimal.NewFromString(trade.GetFiatTotal())

		baseCurrency := trade.GetCurrencyPair().Base
		quoteCurrency := trade.GetCurrencyPair().Quote

		if txType == common.DEPOSIT_ORDER_TYPE {
			coinlot := Coinlot{
				Date:      trade.GetDate(),
				Currency:  baseCurrency,
				Quantity:  quantity,
				UnitPrice: fiatPrice,
				CostBasis: fiatTotal}
			buyLots[baseCurrency] = append(buyLots[baseCurrency], coinlot)
			continue
		}

		if txType == common.BUY_ORDER_TYPE {

			buyCoinlot := Coinlot{
				Date:      trade.GetDate(),
				Currency:  baseCurrency,
				Quantity:  quantity,
				UnitPrice: quoteFiatPrice,
				CostBasis: fiatTotal}
			buyLots[baseCurrency] = append(buyLots[baseCurrency], buyCoinlot)

			if _, ok := common.FiatCurrencies[quoteCurrency]; ok {
				continue
			}
			saleCoinlot := Coinlot{
				Date:      trade.GetDate(),
				Currency:  quoteCurrency,
				Quantity:  total,
				UnitPrice: quoteFiatPrice,
				SalePrice: fiatQuantity,
				CostBasis: fiatTotal}
			saleLots[quoteCurrency] = append(saleLots[quoteCurrency], saleCoinlot)
		}

		if txType == common.SELL_ORDER_TYPE {

			saleCoinlot := Coinlot{
				Date:      trade.GetDate(),
				Currency:  baseCurrency,
				Quantity:  quantity,
				UnitPrice: quoteFiatPrice,
				SalePrice: fiatQuantity,
				CostBasis: fiatTotal}
			saleLots[baseCurrency] = append(saleLots[baseCurrency], saleCoinlot)

			buyCoinlot := Coinlot{
				Date:      trade.GetDate(),
				Currency:  quoteCurrency,
				Quantity:  total,
				UnitPrice: quoteFiatPrice.Mul(total),
				CostBasis: fiatTotal}
			buyLots[quoteCurrency] = append(buyLots[quoteCurrency], buyCoinlot)
		}

	}

	var shorts, longs []Form8949LineItem

	/*
		fmt.Println("buys")
		for _, buy := range buyLots["BTC"] {
			fmt.Printf("%+v\n", buy)
		}
		fmt.Println("sales")
		for _, sale := range saleLots["BTC"] {
			fmt.Printf("%+v\n", sale)
		}
		_shorts, _longs := report.process(buyLots["BTC"], saleLots["BTC"])
		shorts = append(shorts, _shorts...)
		longs = append(longs, _longs...)
	*/

	for currency, _ := range buyLots {
		_shorts, _longs := report.process(buyLots[currency], saleLots[currency])
		shorts = append(shorts, _shorts...)
		longs = append(longs, _longs...)
	}

	form := &Form8949{
		ShortHolds: shorts,
		LongHolds:  longs}
	form.sort()
	return form
}

func (report *Report) process(buyLots, saleLots []Coinlot) ([]Form8949LineItem, []Form8949LineItem) {

	var sublot *Coinlot
	var lineItem Form8949LineItem
	var shorts, longs []Form8949LineItem

	zero := decimal.NewFromFloat(0)

	for _, saleLot := range saleLots {

		if sublot != nil {

			if len(buyLots) == 0 {
				report.ctx.GetLogger().Debugf("[FIFOReport.process] Out of buy lots with remaining saleLot: %s\n", saleLot)
				//closingPosition[sublot.Currency] = sublot
				continue
			}

			report.ctx.GetLogger().Debugf("[FIFOReport.process] Processing NEW sublot: %s\n", sublot)
			report.ctx.GetLogger().Debugf("[FIFOReport.process] buyLots[0]: %+v\n", buyLots[0])
			report.ctx.GetLogger().Debugf("[FIFOReport.process] saleLot=%+v\n", saleLot)

			if saleLot.Quantity.GreaterThan(sublot.Quantity) {

				report.ctx.GetLogger().Debug("[FIFOReport.process] Sale lot quantity is greater than or equal to sublot quantity")

				proceeds := sublot.Quantity.Mul(saleLot.UnitPrice)
				lineItem = Form8949LineItem{
					Currency:         sublot.Currency,
					Description:      fmt.Sprintf("%s %s", sublot.Quantity, sublot.Currency),
					DateAcquired:     sublot.Date,
					DateSold:         saleLot.Date,
					Proceeds:         proceeds.StringFixed(2),
					CostBasis:        sublot.CostBasis.StringFixed(2),
					AdjustmentCode:   "",
					AdjustmentAmount: "",
					GainOrLoss:       proceeds.Sub(sublot.CostBasis).StringFixed(2)}
				report.ctx.GetLogger().Debugf("[FIFOReport.process] (sublot) lineItem: %+v\n", lineItem)
				if report.isShortSale(&lineItem) {
					shorts = append(shorts, lineItem)
				} else {
					longs = append(longs, lineItem)
				}

				remainder := saleLot.Quantity.Sub(sublot.Quantity)
				if remainder.GreaterThan(zero) {
					report.ctx.GetLogger().Debugf("[FIFOReport.process] (sublot) remainder=%s", remainder)
					proceeds := remainder.Mul(saleLot.UnitPrice)
					lineItem = Form8949LineItem{
						Currency:         saleLot.Currency,
						Description:      fmt.Sprintf("%s %s", remainder, saleLot.Currency),
						DateAcquired:     saleLot.Date,
						DateSold:         saleLot.Date,
						Proceeds:         proceeds.StringFixed(2),
						CostBasis:        saleLot.UnitPrice.Mul(remainder).StringFixed(2),
						AdjustmentCode:   "",
						AdjustmentAmount: "",
						GainOrLoss:       proceeds.Sub(saleLot.CostBasis).StringFixed(2)}
					if report.isShortSale(&lineItem) {
						shorts = append(shorts, lineItem)
					} else {
						longs = append(longs, lineItem)
					}
					report.ctx.GetLogger().Debugf("[FIFOReport.process] (sublot) lineItem=%s", lineItem)
				}
				sublot = nil

			} else {

				report.ctx.GetLogger().Debug("Sublot quantity is greater than or equal to the sale lot quantity")

				sublot, lineItem = report.calculate(sublot, &saleLot)
				report.ctx.GetLogger().Debugf("[FIFOReport.process] (sublot-else) (calculate) lineItem: %+v\n", lineItem)
				report.ctx.GetLogger().Debugf("[FIFOReport.process] (sublot-else) (calculate) sublot: %+v\n", sublot)

				if report.isShortSale(&lineItem) {
					shorts = append(shorts, lineItem)
				} else {
					longs = append(longs, lineItem)
				}

			}

			if sublot != nil && len(buyLots) > 0 {
				buyLots = buyLots[1:]
			}

			continue
		}

		report.ctx.GetLogger().Debugf("[FIFOReport.process] -- New Sale Lot --")
		report.ctx.GetLogger().Debugf("sublot=%+v\n", sublot)
		report.ctx.GetLogger().Debugf("buyLot=%+v\n", buyLots[0])
		report.ctx.GetLogger().Debugf("saleLot=%+v\n", saleLot)

		if saleLot.Quantity.LessThanOrEqual(buyLots[0].Quantity) {
			sublot, lineItem = report.calculate(&buyLots[0], &saleLot)
			report.ctx.GetLogger().Debugf("[FIFOReport.process] (calculate) lineItem: %+v\n", lineItem)
			report.ctx.GetLogger().Debugf("[FIFOReport.process] (calculate) sublot: %+v\n", sublot)
			if report.isShortSale(&lineItem) {
				shorts = append(shorts, lineItem)
			} else {
				longs = append(longs, lineItem)
			}
			buyLots = buyLots[1:]

		} else {

			report.ctx.GetLogger().Debugf("Need multiple lots to fulfil saleLot: %+v\n", saleLot)

			var lots []Coinlot
			var sum decimal.Decimal
			for sum.LessThan(saleLot.Quantity) {
				fmt.Printf("Adding buyLot=%+v\n", buyLots[0])
				sum = sum.Add(buyLots[0].Quantity)
				lots = append(lots, buyLots[0])

				report.ctx.GetLogger().Debugf("[FIFOReport.process] sale quantity=%s, buy quantity=%s, sum=%s, costBasis:%s, lot size: %d",
					saleLot.Quantity, buyLots[0].Quantity, sum, buyLots[0].CostBasis, len(lots))

				buyLots = buyLots[1:]
			}
			sublot, lineItem = report.calculateLots(&lots, &saleLot)
			report.ctx.GetLogger().Debugf("[FIFOReport.process] (calculateLots) lineItem: %+v\n", lineItem)
			report.ctx.GetLogger().Debugf("[FIFOReport.process] (calculateLots) sublot: %+v\n", sublot)
			if report.isShortSale(&lineItem) {
				shorts = append(shorts, lineItem)
			} else {
				longs = append(longs, lineItem)
			}
		}

	}

	return shorts, longs
}

func (report *Report) calculate(buyLot *Coinlot, saleLot *Coinlot) (*Coinlot, Form8949LineItem) {
	report.ctx.GetLogger().Debugf("[FIFOReport.calculate] Calculating single lot (%s) against sale: %s\n", buyLot.Quantity, saleLot.Quantity)
	report.ctx.GetLogger().Debugf("[FIFOReport.calculate] buyLot: %s\n", buyLot)
	report.ctx.GetLogger().Debugf("[FIFOReport.calculate] saleLot: %s\n", saleLot)

	var sublot *Coinlot

	if buyLot.Quantity.GreaterThan(saleLot.Quantity) {
		subqty := buyLot.Quantity.Sub(saleLot.Quantity)
		sublot = &Coinlot{
			Date:      buyLot.Date,
			Currency:  saleLot.Currency,
			Quantity:  subqty,
			UnitPrice: buyLot.UnitPrice,
			CostBasis: buyLot.UnitPrice.Mul(subqty)}
	}

	costBasis := buyLot.UnitPrice.Mul(saleLot.Quantity)
	return sublot, Form8949LineItem{
		Currency:         saleLot.Currency,
		Description:      fmt.Sprintf("%s %s", saleLot.Quantity, saleLot.Currency),
		DateAcquired:     buyLot.Date,
		DateSold:         saleLot.Date,
		Proceeds:         saleLot.SalePrice.StringFixed(2),
		CostBasis:        costBasis.StringFixed(2),
		AdjustmentCode:   "",
		AdjustmentAmount: "",
		GainOrLoss:       saleLot.CostBasis.Sub(costBasis).StringFixed(2)}
}

func (report *Report) calculateLots(buyLots *[]Coinlot, saleLot *Coinlot) (*Coinlot, Form8949LineItem) {
	report.ctx.GetLogger().Debugf("[FIFOReport.calculateLots] Calculating %d lots\n", len(*buyLots))

	for _, l := range *buyLots {
		report.ctx.GetLogger().Debugf("[FIFOReport.calculate] buyLot: %s\n", l)
	}
	report.ctx.GetLogger().Debugf("[FIFOReport.calculate] saleLot: %s\n", saleLot)

	var sublot *Coinlot
	var _lots = *buyLots
	var deductedQty, costBasis decimal.Decimal
	dateAcquired := _lots[0].Date

	for _, lot := range *buyLots {

		deductedQty = deductedQty.Add(lot.Quantity)
		//basis := lot.UnitPrice.Mul(saleLot.Quantity)
		basis := lot.CostBasis
		costBasis = costBasis.Add(basis)

		if deductedQty.GreaterThan(saleLot.Quantity) {
			subqty := deductedQty.Sub(saleLot.Quantity)
			sublot = &Coinlot{
				Date:      lot.Date,
				Currency:  saleLot.Currency,
				Quantity:  subqty,
				UnitPrice: lot.UnitPrice,
				CostBasis: basis}
			break
		}
	}
	return sublot, Form8949LineItem{
		Currency:         saleLot.Currency,
		Description:      fmt.Sprintf("%s %s", saleLot.Quantity, saleLot.Currency),
		DateAcquired:     dateAcquired,
		DateSold:         saleLot.Date,
		Proceeds:         saleLot.SalePrice.StringFixed(2),
		CostBasis:        costBasis.StringFixed(2),
		AdjustmentCode:   "",
		AdjustmentAmount: "",
		GainOrLoss:       saleLot.SalePrice.Sub(costBasis).StringFixed(2)}
	//GainOrLoss: saleLot.CostBasis.Sub(costBasis).StringFixed(2)}
}

func (report *Report) isShortSale(lineItem *Form8949LineItem) bool {
	diff := lineItem.DateSold.Sub(lineItem.DateAcquired)
	return int(diff.Hours()/24) < 365
}

func (report *Report) isTaxable(tx common.Transaction) bool {
	if tx.GetCategory() == common.TX_CATEGORY_TRANSFER ||
		tx.GetCategory() == common.TX_CATEGORY_LOST ||
		tx.GetCategory() == common.TX_CATEGORY_GIFT ||
		tx.GetCategory() == common.TX_CATEGORY_DONATION {
		return false
	}
	return true
}

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
