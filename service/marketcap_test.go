package service

/*
func TestMarketCapService_GetMarkets(t *testing.T) {
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

func TestMarketCapService_GetGlobalMarkets(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	if marketcap.GetGlobalMarket("USD").LastUpdated <= 0 {
		t.Fatal("[TestMarketCapService.GetMarkets] Unable to get global market cap")
	}

	test.CleanupMockContext()
}

func TestMarketCapService_GetMarketsByPrice(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	asc := marketcap.GetMarketsByPrice("asc")
	desc := marketcap.GetMarketsByPrice("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	//ioutil.WriteFile("/tmp/asc", ascData, 0644)
	//ioutil.WriteFile("/tmp/desc", descData, 0644)

	priceI, _ := strconv.ParseFloat(asc[0].PriceUSD, 64)
	priceJ, _ := strconv.ParseFloat(desc[0].PriceUSD, 64)

	if priceI > priceJ {
		t.Fatal("[TestMarketCapService.GetMarkets] Unable to sort market cap by price")
	}

	test.CleanupMockContext()
}

func TestMarketCapService_GetMarketsByPercentChange1H(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	asc := marketcap.GetMarketsByPercentChange1H("asc")
	desc := marketcap.GetMarketsByPercentChange1H("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/asc", ascData, 0644)
	ioutil.WriteFile("/tmp/desc", descData, 0644)

	fi, _ := strconv.ParseFloat(asc[0].PercentChange1h, 64)
	fj, _ := strconv.ParseFloat(desc[0].PercentChange1h, 64)

	if fi > fj {
		t.Fatal("[TestMarketCapService.GetMarkets] Unable to sort market cap by percent changed 1h")
	}

	test.CleanupMockContext()
}
*/

/*
func TestMarketCapService_GetMarketsByPercentChange24H(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	asc := marketcap.GetMarketsByPercentChange24H("asc")
	desc := marketcap.GetMarketsByPercentChange24H("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/asc", ascData, 0644)
	ioutil.WriteFile("/tmp/desc", descData, 0644)

	fi, _ := strconv.ParseFloat(asc[0].PercentChange24h, 64)
	fj, _ := strconv.ParseFloat(desc[0].PercentChange24h, 64)

	if fi > fj {
		t.Fatal("[TestMarketCapService.GetMarkets] Unable to sort market cap by percent changed 24h")
	}

	test.CleanupMockContext()
}

func TestMarketCapService_GetMarketsByPercentChange7D(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	asc := marketcap.GetMarketsByPercentChange7D("asc")
	desc := marketcap.GetMarketsByPercentChange7D("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/asc", ascData, 0644)
	ioutil.WriteFile("/tmp/desc", descData, 0644)

	fi, _ := strconv.ParseFloat(asc[0].PercentChange7d, 64)
	fj, _ := strconv.ParseFloat(desc[0].PercentChange7d, 64)

	if fi > fj {
		t.Fatal("[TestMarketCapService.GetMarkets] Unable to sort market cap by percent changed 7d")
	}

	test.CleanupMockContext()
}

func TestMarketCapService_GetMarketsByTopPerformers(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	asc := marketcap.GetMarketsByTopPerformers("asc")
	desc := marketcap.GetMarketsByTopPerformers("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/top-performers-asc", ascData, 0644)
	ioutil.WriteFile("/tmp/top-performers-desc", descData, 0644)

	test.CleanupMockContext()
}

func TestMarketCapService_GetMarketsByTrending(t *testing.T) {
	ctx := test.NewTestContext()
	marketcap := NewMarketCapService(ctx.Logger)

	asc := marketcap.GetMarketsByTrending("asc")
	desc := marketcap.GetMarketsByTrending("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/trending-asc", ascData, 0644)
	ioutil.WriteFile("/tmp/trending-desc", descData, 0644)

	test.CleanupMockContext()
}
*/
