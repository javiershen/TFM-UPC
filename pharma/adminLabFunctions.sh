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
options=("1- Read all medicaments" "2- Read all users" "3- Read a medicament info" "4- Quit")
select opt in "${options[@]}"
do
    case $opt in
        "1- Read all medicaments")
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"GetAllMedicaments","Args":["'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'"]}'
            break
            ;;
        "2- Read all users")
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"GetAllUsers","Args":["'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'"]}'
            break
            ;;
        "3- Read a medicament info")
            echo READ PRODUCT CODE:
            read read_product_code
            export READ_PRODUCT_CODE="$read_product_code"
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"GetMedicament","Args":["'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${READ_PRODUCT_CODE}"'"]}'
            break
            ;;
        "4- Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
popd

