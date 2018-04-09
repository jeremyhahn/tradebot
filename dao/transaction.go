package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type TransactionDAO interface {
	Create(Transaction entity.TransactionEntity) error
	Save(Transaction entity.TransactionEntity) error
	Update(tx entity.TransactionEntity, field, value string) error
	Get(id string) (entity.TransactionEntity, error)
	Find() ([]entity.Transaction, error)
}

type TransactionDAOImpl struct {
	ctx common.Context
	TransactionDAO
}

func NewTransactionDAO(ctx common.Context) TransactionDAO {
	return &TransactionDAOImpl{ctx: ctx}
}

func (dao *TransactionDAOImpl) Create(tx entity.TransactionEntity) error {
	return dao.ctx.GetCoreDB().Create(tx).Error
}

func (dao *TransactionDAOImpl) Save(tx entity.TransactionEntity) error {
	return dao.ctx.GetCoreDB().Save(tx).Error
}

func (dao *TransactionDAOImpl) Update(tx entity.TransactionEntity, field, value string) error {
	return dao.ctx.GetCoreDB().Model(tx).Update(field, value).Error
}

func (dao *TransactionDAOImpl) Get(id string) (entity.TransactionEntity, error) {
	tx := &entity.Transaction{Id: id}
	if err := dao.ctx.GetCoreDB().First(tx).Error; err != nil {
		return nil, err
	}
	return tx, nil
}

func (dao *TransactionDAOImpl) Find() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	daoUser := &entity.User{Id: dao.ctx.GetUser().GetId()}
	if err := dao.ctx.GetCoreDB().Model(daoUser).Related(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
