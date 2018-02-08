FROM bitnami/minideb:stretch

MAINTAINER Jeremy Hahn version: 0.1

RUN install_packages wget git build-essential ca-certificates locales

RUN wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.9.4.linux-amd64.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"
ENV GOPATH="/go"

RUN git clone https://c2d44921eeb69d65f8eb297e721f463039d0453d:x-oauth-basic@github.com/jeremyhahn/tradebot.git "${GOPATH}/src/github.com/jeremyhahn/tradebot"
RUN cd "${GOPATH}/src/github.com/jeremyhahn/tradebot" && make deps && make build

EXPOSE 8080
