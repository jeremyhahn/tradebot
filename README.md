# tradebot

Ethereum powered automated cryptocurrency trading platform.

## Features

* Plugin system supports financial indicators, trading strategies and exchanges
* Install / Uninstall plugins without restarting the trading platform
* Configurable financial indicators
* Configurable trading strategies
* Trading bot automatically execute trades based on configured trading strategies / indicators
* Marketplace for custom financial indicators and trading strategies

## Build

Due to limitations of [Golang plugins](https://golang.org/pkg/plugin/), this software requires a Linux operating. [Debian](https://www.debian.org/) is a great choice.

#### Dependencies

1. [Golang](https://golang.org/)
2. [Make](https://www.gnu.org/software/make/)

#### Optional Dependencies

1. [Docker](https://www.docker.com/)
