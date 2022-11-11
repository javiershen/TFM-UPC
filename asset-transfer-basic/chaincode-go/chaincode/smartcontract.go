package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	/*
		 	"github.com/hyperledger/fabric/core/chaincode/shim"
			pb "github.com/hyperledger/fabric/protos/peer"
	*/)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Entity struct {
	Entity_Name  string        `json:"Name"`
	Entity_ID    string        `json:"ID"`
	Type         string        `json:"Type"`
	Entity_Users []Entity_User `json:"Entity_Users"`
}

type Entity_User struct {
	User_Name string `json:"Name"`
	User_ID   string `json:"ID"`
	Email     string `json:"Email"`
	Rol       string `json:"Type"`
	Address   string `json:"Address"`
	Password  string `json:"Password"`
}

type Medicament struct {
	Expiration_Month int    `json:"Expiration_Month"`
	Expiration_Year  int    `json:"Expiration_Year"`
	Lot_Number       string `json:"LotNumber"`
	Medicament_Name  string `json:"Name"`
	Product_Code     int    `json:"Product_Code"`
	Serial_Number    string `json:"Serial_Number"`
	Status           int    `json:"Status"` // 1: creado | 2: despachado de lab | 3: recibido por farmacia | 4: dispensado | 5: indispensable por motivo que sea
	Entity_Producer  Entity `json:"Entity_Producer"`
	Entity_Owner     Entity `json:"Entity_Owner"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	lab_Users := []Entity_User{
		{User_Name: "lab", User_ID: "lab", Email: "lab@pg.com", Rol: "admin", Address: "bangalore", Password: "adminpw"},
	}

	for _, lab_User := range lab_Users {
		assetJSON, err := json.Marshal(lab_User)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(lab_User.User_ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	pharmacy_Users := []Entity_User{
		{User_Name: "pharmacy", User_ID: "pharmacy", Email: "pharmacy@pg.com", Rol: "admin", Address: "bangalore", Password: "adminpw"},
	}

	for _, pharmacy_User := range pharmacy_Users {
		assetJSON, err := json.Marshal(pharmacy_User)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(pharmacy_User.User_ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	entities := []Entity{
		{Entity_Name: "lab1", Entity_ID: "lab1", Type: "lab", Entity_Users: lab_Users},
		{Entity_Name: "pharmacy2", Entity_ID: "pharmacy2", Type: "pharmacy", Entity_Users: pharmacy_Users},
	}

	for _, entity := range entities {
		assetJSON, err := json.Marshal(entity)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(entity.Entity_ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	medicaments := []Medicament{
		{Medicament_Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352687", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Medicament_Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352688", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Medicament_Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352689", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Medicament_Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352687", Lot_Number: "L101JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Medicament_Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352688", Lot_Number: "L101JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Medicament_Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352689", Lot_Number: "L201JX30", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
	}

	for _, medicament := range medicaments {
		medJSON, err := json.Marshal(medicament)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(medicament.Serial_Number, medJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) getStatus(ctx contractapi.TransactionContextInterface, _currentStatus int, _function string) (int, error) {

	if (_function == "DispatchMedicament" && _currentStatus == 1) ||
		(_function == "ReceiveMedicament" && _currentStatus == 2) ||
		(_function == "DispenseMedicament" && _currentStatus == 3) {
		newStatus := _currentStatus + 1
		return newStatus, nil
	} else {
		return _currentStatus, fmt.Errorf("Status can not be modified")
	}

}

func (s *SmartContract) isUserInEntity(ctx contractapi.TransactionContextInterface, _user Entity_User, _entity Entity) (bool, error) {

	for i := 0; i < len(_entity.Entity_Users); i++ {
		if _entity.Entity_Users[i] == _user {
			return true, nil
		}
	}
	return false, nil

}

func (s *SmartContract) isFunctionAccessible(ctx contractapi.TransactionContextInterface, _function string, _userID string, _entityID string) (bool, error) {
	entity, err := s.ReadEntity(ctx, _entityID)

	if err != nil {
		return false, err
	}
	user, err := s.ReadUser(ctx, _userID)
	if err != nil {
		return false, err
	}
	isValidUser, err := s.isUserInEntity(ctx, user, entity)
	if err != nil {
		return false, err
	}
	if isValidUser {

		if entity.Type == "lab" {
			if _function == "RegisterMedicament" || _function == "DispatchMedicament" {
				return true, nil
			}

		} else if entity.Type == "pharmacy" {
			if _function == "ReceiveMedicament" || _function == "DispenseMedicament" {
				return true, nil
			}
		}
		return false, nil

	}
	return false, nil

}

func (s *SmartContract) isExpired(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return true, err
	}
	expYear := medicament.Expiration_Year
	currentDate := time.Now()
	if expYear < int(currentDate.Year()) {
		return true, nil
	} else if expYear == int(currentDate.Year()) {
		expMonth := medicament.Expiration_Month
		if expMonth < int(currentDate.Month()) {
			return true, nil
		}
	}
	return false, nil
}
func (s *SmartContract) RegisterMedicament(ctx contractapi.TransactionContextInterface, _Name string, _Product_Code int, _Serial_Number string, _Lot_Number string, _Expiration_Year int, _Expiration_Month int) (bool, error) {

	//REVISAR QUE NO ESTÉ CREADO
	exists, err := s.AssetExists(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if exists {
		return false, fmt.Errorf("This medicament already exists")
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//SI TODO VA BIEN, MEDICAMENTO ES REGISTRADO
	medicament := Medicament{
		Medicament_Name:  _Name,
		Product_Code:     _Product_Code,
		Serial_Number:    _Serial_Number,
		Lot_Number:       _Lot_Number,
		Expiration_Year:  _Expiration_Year,
		Expiration_Month: _Expiration_Month,
		Status:           1,
	}
	medJSON, err := json.Marshal(medicament)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(medicament.Serial_Number, medJSON)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *SmartContract) DispatchMedicament(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {

	//OBTENGO EL ASSET Y COMPRUEBO QUE EXISTE
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//SI TODO VA BIEN:

	//SE SETEA EL FUTURE OWNER

	//EL STATUS CAMBIA
	status, err := s.getStatus(ctx, medicament.Status, "DispatchMedicament")
	if err != nil {
		return false, err
	}

	medicament.Status = status

	//ASSET ACTUALIZADO
	assetJSON, err := json.Marshal(medicament)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(_Serial_Number, assetJSON)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (s *SmartContract) ReceiveMedicament(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {

	//OBTENGO EL ASSET Y COMPRUEBO QUE EXISTE
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//REVISAR QUE EL FUTURE OWNER DEL FARMACO CORRESPONDE A LA FARMACIA QUE EJECUTA ESTA FUNCION

	//EL STATUS CAMBIA
	status, err := s.getStatus(ctx, medicament.Status, "DispatchMedicament")
	if err != nil {
		return false, err
	}

	medicament.Status = status

	//ASSET ACTUALIZADO
	assetJSON, err := json.Marshal(medicament)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(_Serial_Number, assetJSON)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SmartContract) DispenseMedicament(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {

	//OBTENGO EL ASSET Y COMPRUEBO QUE EXISTE
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//REVISAR QUE EL CURRENT OWNER DEL FARMACO CORRESPONDE A LA FARMACIA QUE EJECUTA ESTA FUNCION

	//EL STATUS CAMBIA
	status, err := s.getStatus(ctx, medicament.Status, "DispatchMedicament")
	if err != nil {
		return false, err
	}

	medicament.Status = status

	//ASSET ACTUALIZADO
	assetJSON, err := json.Marshal(medicament)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(_Serial_Number, assetJSON)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadMedicament(ctx contractapi.TransactionContextInterface, _Serial_Number string) (*Medicament, error) {
	assetJSON, err := ctx.GetStub().GetState(_Serial_Number)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", _Serial_Number)
	}

	var medicament Medicament
	err = json.Unmarshal(assetJSON, &medicament)
	if err != nil {
		return nil, err
	}

	return &medicament, nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, _User_ID string) (*Entity_User, error) {
	userJSON, err := ctx.GetStub().GetState(_User_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the user %s does not exist", _User_ID)
	}

	var user Entity_User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadEntity(ctx contractapi.TransactionContextInterface, _Entity_ID string) (*Entity, error) {
	entityJSON, err := ctx.GetStub().GetState(_Entity_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if entityJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", _Entity_ID)
	}

	var entity Entity
	err = json.Unmarshal(entityJSON, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(_Serial_Number)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

/* // GetAllAssets returns all assets found in world state
func (t *SmartContract) GetAllMedicaments(ctx contractapi.TransactionContextInterface) ([]*Medicament, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var medicaments []*Medicament
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var medicament Medicament
		err = json.Unmarshal(queryResponse.Value, &medicament)
		if err != nil {
			return nil, err
		}
		medicaments = append(medicaments, &medicament)
	}

	return medicaments, nil
} */

/* func (t *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*Entity_User, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var entity_Users []*Entity_User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var entity_User Entity_User
		err = json.Unmarshal(queryResponse.Value, &entity_User)
		if err != nil {
			return nil, err
		}
		entity_Users = append(entity_Users, &entity_User)
	}

	return entity_Users, nil
} */
