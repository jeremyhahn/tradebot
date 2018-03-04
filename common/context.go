package common

import (
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
)

type Context interface {
	GetLogger() *logging.Logger
	GetCoreDB() *gorm.DB
	GetPriceDB() *gorm.DB
	GetUser() UserContext
	SetUser(user UserContext)
	GetDebug() bool
	GetSSL() bool
}

type Ctx struct {
	Logger  *logging.Logger
	CoreDB  *gorm.DB
	PriceDB *gorm.DB
	User    UserContext
	Debug   bool
	SSL     bool
	Context
}

func (c *Ctx) GetLogger() *logging.Logger {
	return c.Logger
}

func (c *Ctx) GetCoreDB() *gorm.DB {
	return c.CoreDB
}

func (c *Ctx) GetPriceDB() *gorm.DB {
	return c.PriceDB
}

func (c *Ctx) GetUser() UserContext {
	return c.User
}

func (c *Ctx) SetUser(user UserContext) {
	c.User = user
}

func (c *Ctx) GetDebug() bool {
	return c.Debug
}

func (c *Ctx) GetSSL() bool {
	return c.SSL
}
