//// +build integration

package dao

import (
	"os"
	"sync"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logging "github.com/op/go-logging"
)

var TEST_CONTEXT *common.Context
var TEST_LOCK sync.Mutex
var TEST_USERNAME = "test"
var TEST_DBPATH = "/tmp/tradebot-integration-testing.db"

func NewIntegrationTestContext() *common.Context {

	TEST_LOCK.Lock()

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	db, err := gorm.Open("sqlite3", TEST_DBPATH)
	db.LogMode(true)
	if err != nil {
		panic(err)
	}

	TEST_CONTEXT = &common.Context{
		DB:     db,
		Logger: logger,
		User: &common.User{
			Id:            1,
			Username:      TEST_USERNAME,
			LocalCurrency: "USD"}}

	userDAO := NewUserDAO(TEST_CONTEXT)
	userDAO.Save(&User{Username: TEST_USERNAME, LocalCurrency: "USD"})

	return TEST_CONTEXT
}

func CleanupIntegrationTest() {
	if TEST_CONTEXT != nil {
		TEST_CONTEXT.DB.Close()
		TEST_LOCK.Unlock()
		os.Remove(TEST_DBPATH)
	}
}
