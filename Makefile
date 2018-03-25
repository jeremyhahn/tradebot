ORG := automatethethingsllc
PACKAGE := tradebot
TARGET_OS := linux

ARCH := $(shell go env GOARCH)
OS := $(shell go env GOOS)
BUILD_DIR = build

export PATH := $(GOPATH)/bin:$(PATH)

default: build

deps:
	go get "github.com/stretchr/testify"
	go get "github.com/op/go-logging"
	go get "github.com/jinzhu/gorm"
	go get "github.com/jinzhu/gorm/dialects/sqlite"
	#go get "github.com/jinzhu/gorm/dialects/mysql"
	go get "golang.org/x/text/encoding/unicode"
	go get "github.com/patrickmn/go-cache"
	go get "github.com/gorilla/websocket"
	go get "github.com/gorilla/mux"
	go get "github.com/codegangsta/negroni"
	go get "github.com/dgrijalva/jwt-go"
	go get "github.com/ethereum/go-ethereum"
	go get "github.com/Zauberstuhl/go-coinbase"
	go get "github.com/preichenberger/go-gdax"
	go get "github.com/toorop/go-bittrex"
	go get "github.com/adshao/go-binance"
	go get "github.com/joho/godotenv"
	go get "github.com/shopspring/decimal"
	go get "golang.org/x/crypto/bcrypt"

certs:
	mkdir -p keys/
	openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 -keyout keys/key.pem -out keys/cert.pem \
          -subj "/C=US/ST=Blockchain/L=Tradebot/O=Cryptoconomy/CN=localhost"
	openssl genrsa -out keys/rsa.key 1024
	openssl rsa -in keys/rsa.key -pubout > keys/rsa.pub

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
	cd plugins/exchanges/src && go build -buildmode=plugin -o ../coinbase.so coinbase.go
	cd plugins/exchanges/src && go build -buildmode=plugin -o ../gdax.so gdax.go
	cd plugins/exchanges/src && go build -buildmode=plugin -o ../bittrex.so bittrex.go
	cd plugins/exchanges/src && go build -buildmode=plugin -o ../binance.so binance.go

strategies:
	cd plugins/strategies/src && go build -buildmode=plugin -o ../default.so default.go

plugins: indicators strategies

clean:
	rm -rf keys/
	cd plugins/indicators && rm -rf *.so
	cd plugins/strategies && rm -rf *.so
	cd plugins/exchanges && rm -rf *.so
	rm -rf tradebot
	cd webapp && yarn clean

debugbuild: clean plugins
	go build -gcflags "-N -l"

quickdebug:
	go build -gcflags "-N -l"

quickbuild:
	go build -ldflags "-w"

dockerbuild: clean deps certs plugins
	go build
	./docker-build.sh

webapp:
	cd webui && npm i && yarn build:dev

build: clean certs plugins test webui
	go build
