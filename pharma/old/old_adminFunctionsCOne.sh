#!/bin/bash
#
#
# Exit on first error

echo What you want to do? Read all medicaments or Read all users or Read a medicament info:
read admin_action
export ADMIN_ACTION="$admin_action" 
pushd ../test-network
if [[ "$ADMIN_ACTION" == "Read all medicaments" ]]; then
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"GetAllMedicaments","Args":["'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'"]}'
elif [[ "$ADMIN_ACTION" == "Read all users" ]]; then
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"GetAllUsers","Args":["'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'"]}'
elif [[ "$ADMIN_ACTION" == "Read a medicament info" ]]; then
    echo READ PRODUCT CODE:
    read read_product_code
    export READ_PRODUCT_CODE="$read_product_code"
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"GetMedicament","Args":["'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${READ_PRODUCT_CODE}"'"]}'
else
    echo Chosen action does not exist. 
fi
# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"GetMedicament","Args":["'"${USERNAME_C1}"'", , "'"${ENTITY_ID_C1}"'", "'"${READ_PRODUCT_CODE}"'"]}'
popd