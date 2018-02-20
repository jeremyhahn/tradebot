FROM golang:1.9.4

MAINTAINER Jeremy Hahn version: 0.1

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install -y build-essential syslog-ng

#RUN git clone https://github.com/ethereum/go-ethereum
#RUN cd go-ethereum && make geth

#ENV PATH="${PATH}:/go/go-ethereum/build/bin"

RUN git clone https://c2d44921eeb69d65f8eb297e721f463039d0453d:x-oauth-basic@github.com/jeremyhahn/tradebot.git "${GOPATH}/src/github.com/jeremyhahn/tradebot"
RUN cd "${GOPATH}/src/github.com/jeremyhahn/tradebot" && make deps && make

RUN "${GOPATH}/src/github.com/jeremyhahn/tradebot/tradebot --initdb"
RUN "${GOPATH}/src/github.com/jeremyhahn/tradebot/test/ethereum/privatenet.sh"
RUN "${GOPATH}/src/github.com/jeremyhahn/tradebot/tradebot -debug"

EXPOSE 8080
