#!/bin/bash

rm -rf blockchain/
geth --datadir ./blockchain init genesis.json --port 30304
geth --datadir ./blockchain --networkid 420 console --port 30304

# Create new account
# personal.newAccount("test")

# Check default account
# eth.coinbase

# Set etherbase
# miner.setEtherbase(eth.accounts[0])

# Check balance
# eth.getBalance(eth.coinbase)

# Start mining
# miner.start()

# Stop mining
# miner.stop()
