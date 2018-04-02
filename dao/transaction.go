package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type TransactionDAO interface {
	Create(Transaction entity.TransactionEntity) error
	Save(Transaction entity.TransactionEntity) error
	Find() ([]entity.Transaction, error)
	GetByTrade(trade entity.TradeEntity) (entity.TransactionEntity, error)
}

type TransactionDAOImpl struct {
	ctx common.Context
	TransactionDAO
}

func NewTransactionDAO(ctx common.Context) TransactionDAO {
	return &TransactionDAOImpl{ctx: ctx}
}

func (dao *TransactionDAOImpl) Create(transaction entity.TransactionEntity) error {
	return dao.ctx.GetCoreDB().Create(transaction).Error
}

func (dao *TransactionDAOImpl) Save(Transaction entity.TransactionEntity) error {
	return dao.ctx.GetCoreDB().Save(Transaction).Error
}

func (dao *TransactionDAOImpl) Find() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	daoUser := &entity.User{Id: dao.ctx.GetUser().GetId()}
	if err := dao.ctx.GetCoreDB().Model(daoUser).Related(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
