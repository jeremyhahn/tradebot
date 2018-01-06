package service

/*
func TestPortfolioService_Build(t *testing.T) {
	ctx := test.NewTestContext()
	service := NewPortfolioService(ctx)
	portfolio := service.Build()

	if len(portfolio.Exchanges) <= 0 {
		t.Fatal("[TestPortfolioService_Build] Unable to get exchanges")
	}
	test.CleanupMockContext()
}

func TestPortfolioService_Stream(t *testing.T) {
	ctx := test.NewTestContext()
	service := NewPortfolioService(ctx)

	channel := service.Stream(ctx.User)
	portfolio := <-channel
	portfolio2 := <-channel
	portfolio3 := <-channel
	service.Stop()

	if len(portfolio.Exchanges) <= 0 || len(portfolio2.Exchanges) <= 0 || len(portfolio3.Exchanges) <= 0 {
		t.Fatal("[TestPortfolioService_Stream] Unable to get stream")
	}

	test.CleanupMockContext()
}
*/
