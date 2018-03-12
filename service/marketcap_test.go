// +build integration

package service

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarketCapService_GetMarkets(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)
	assert.Equal(t, true, len(marketcap.GetMarkets()) > 0)
	CleanupIntegrationTest()
}

func TestMarketCapService_GetGlobalMarkets(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)
	assert.Equal(t, true, marketcap.GetGlobalMarket("USD").LastUpdated > 0)
	CleanupIntegrationTest()
}

func TestMarketCapService_GetMarketsByPrice(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)

	asc := marketcap.GetMarketsByPrice("asc")
	desc := marketcap.GetMarketsByPrice("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/GetMarketsByPrice", ascData, 0644)
	ioutil.WriteFile("/tmpGetMarketsByPrice", descData, 0644)

	priceI, _ := strconv.ParseFloat(asc[0].PriceUSD, 64)
	priceJ, _ := strconv.ParseFloat(desc[0].PriceUSD, 64)

	assert.Equal(t, false, priceI > priceJ)

	CleanupIntegrationTest()
}

func TestMarketCapService_GetMarketsByPercentChange1H(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)

	asc := marketcap.GetMarketsByPercentChange1H("asc")
	desc := marketcap.GetMarketsByPercentChange1H("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/GetMarketsByPercentChange1H", ascData, 0644)
	ioutil.WriteFile("/tmp/GetMarketsByPercentChange1H", descData, 0644)

	fi, _ := strconv.ParseFloat(asc[0].PercentChange1h, 64)
	fj, _ := strconv.ParseFloat(desc[0].PercentChange1h, 64)

	assert.Equal(t, false, fi > fj)

	CleanupIntegrationTest()
}

func TestMarketCapService_GetMarketsByPercentChange24H(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)

	asc := marketcap.GetMarketsByPercentChange24H("asc")
	desc := marketcap.GetMarketsByPercentChange24H("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/GetMarketsByPercentChange24H", ascData, 0644)
	ioutil.WriteFile("/tmp/GetMarketsByPercentChange24H", descData, 0644)

	fi, _ := strconv.ParseFloat(asc[0].PercentChange24h, 64)
	fj, _ := strconv.ParseFloat(desc[0].PercentChange24h, 64)

	assert.Equal(t, false, fi > fj)

	CleanupIntegrationTest()
}

func TestMarketCapService_GetMarketsByPercentChange7D(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)

	asc := marketcap.GetMarketsByPercentChange7D("asc")
	desc := marketcap.GetMarketsByPercentChange7D("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/GetMarketsByPercentChange7D", ascData, 0644)
	ioutil.WriteFile("/tmp/GetMarketsByPercentChange7D", descData, 0644)

	fi, _ := strconv.ParseFloat(asc[0].PercentChange7d, 64)
	fj, _ := strconv.ParseFloat(desc[0].PercentChange7d, 64)

	assert.Equal(t, false, fi > fj)

	CleanupIntegrationTest()
}

func TestMarketCapService_GetMarketsByTopPerformers(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)

	asc := marketcap.GetMarketsByTopPerformers("asc")
	desc := marketcap.GetMarketsByTopPerformers("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/GetMarketsByTopPerformers", ascData, 0644)
	ioutil.WriteFile("/tmp/GetMarketsByTopPerformers", descData, 0644)

	CleanupIntegrationTest()
}

func TestMarketCapService_GetMarketsByTrending(t *testing.T) {
	ctx := NewIntegrationTestContext()
	marketcap := NewMarketCapService(ctx)

	asc := marketcap.GetTrendingMarkets("asc")
	desc := marketcap.GetTrendingMarkets("desc")

	ascData, _ := json.MarshalIndent(asc, "", "    ")
	descData, _ := json.MarshalIndent(desc, "", "    ")

	ioutil.WriteFile("/tmp/GetMarketsByTrending", ascData, 0644)
	ioutil.WriteFile("/tmp/GetMarketsByTrending", descData, 0644)

	CleanupIntegrationTest()
}
