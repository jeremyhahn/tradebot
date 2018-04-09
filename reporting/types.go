package reporting

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type Remainder struct {
	Quantity  decimal.Decimal
	CostBasis decimal.Decimal
	Tx        common.Transaction
}

func (remainder *Remainder) String() string {
	return fmt.Sprintf("[Remainder] Amount: %s, CostBasis: %s, Tx: %s",
		remainder.Quantity, remainder.CostBasis, remainder.Tx)
}

type NeedsMore struct {
	SaleQuantity  decimal.Decimal
	SalesTx       common.Transaction
	Needed        decimal.Decimal
	PurchasePrice decimal.Decimal
}

func (needsMore *NeedsMore) String() string {
	return fmt.Sprintf("[NeedsMore] SaleQuantity: %s, Needed: %s, PurchasePrice: %s, Tx: %s",
		needsMore.SaleQuantity, needsMore.Needed, needsMore.PurchasePrice, needsMore.SalesTx)
}
