#!/bin/bash
#
#
# Exit on first error


echo PRODUCT CODE:
read product_code
export PRODUCT_CODE="$product_code"
echo Medicament destination:
read medicament_destination
export MEDICAMENT_DESTINATION="$medicament_destination"

pushd ../test-network
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["DispatchMedicament", "'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${MEDICAMENT_DESTINATION}"'", "'"${PRODUCT_CODE}"'"]}'
popd