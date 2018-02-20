#!/bin/bash

cd "${TRADEHOME}/test/ethereum/" && ./dockernet.sh

sleep 5

cd $TRADEHOME
./tradebot --initdb
./tradebot --debug
