#!/bin/bash
FLAGS='--networkid=8888 --preload=identities.js'
DEV_FLAGS='--nodiscover --verbosity=4'

BASE_PORT=30303
BASE_RPC_PORT=8545

/Users/hungngo/Documents/TSF/go-ethereum/build/bin/geth --testnet --datadir="/Users/hungngo/LocalDocs/ethereum-data/ropsten" --port=30303 --rpcport=8545 --rpc --rpcapi=eth,net,web3,personal,clique,admin,debug,txpool --rpcaddr=0.0.0.0 --syncmode "full" --v5disc --bootnodes=enode://94c15d1b9e2fe7ce56e458b9a3b672ef11894ddedd0c6f247e0f1d3487f52b66208fb4aeb8179fce6e3a749ea93ed147c37976d67af557508d199d9594c35f09@192.81.208.223:30303,enode://6332792c4a00e3e4ee0926ed89e0d27ef985424d97b6a45bf0f23e51f0dcb5e66b875777506458aea7af6f9e4ffb69f43f3778ee73c81ed9d34c51c4b16b0b0f@52.232.243.152:30303