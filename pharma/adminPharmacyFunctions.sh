#!/bin/bash
# Bash Menu Script Example
pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null
}

pushd ../test-network

PS3='Which of the following actions you want to do: '
options=("1- Get pharmacy Stock" "2- Read all users" "3- Quit")
select opt in "${options[@]}"
do
    case $opt in
        "1- Get pharmacy Stock")     
            
           peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"GetPharmacyStock","Args":["'"${USERNAME_C2}"'", "'"${ENTITY_ID_C2}"'"]}'
            break
            ;;
        "2- Read all users")
           
           peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"GetAllUsers","Args":["'"${USERNAME_C2}"'", "'"${ENTITY_ID_C2}"'"]}'
            break
            ;;
        "3- Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
popd
