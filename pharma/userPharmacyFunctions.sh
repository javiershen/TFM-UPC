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
options=("Receive a medicament" "Use prescription and dispense medicament" "Quit")
select opt in "${options[@]}"
do
    case $opt in
        "Receive a medicament")
            echo MEDICAMENT SERIAL NUMBER:
            read medicament_sn
            export MEDICAMENT_SN="$medicament_sn"
            echo MEDICAMENT CODE:
            read medicament_code
            export MEDICAMENT_CODE="$medicament_code"
            echo MEDICAMENT NAME:
            read medicine_name
            export MEDICINE_NAME="$medicine_name"           
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["ReceiveMedicament", "'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${MEDICAMENT_SN}"'"]}'
              
           peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["AddMedicamentToStock", "'"${USERNAME_C2}"'", "'"${ENTITY_ID_C2}"'", "'"${MEDICAMENT_CODE}"'", "'"${MEDICINE_NAME}"'"]}'
            break
            ;;
        "Use prescription and dispense medicament")
            echo MEDICAMENT CODE:
            read medicament_code
            export MEDICAMENT_CODE="$medicament_code"
            echo MEDICAMENT SERIAL NUMBER:
            read medicament_sn
            export MEDICAMENT_SN="$medicament_sn"
            echo PATIENT ID:
            read patient_ID
            export PATIENT_ID="$patient_ID"
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["ConsumePrescription", "'"${USERNAME_C2}"'", "'"${ENTITY_ID_C2}"'", "'"${MEDICAMENT_CODE}"'", "'"${PATIENT_ID}"'"]}'
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel1 -n basic1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.pharma.com/peers/peer0.org2.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["DispenseMedicament", "'"${USERNAME_C1}"'", "'"${ENTITY_ID_C1}"'", "'"${MEDICAMENT_SN}"'"]}'
            break
            ;;
        
        "Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
popd

