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
options=("1- Generate prescription" "2- See prescription" "3- Quit")
select opt in "${options[@]}"
do
    case $opt in
	"1- Generate prescription")
            echo MEDICAMENT CODE:
            read medicament_code
            export MEDICAMENT_CODE="$medicament_code"
            echo PATIENT ID:
            read patient_ID
            export PATIENT_ID="$patient_ID"
            echo EXPIRATION YEAR:
            read expiration_year
            export EXPIRATION_YEAR="$expiration_year"
            echo EXPIRATION MONTH:
            read expiration_month
            export EXPIRATION_MONTH="$expiration_month"      
                           
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"Invoke","Args":["GeneratePrescription", "'"${USERNAME_C2}"'", "'"${ENTITY_ID_C2}"'", "'"${MEDICAMENT_CODE}"'", "'"${PATIENT_ID}"'", "'"${EXPIRATION_YEAR}"'", "'"${EXPIRATION_MONTH}"'"]}'       	    
	    break
	    ;;
        "2- See prescription")
            echo PATIENT ID:
            read patient_id
            export PATIENT_ID="$patient_id"            
            echo MEDICAMENT CODE:
            read medicament_code
            export MEDICAMENT_CODE="$medicament_code"
            
            peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem" -C channel2 -n basic2 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.pharma.com/peers/peer0.org1.pharma.com/tls/ca.crt" --peerAddresses localhost:11051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.pharma.com/peers/peer0.org3.pharma.com/tls/ca.crt" -c '{"function":"GetPrescription","Args":["'"${USERNAME_C2}"'", "'"${ENTITY_ID_C2}"'", "'"${MEDICAMENT_CODE}"'", "'"${PATIENT_ID}"'"]}'       
            break
            ;;
        "3- Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
popd

