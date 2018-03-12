FROM golang:1.10-stretch

MAINTAINER Jeremy Hahn version: 0.2

RUN apt-get update -y && apt-get install -y build-essential

RUN git clone https://github.com/ethereum/go-ethereum
RUN cd go-ethereum && make geth

ENV PATH="${PATH}:/go/go-ethereum/build/bin"
ENV TRADEHOME="${GOPATH}/src/github.com/jeremyhahn/tradebot"

RUN mkdir -p "${TRADEHOME}/db"
COPY ./ "${TRADEHOME}"
RUN cd "${TRADEHOME}" && make dockerbuild

#RUN curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
#RUN cd "${TRADEHOME}/webui" && npm i && yarn build:dev

CMD "${TRADEHOME}/test/docker/start.sh"

EXPOSE 8080
EXPOSE 30303
