#!/bin/bash
#
#
# Exit on first error


echo ENTITY ID - pharmacy2 hospital1:
read entity
export ENTITY_ID_C2="$entity"
echo USER NAME - userSanitary adminSanitary userPharmacy adminPharmacy:
read username
export USERNAME_C2="$username"
echo PASSWORD:
read userpassword
export USER_PASSWORD_C2="$userpassword"
echo You are logged with user "${USERNAME_C2}"
pushd ../test-network
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"LogIn","Args":["'"${ENTITY_ID_C2}"'","'"${USERNAME_C2}"'","'"${USER_PASSWORD_C2}"'"]}'
popd
