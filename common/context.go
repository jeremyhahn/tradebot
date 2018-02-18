package common

import (
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
)

type Context struct {
	Logger  *logging.Logger
	CoreDB  *gorm.DB
	PriceDB *gorm.DB
	User    User
	Debug   bool
	SSL     bool
}

func (c *Context) GetUser() User {
	return c.User
}

func (c *Context) SetUser(user User) {
	c.User = user
}
