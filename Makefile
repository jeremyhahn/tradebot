ORG := automatethethingsllc
PACKAGE := tradebot
TARGET_OS := linux

ARCH := $(shell go env GOARCH)
OS := $(shell go env GOOS)
BUILD_DIR = build

export PATH := $(GOPATH)/bin:$(PATH)

default: build

certgen:
	mkdir -p ssl
	openssl req -newkey rsa:2048 -nodes -keyout ssl/key.pem -x509 -days 365 -out ssl/cert.pem

deps:
	go get "github.com/stretchr/testify"
	go get "github.com/op/go-logging"
	go get "github.com/jinzhu/gorm"
	go get "github.com/jinzhu/gorm/dialects/sqlite"
#	go get "github.com/jinzhu/gorm/dialects/mysql"
	go get "github.com/gorilla/websocket"
	go get "github.com/preichenberger/go-gdax"
	go get "github.com/toorop/go-bittrex"
	go get "github.com/adshao/go-binance"

unittest:
	cd dao && go test -v
	cd plugins/indicators/src && go test -v
	cd service && go test -v

integrationtest:
	cd dao && go test -v -tags integration
	cd plugins/indicators/src && go test -v -tags integration
	cd service && go test -v -tags integration

test: unittest integrationtest

indicators:
	cd plugins/indicators/src && go build -buildmode=plugin -o ../example.so example.go
	cd plugins/indicators/src && go build -buildmode=plugin -o ../sma.so sma.go
	cd plugins/indicators/src && go build -buildmode=plugin -o ../ema.so ema.go
	cd plugins/indicators/src && go build -buildmode=plugin -o ../rsi.so rsi.go sma.go
	cd plugins/indicators/src && go build -buildmode=plugin -o ../bollinger_bands.so bollinger_bands.go sma.go
	cd plugins/indicators/src && go build -buildmode=plugin -o ../macd.so macd.go ema.go
	cd plugins/indicators/src && go build -buildmode=plugin -o ../obv.so obv.go

strategies:
	cd plugins/strategies/src && go build -buildmode=plugin -o ../default.so default.go

plugins: indicators strategies

clean:
	rm -rf ssl/
	cd plugins/indicators && rm -rf *.so
	cd plugins/strategies && rm -rf *.so
	rm -rf tradebot

builddebug: clean plugins
	go build -gcflags "-N -l"

quickdebug:
	go build -gcflags "-N -l"

quickbuild:
	go build -ldflags "-w"

build: clean plugins quickbuild test
