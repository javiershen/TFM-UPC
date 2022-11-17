#!/bin/bash
#
#
# Exit on first error


echo MEDICAMENT NAME:
read medicament_name
export MEDICAMENT_NAME="$medicament_name"
echo MEDICAMENT SERIAL NUMBER:
read medicament_serial_number
export MEDICAMENT_SERIAL_NUMBER="$medicament_serial_number"
echo PRODUCT CODE:
read product_code
export PRODUCT_CODE="$product_code"
echo MEDICAMENT LOT NUMBER:
read medicament_lot_number
export MEDICAMENT_LOT_NUMBER="$medicament_lot_number"
echo EXPIRATION YEAR:
read expiration_year
export EXPIRATION_YEAR="$expiration_year"
echo EXPIRATION MONTH:
read expiration_month
export EXPIRATION_MONTH="$expiration_month"
pushd ../test-network
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["RegisterMedicament","'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${MEDICAMENT_NAME}"'", "'"${MEDICAMENT_LOT_NUMBER}"'", "'"${PRODUCT_CODE}"'", "'"${MEDICAMENT_SERIAL_NUMBER}"'", "'"${EXPIRATION_YEAR}"'", "'"${EXPIRATION_MONTH}"'"]}'
popd