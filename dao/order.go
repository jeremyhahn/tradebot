package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type OrderDAO interface {
	Create(Order entity.OrderEntity) error
	Save(Order entity.OrderEntity) error
	Find() ([]entity.Order, error)
	GetByTrade(trade entity.TradeEntity) (entity.OrderEntity, error)
}

type OrderDAOImpl struct {
	ctx   *common.Context
	Items []entity.Order
	OrderDAO
}

func NewOrderDAO(ctx *common.Context) OrderDAO {
	ctx.CoreDB.AutoMigrate(&entity.Order{})
	return &OrderDAOImpl{ctx: ctx}
}

func (dao *OrderDAOImpl) Create(Order entity.OrderEntity) error {
	return dao.ctx.CoreDB.Create(Order).Error
}

func (dao *OrderDAOImpl) Save(Order entity.OrderEntity) error {
	return dao.ctx.CoreDB.Save(Order).Error
}

func (dao *OrderDAOImpl) Find() ([]entity.Order, error) {
	var Orders []entity.Order
	daoUser := &entity.User{Id: dao.ctx.User.GetId()}
	if err := dao.ctx.CoreDB.Model(daoUser).Related(&Orders).Error; err != nil {
		return nil, err
	}
	return Orders, nil
}

/*
func BulkInsert(unsavedRows []entity.OrderEntity) error {
	valueStrings := make([]string, 0, len(unsavedRows))
	valueArgs := make([]interface{}, 0, len(unsavedRows)*3)
	for _, post := range unsavedRows {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, post.Column1)
		valueArgs = append(valueArgs, post.Column2)
		valueArgs = append(valueArgs, post.Column3)
	}
	stmt := fmt.Sprintf("INSERT INTO my_sample_table (column1, column2, column3) VALUES %s", strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}
*/
