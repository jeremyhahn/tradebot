package service

/*
func TestChartService_Stream(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	if len(marketcap.GetMarkets()) <= 0 {
		t.Fatal("[TestMarketCapService.GetMarkets] Unable to get market cap list")
	}

	if len(marketcap.GetMarkets()) <= 0 {
		t.Fatal("[TestMarketCapService.GetMarkets] Unable to get market cap list")
	}

	test.CleanupMockContext()
}

func TestChartService_GetLastTrade(t *testing.T) {
	ctx := test.NewTestContext()
	now := time.Now()
	sampleTrades := make([]dao.Trade, 0, 5)
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, -1, 0),
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    15000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -20),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    16000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -15),
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    12000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -5),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    19000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     now,
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    9000,
		UserID:   ctx.User.Id})

	autoTradeCoin := &dao.AutoTradeCoin{
		UserID:   ctx.User.Id,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Period:   900,
		Trades:   sampleTrades}

	autoTradeDAO := dao.NewAutoTradeDAO(ctx)
	autoTradeDAO.Save(autoTradeCoin)
	trade := autoTradeDAO.GetLastTrade(autoTradeCoin)

	if trade.AutoTradeID != 1 && trade.Date == now {
		t.Fatal("[TestChartService_GetLastTrade] Failed to get expected last trade")
	}

	trades := autoTradeDAO.FindByCurrency(ctx.User, &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"})

	if len(trades) != 5 {
		t.Fatal("[TestChartService_GetLastTrade] Failed to get list of expected trades")
	}

	test.CleanupMockContext()
}

func TestChartService_GetLastTrade(t *testing.T) {
	ctx := test.NewTestContext()
	now := time.Now()
	sampleTrades := make([]dao.Trade, 0, 5)
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, -1, 0),
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    15000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -20),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    16000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -15),
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    12000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     time.Now().AddDate(0, 0, -5),
		Type:     "sell",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    19000,
		UserID:   ctx.User.Id})
	sampleTrades = append(sampleTrades, dao.Trade{
		Date:     now,
		Type:     "buy",
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Amount:   1,
		Price:    9000,
		UserID:   ctx.User.Id})

	autoTradeCoin := &dao.AutoTradeCoin{
		UserID:   ctx.User.Id,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Period:   900,
		Trades:   sampleTrades}

	autoTradeDAO := dao.NewAutoTradeDAO(ctx)
	autoTradeDAO.Save(autoTradeCoin)
	trade := autoTradeDAO.GetLastTrade(autoTradeCoin)

	if trade.AutoTradeID != 1 && trade.Date == now {
		t.Fatal("[TestChartService_GetLastTrade] Failed to get expected last trade")
	}

	trades := autoTradeDAO.FindByCurrency(ctx.User, &common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"})

	if len(trades) != 5 {
		t.Fatal("[TestChartService_GetLastTrade] Failed to get list of expected trades")
	}

	test.CleanupMockContext()
}
*/
