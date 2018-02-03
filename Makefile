ORG := jeremyhahn
PACKAGE := tradebot
TARGET_OS := linux

ARCH := $(shell go env GOARCH)
OS := $(shell go env GOOS)
BUILD_DIR = build

export PATH := $(GOPATH)/bin:$(PATH)

default: all

deps:
	go get "github.com/op/go-logging"
	go get "github.com/jinzhu/gorm"
	go get "github.com/jinzhu/gorm/dialects/sqlite"
	go get "github.com/op/go-logging"
	go get "github.com/gorilla/websocket"
	go get "github.com/preichenberger/go-gdax"
	go get "github.com/toorop/go-bittrex"
	go get "github.com/adshao/go-binance"
	go get "github.com/stretchr/testify"

test:
	cd common && go test common/*.go
	cd dao && go test dao/*.go
	cd exchange && go test exchange/*.go
	cd indicators && go test indicators/*.go
	cd restapi && go test restapi/*.go
	cd service && go test service/*.go
	cd strategy && go test strategy/*.go
	cd webservice && go test webservice/*.go

build:
	cd plugins/indicators && go build -buildmode=plugin ExampleIndicator.go
	go build
