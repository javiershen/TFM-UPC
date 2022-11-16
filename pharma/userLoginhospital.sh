#!/bin/bash
#
#
# Exit on first error


echo ENTITY_ID pharmacy2 hospital1:
read entity
export ENTITY_ID="$entity"
echo USERNAME sanitaryUser adminLab pharmacyUser pharmacyAdmin:
read username
export USERNAME_CHANNELTWO="$username"
echo USER_PASSWORD adminpw:
read userpassword
export USER_PASSWORD="$userpassword"
echo You are logged with user "${USERNAME_CHANNELTWO}"
pushd ../test-network
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"LogIn","Args":["'"${ENTITY_ID}"'","'"${USERNAME_CHANNELTWO}"'","'"${USER_PASSWORD}"'"]}'
popd