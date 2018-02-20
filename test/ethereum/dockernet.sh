#!/bin/bash

rm -rf blockchain/
geth --datadir ./blockchain init genesis.json
geth --datadir ./blockchain --networkid 420 &
