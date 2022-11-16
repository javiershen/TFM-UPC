#!/bin/bash
#
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_LANGUAGE=${1:-"go"}
CC_SRC_LANGUAGE=`echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:]`

if [ "$CC_SRC_LANGUAGE" = "go" -o "$CC_SRC_LANGUAGE" = "golang" ] ; then
	CC_SRC_PATH="../asset-transfer-basic/chaincode-go/chaincode/"
else
	echo The chaincode language ${CC_SRC_LANGUAGE} is not supported by this script
	echo Supported chaincode languages are: go
	exit 1
fi

# clean out any old identites in the wallets
# rm -rf javascript/wallet/*
# rm -rf java/wallet/*
# rm -rf typescript/wallet/*
# rm -rf go/wallet/*

# launch network; create channel and join peer to channel
pushd ../test-network
./network.sh down
./network.sh up createChannel -c channel1
popd
pushd ../test-network/addOrg3
./addOrg3.sh up
popd
pushd ../test-network
./network.sh createChannelTwo -c channel2
./network.sh deployCC -ccn basic1 -ccp ../asset-transfer-basic/chaincode-ch1-go -ccl go -c channel1
./network.sh deployCCTwo -ccn basic2 -ccp ../asset-transfer-basic/chaincode-ch2-go -ccl go -c channel2
popd

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...
EOF