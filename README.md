# tradebot

Ethereum powered automated cryptocurrency trading platform.

## Current Features

* User friendly web interface
* User friendly portfolio with integrated buying / selling
* Exchange order / trade history reporting; Generate accounting reports
* Plugin system supports financial indicators, trading strategies, exchanges and wallets
* Add / remove functionality without interrupting underlying platform or current trades
* Trading bot automatically executes trades based on configured trading strategies / indicators
* Json Web Token (JWT) protected APIs
* REST and WebSocket APIs

## Roadmap
* Develop initial set of smart contracts to enable trading
* Live Charts
* Marketplace for custom plugins
* Distributed trading protocol
* Exchange features: distributed order book, buy / sell between network users, Atomic Swaps
* Financial lending

## Build

Due to limitations of [Golang plugins](https://golang.org/pkg/plugin/), this software requires a Linux or Mac operating system to run natively. [Debian](https://www.debian.org/) is a great choice. [Docker](https://www.docker.com/) support is included for Windows users.

#### Dependencies

1. [Golang](https://golang.org/)
2. [Make](https://www.gnu.org/software/make/)
3. [OpenSSL](https://www.openssl.org/)
4. [Docker](https://www.docker.com/) (Optional)

#### Windows

1. [Docker Toolkit](https://docs.docker.com/toolbox/toolbox_install_windows/) (Windows Home Edition or < Windows 10)
2. [Docker Community Edition](https://store.docker.com/editions/community/docker-ce-desktop-windows) (Windows 10 Professional or Enterprise)

Windows Home Edition and versions prior to Windows 10 do not support native virtualization. As such, they require `docker-machine` included in the Docker Toolkit. Newer versions of Windows that support native virtualization can take advantage of the latest Docker CE with enhanced performance.

## Tech Stack

#### User Interface
* [React.js](https://reactjs.org/)
* [Material-UI](https://material-ui-next.com/)

#### Backend
* [Golang](https://golang.org/)
* [Ethereum](https://www.ethereum.org/)
* [GORM](http://gorm.io/)
* [SQLite](https://www.sqlite.org/)
