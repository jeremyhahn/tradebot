#!/bin/bash

# docker run -p 8000:8080 tradebot

docker run -d \
  --name tradebot \
  --mount source=tradebot,target=/go/src/github.com/jeremyhahn/tradebot \
  -p 8000:8080 \
  tradebot:latest
