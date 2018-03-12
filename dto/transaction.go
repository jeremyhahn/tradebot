package dto

import (
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type TransactionDTO struct {
	Date   time.Time
	Type   string
	Amount float64
	Fee    float64
}

func NewTransactionDTO() common.Transaction {
	return &TransactionDTO{}
}

func (t *TransactionDTO) GetDate() time.Time {
	return t.Date
}

func (t *TransactionDTO) GetType() string {
	return t.Type
}

func (t *TransactionDTO) GetAmount() float64 {
	return t.Amount
}

func (t *TransactionDTO) GetFee() float64 {
	return t.Fee
}

func (t *TransactionDTO) String() string {
	return fmt.Sprintf("[TransactionDTO] Date: %s, Type: %s, Amount: %.8f, Fee: %.8f",
		t.Date, t.Type, t.Amount, t.Fee)
}
