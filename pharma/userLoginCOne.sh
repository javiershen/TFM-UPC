#!/bin/bash
#
#
# Exit on first error


echo ENTITY_ID_C1 pharmacy2 lab1:
read entity
export ENTITY_ID_C1="$entity"
echo USERNAME_C1 userLab adminLab userPharmacy adminPharmacy:
read username
export USERNAME_C1="$username"
echo USER_PASSWORD_C1 adminpw:
read userpassword
export USER_PASSWORD_C1="$userpassword"
echo You are logged with user "${USERNAME_C1}"
pushd ../test-network
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"LogIn","Args":["'"${ENTITY_ID_C1}"'","'"${USERNAME_C1}"'","'"${USER_PASSWORD_C1}"'"]}'
popd