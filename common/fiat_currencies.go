package common

var FiatCurrencies = map[string]*Currency{
	"USD": &Currency{
		ID:           "USD",
		Name:         "United States dollar",
		Symbol:       "$",
		BaseUnit:     100,
		DecimalPlace: 2},
	"CAD": &Currency{
		ID:           "CAD",
		Name:         "Canadian dollar",
		Symbol:       "$",
		BaseUnit:     100,
		DecimalPlace: 2},
	"EUR": &Currency{
		ID:           "EUR",
		Name:         "Euro",
		Symbol:       "€",
		BaseUnit:     100,
		DecimalPlace: 2},
	"JPY": &Currency{
		ID:           "JPY",
		Name:         "Japanese yen",
		Symbol:       "¥",
		BaseUnit:     100,
		DecimalPlace: 2},
	"GBP": &Currency{
		ID:           "GBP",
		Symbol:       "£",
		BaseUnit:     100,
		DecimalPlace: 2},
	"CHF": &Currency{
		ID:           "CHF",
		Name:         "Swiss franc",
		Symbol:       "Fr",
		BaseUnit:     100,
		DecimalPlace: 2},
	"KRW": &Currency{
		ID:           "KRW",
		Name:         "South Korean won",
		Symbol:       "₩",
		BaseUnit:     100,
		DecimalPlace: 2}}
