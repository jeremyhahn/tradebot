package main

/*
func TestCoinbase_GetBalance(t *testing.T) {
	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")
	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	coins, sum := cb.GetBalances()

	assert.Equal(t, true, len(coins) > 0)
	assert.Equal(t, true, sum.GreaterThan(decimal.NewFromFloat(0)))

	test.CleanupIntegrationTest()
}

func TestCoinbase_GetPriceAt(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	atDate := time.Date(2014, 03, 22, 23, 38, 8, 0, time.Now().Location())
	candle, err := cb.GetPriceAt("BTC", atDate)

	assert.Nil(t, err)
	assert.Equal(t, atDate, candle.Date)
	assert.Equal(t, true, candle.Close.GreaterThan(decimal.NewFromFloat(0)))
}

func TestCoinbase_GetOrderHistory(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	orders := cb.GetOrderHistory(&common.CurrencyPair{
		Base:          "BTC",
		Quote:         "USD",
		LocalCurrency: "USD"})

	assert.Nil(t, err)
	assert.Equal(t, true, len(orders) > 0)

	for _, o := range orders {
		util.DUMP(o)
	}

	test.CleanupIntegrationTest()
}

func TestCoinbase_GetDeposits(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	deposits, err := cb.GetDepositHistory()
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, true, len(deposits) > 0)

	test.CleanupIntegrationTest()
}

func TestCoinbase_GetWithdrawls(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	withdrawls, err := cb.GetWithdrawalHistory()
	assert.Nil(t, err)
	assert.Equal(t, true, len(withdrawls) > 0)

	for _, w := range withdrawls {
		util.DUMP(w)
	}

	test.CleanupIntegrationTest()
}

/*
func TestCoinbase_GetCurrencies(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	currencies, err := cb.GetCurrencies()
	assert.Nil(t, err)
	assert.Equal(t, true, len(currencies) > 0)

	test.CleanupIntegrationTest()
}
*/

/*
func TestCoinbase_GetTransactions(t *testing.T) {

	ctx := test.CreateIntegrationTestContext("../../../.env", "../../../")

	userDAO := dao.NewUserDAO(ctx)
	userEntity := &entity.User{Id: ctx.GetUser().GetId()}

	cryptoExchange, err := userDAO.GetExchange(userEntity, "Coinbase")
	assert.Nil(t, err)

	cb := CreateCoinbase(ctx, cryptoExchange).(*Coinbase)
	txs, err := cb.getTransactions()
	assert.Nil(t, err)
	assert.Equal(t, true, len(txs) > 0)

	for _, t := range txs {
		util.DUMP(t)
	}

	util.DUMP(len(txs))

	test.CleanupIntegrationTest()
}
*/
