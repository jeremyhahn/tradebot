package viewmodel

type UserCryptoExchange struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Key   string `json:"key"`
	Extra string `json:"extra"`
}

func (entity *UserCryptoExchange) GetId() string {
	return entity.Id
}

func (entity *UserCryptoExchange) GetName() string {
	return entity.Name
}

func (entity *UserCryptoExchange) GetURL() string {
	return entity.URL
}

func (entity *UserCryptoExchange) GetKey() string {
	return entity.Key
}

func (entity *UserCryptoExchange) GetExtra() string {
	return entity.Extra
}
