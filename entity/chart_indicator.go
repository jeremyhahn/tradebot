package entity

type ChartIndicator struct {
	Id         uint   `gorm:"primary_key"`
	ChartId    uint   `gorm:"foreign_key;unique_index:idx_chart_indicator"`
	Name       string `gorm:"unique_index:idx_chart_indicator"`
	Parameters string `gorm:"not null"`
}

func (entity *ChartIndicator) GetId() uint {
	return entity.Id
}

func (entity *ChartIndicator) GetChartId() uint {
	return entity.ChartId
}

func (entity *ChartIndicator) GetName() string {
	return entity.Name
}

func (entity *ChartIndicator) GetParameters() string {
	return entity.Parameters
}
