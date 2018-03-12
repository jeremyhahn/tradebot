package service

/*
func TestBlockCypher_GetBalance(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	marketcapService := NewMarketCapService(ctx)
	service := NewBlockCypherService(ctx, marketcapService)
	balance := service.GetBalance(os.Getenv("ETH_ADDRESS"))

	util.DUMP(balance)

	assert.Equal(t, true, len(balance.GetAddress()) > 0)
	assert.Equal(t, true, balance.GetBalance() > 0)
	assert.Equal(t, true, balance.GetValue() > 0)

	test.CleanupIntegrationTest()
}

func TestBlockCypher_GetTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service := NewBlockCypherService(ctx, NewMarketCapService(ctx))
	transactions, err := service.GetTransactions(os.Getenv("ETH_ADDRESS"))

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	totalWithdrawl := 0.0
	totalDeposit := 0.0

	ctx.GetLogger().Debugf("%s\n", transactions)

	for _, tx := range transactions {
		if tx.GetType() == "deposit" {
			totalDeposit += tx.GetAmount()
		} else if tx.GetType() == "withdrawl" {
			totalWithdrawl += tx.GetAmount()
		}
		assert.Equal(t, true, tx.GetDate().Before(time.Now()))
	}
	assert.Equal(t, true, totalDeposit > 0)
	assert.Equal(t, true, totalWithdrawl > 0)

	test.CleanupIntegrationTest()
}

func TestBlockCypher_GetTokenTransactions(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service := NewBlockCypherService(ctx, NewMarketCapService(ctx))
	transactions, err := service.GetTokenTransactions(os.Getenv("ETH_ADDRESS"), os.Getenv("TOKEN_ADDRESS"))

	assert.Nil(t, err)
	assert.NotNil(t, transactions)

	balance := 0.0
	for _, tx := range transactions {
		//		ctx.GetLogger().Debugf("%f\n", tx.GetAmount())
		balance += tx.GetAmount()
	}

	ctx.GetLogger().Debugf("%f\n", balance)

	//ctx.GetLogger().Debugf("%s\n", transactions)

	test.CleanupIntegrationTest()
}

func TestBlockCypher_GetContract(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	service := NewBlockCypherService(ctx, NewMarketCapService(ctx))

	contract, err := service.GetContract(os.Getenv("TOKEN_ADDRESS"))
	assert.Nil(t, err)
	util.DUMP(contract)

	test.CleanupIntegrationTest()
}
*/
