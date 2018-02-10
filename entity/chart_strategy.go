package entity

type ChartStrategy struct {
	Id         uint   `gorm:"primary_key"`
	ChartId    uint   `gorm:"foreign_key;unique_index:idx_chart_strategy"`
	Name       string `gorm:"unique_index:idx_chart_strategy"`
	Parameters string `gorm:"not null"`
}

func (entity *ChartStrategy) GetId() uint {
	return entity.Id
}

func (entity *ChartStrategy) GetChartId() uint {
	return entity.ChartId
}

func (entity *ChartStrategy) GetName() string {
	return entity.Name
}

func (entity *ChartStrategy) GetParameters() string {
	return entity.Parameters
}
