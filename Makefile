ORG := automatethethingsllc
PACKAGE := tradebot
TARGET_OS := linux

ARCH := $(shell go env GOARCH)
OS := $(shell go env GOOS)
BUILD_DIR = build

export PATH := $(GOPATH)/bin:$(PATH)

default: build

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

unittest:
	go test -v dao/* -tags !integration

indicators:
	cd plugins/indicators/src && go build -buildmode=plugin example.go && mv example.so ../
	cd plugins/indicators/src && go build -buildmode=plugin sma.go && mv sma.so ../
	cd plugins/indicators/src && go build -buildmode=plugin ema.go && mv ema.so ../
	cd plugins/indicators/src && go build -buildmode=plugin rsi.go sma.go && mv rsi.so ../
	cd plugins/indicators/src && go build -buildmode=plugin bollinger_bands.go sma.go && mv bollinger_bands.so ../
	cd plugins/indicators/src && go build -buildmode=plugin macd.go ema.go && mv macd.so ../
	cd plugins/indicators/src && go build -buildmode=plugin obv.go && mv obv.so ../

strategies:
	cd plugins/strategies/src && go build -buildmode=plugin default.go && mv default.so ../

clean:
	cd plugins/indicators && rm -rf *.so
	cd plugins/strategies && rm -rf *.so
	rm -rf tradebot

build: clean test indicators strategies
	go build
