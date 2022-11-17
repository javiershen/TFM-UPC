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
options=("1- Receive a medicament" "2- Dispense a medicament" "3- Use prescription" "4- Quit")
select opt in "${options[@]}"
do
    case $opt in
        "1- Receive a medicament")
            
            echo PRODUCT CODE:
            read product_code
            export PRODUCT_CODE="$product_code"
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["ReceiveMedicament", "'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${PRODUCT_CODE}"'"]}'
            break
            ;;
        "2- Dispense a medicament")
            echo PRODUCT CODE:
            read product_code
            export PRODUCT_CODE="$product_code"
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["DispenseMedicament", "'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${PRODUCT_CODE}"'"]}'
            break
            ;;
        "3- Use prescription")
            echo MEDICAMENT CODE:
            read medicament_code
            export MEDICAMENT_CODE="$medicament_code"
            echo PATIENT ID:
            read patient_ID
            export PATIENT_ID="$patient_ID"
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["ConsumePrescription", "'"${USERNAME_C2}"'", "'"${ENTITY_ID_C2}"'", "'"${MEDICAMENT_CODE}"'", "'"${PATIENT_ID}"'"]}'
            break
            ;;
        "4- Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
popd

