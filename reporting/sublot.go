package reporting

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Sublot struct {
	Date      time.Time
	Quantity  decimal.Decimal
	Price     decimal.Decimal
	CostBasis decimal.Decimal
}

func (remainder *Sublot) String() string {
	return fmt.Sprintf("[Sublot] Date: %s, Quantity: %s, Price: %s, CostBasis: %s",
		remainder.Date, remainder.Quantity, remainder.Price, remainder.CostBasis)
}
