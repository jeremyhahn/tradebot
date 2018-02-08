FROM golang:1.9.4

MAINTAINER Jeremy Hahn version: 0.1

RUN git clone https://c2d44921eeb69d65f8eb297e721f463039d0453d:x-oauth-basic@github.com/jeremyhahn/tradebot.git "${GOPATH}/src/github.com/jeremyhahn/tradebot"
RUN cd "${GOPATH}/src/github.com/jeremyhahn/tradebot" && make deps && make build

EXPOSE 8080
