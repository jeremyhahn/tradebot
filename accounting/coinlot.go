package accounting

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Coinlot struct {
	Date      time.Time
	Currency  string
	Quantity  decimal.Decimal
	UnitPrice decimal.Decimal
	SalePrice decimal.Decimal
	CostBasis decimal.Decimal
}

func (coinlot *Coinlot) String() string {
	return fmt.Sprintf("[Coinlot] Date: %s, Currency: %s, Quantity: %s, UnitPrice: %s, SalePrice: %s, CostBasis: %s",
		coinlot.Date, coinlot.Currency, coinlot.Quantity, coinlot.UnitPrice, coinlot.SalePrice, coinlot.CostBasis)
}
