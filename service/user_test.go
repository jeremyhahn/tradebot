package service

/*
func TestUserService_CreateUser(t *testing.T) {
	service := createTestUserService()

	userById := service.GetUserById(1)
	if userById.Id != 1 {
		t.Errorf("[TestUserService_CreateUser] Unexpected user id: %d, expected: %d", userById.Id, 1)
	}
	if userById.Username != "test" {
		t.Errorf("[TestUserService_CreateUser] Unexpected username: %s, expected: %s", userById.Username, "test")
	}

	userByName := service.GetUserByName("test")
	if userByName.Id != 1 {
		t.Errorf("[TestUserService_CreateUser] Unexpected user id: %d, expected: %d", userByName.Id, 1)
	}
	if userByName.Username != "test" {
		t.Errorf("[TestUserService_CreateUser] Unexpected username: %s, expected: %s", userByName.Username, "test")
	}

	test.CleanupMockContext()
}

func TestUserService_GetWallets(t *testing.T) {
	service := createTestUserService()

	user := service.GetCurrenctUser()
	wallets := service.GetWallets(user)

	if wallets[0].Address != test.BTC_ADDRESS {
		t.Fatal("[TestUserService_GetWallets] Unexpected BTC wallet address")
	}

	if wallets[1].Address != test.XRP_ADDRESS {
		t.Fatal("[TestUserService_GetWallets] Unexpected XRP wallet address")
	}

	test.CleanupMockContext()
}

func TestUserService_GetWallet(t *testing.T) {
	service := createTestUserService()

	user := service.GetCurrenctUser()
	wallet := service.GetWallet(user, "BTC")

	if wallet.Address != test.BTC_ADDRESS {
		t.Fatal("[TestUserService_GetWallet] Unexpected BTC wallet address")
	}

	test.CleanupMockContext()
}

func TestUserService_GetExchanges(t *testing.T) {
	service := createTestUserService()

	user := service.GetCurrenctUser()
	exchanges := service.GetExchanges(user)

	service.ctx.Logger.Debugf("[Fuck] %+v\n", exchanges)

	if exchanges[0].Name != "GDAX" {
		t.Fatalf("[TestUserService_GetExchanges] Failed to get expected GDAX exchange, got: %s", exchanges[0].Name)
	}
	if exchanges[0].Total <= 0 {
		t.Fatalf("[TestUserService_GetExchanges] Failed to get GDAX total, got: %f", exchanges[0].Total)
	}

	if exchanges[1].Name != "bittrex" {
		t.Fatalf("[TestUserService_GetExchanges] Failed to get expected Bittrex exchange, got: %s", exchanges[1].Name)
	}
	if exchanges[1].Total <= 0 {
		t.Fatalf("[TestUserService_GetExchanges] Failed to get Bittrex total, got: %f", exchanges[1].Total)
	}

	if exchanges[2].Name != "binance" {
		t.Fatalf("[TestUserService_GetExchanges] Failed to get expected Binance exchange, got: %s", exchanges[2].Name)
	}
	if exchanges[2].Total <= 0 {
		t.Fatalf("[TestUserService_GetExchanges] Failed to get Binance total, got: %f", exchanges[2].Total)
	}

	test.CleanupMockContext()
}

func createTestUserService() *UserService {
	ctx := test.NewTestContext()
	dao := dao.NewUserDAO(ctx)
	return NewUserService(ctx, dao)
}
*/
