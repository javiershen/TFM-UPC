package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	//"github.com/hyperledger/fabric-contract-api-go@v1.2.0/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	//"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Entity struct {
	Entity_Name  string   `json:"Entity_Name"`
	Entity_ID    string   `json:"Entity_ID"`
	Type         string   `json:"Type"`
	Entity_Users []string `json:"Entity_Users"`
}

type Entity_User struct {
	User_Name string `json:"User_Name"`
	User_ID   string `json:"User_ID"`
	Email     string `json:"Email"`
	Rol       string `json:"Rol"`
	Address   string `json:"Address"`
	Password  string `json:"Password"`
}

type Medicament struct {
	Expiration_Month int    `json:"Expiration_Month"`
	Expiration_Year  int    `json:"Expiration_Year"`
	Lot_Number       string `json:"Lot_Number"`
	Medicament_Name  string `json:"Medicament_Name"`
	Product_Code     int    `json:"Product_Code"`
	Serial_Number    string `json:"Serial_Number"`
	Status           int    `json:"Status"` // 1: creado | 2: despachado de lab | 3: recibido por farmacia | 4: dispensado | 5: indispensable por motivo que sea
	Producer_Lab     string `json:"Producer_Lab"`
	Seller_Pharmacy  string `json:"Seller_Pharmacy"`
	Current_Owner    string `json:"Current_Owner"`
}

// invoke function to call tracking points functions
func (s *SmartContract) InvokeTrackingPoint(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()

	function := args[1]
	userID := args[2]
	entityID := args[3]
	accessible, err := s.isFunctionAccessible(ctx, function, userID, entityID)
	if err != nil {
		return err
	}
	if accessible {
		allArgs := args[4:]
		correctArgs, err := s.areArgumentsCorrect(ctx, function, allArgs)
		if err != nil {
			return err
		}
		if correctArgs {
			if function == "RegisterMedicament" {
				return s.RegisterMedicament(ctx, allArgs)
			} else if function == "DispatchMedicament" {
				return s.RegisterMedicament(ctx, allArgs)
			} else if function == "ReceiveMedicament" {
				return s.RegisterMedicament(ctx, allArgs)
			} else if function == "DispenseMedicament" {
				return s.RegisterMedicament(ctx, allArgs)
			}
			return fmt.Errorf("Invalid function")
		}
		return fmt.Errorf("Incorrect Args")
	}
	return fmt.Errorf("Function inaccessible")
}

// InitLedger adds a base set of medicaments, entities and entity users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	lab_User := Entity_User{
		User_Name: "lab", User_ID: "lab_ID", Email: "lab@pg.com", Rol: "admin", Address: "bangalore", Password: "adminpw",
	}

	lab_UserJSON, err := json.Marshal(lab_User)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(lab_User.User_ID, lab_UserJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	pharmacy_User := Entity_User{
		User_Name: "pharmacy", User_ID: "pharmacy_ID", Email: "pharmacy@pg.com", Rol: "admin", Address: "bangalore", Password: "adminpw",
	}

	pharmacy_UserJSON, err := json.Marshal(pharmacy_User)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(pharmacy_User.User_ID, pharmacy_UserJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	entities := []Entity{
		{Entity_Name: "lab1", Entity_ID: "lab1", Type: "lab", Entity_Users: []string{"lab_ID"}},
		{Entity_Name: "pharmacy2", Entity_ID: "pharmacy2", Type: "pharmacy", Entity_Users: []string{"pharmacy_ID"}},
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
		{Medicament_Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352687", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1, Producer_Lab: "lab1", Seller_Pharmacy: "", Current_Owner: "lab1"},
		{Medicament_Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352688", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1, Producer_Lab: "lab1", Seller_Pharmacy: "", Current_Owner: "lab1"},
		{Medicament_Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352687", Lot_Number: "L101JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1, Producer_Lab: "lab1", Seller_Pharmacy: "", Current_Owner: "lab1"},
		{Medicament_Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352688", Lot_Number: "L101JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1, Producer_Lab: "lab1", Seller_Pharmacy: "", Current_Owner: "lab1"},
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

// function to register a medicament
func (s *SmartContract) RegisterMedicament(ctx contractapi.TransactionContextInterface, args []string) error {

	//REVISAR QUE NO ESTÉ CREADO
	exists, err := s.AssetExists(ctx, args[2])
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("This medicament already exists")
	}
	_Product_Code, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	_Expiration_Year, err := strconv.Atoi(args[4])
	if err != nil {
		return err
	}
	_Expiration_Month, err := strconv.Atoi(args[5])
	if err != nil {
		return err
	}
	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Expiration_Year, _Expiration_Month)
	if err != nil {
		return err
	}
	if expired {
		return fmt.Errorf("This medicament is expired, it can not be processed")
	}
	fmt.Printf("creating medicament with name: " + args[0] + "Serial_Number: " + args[2])
	//SI TODO VA BIEN, MEDICAMENTO ES REGISTRADO
	medicament := Medicament{
		Medicament_Name:  args[0],
		Product_Code:     _Product_Code,
		Serial_Number:    args[2],
		Lot_Number:       args[3],
		Expiration_Year:  _Expiration_Year,
		Expiration_Month: _Expiration_Month,
		Status:           1,
		Producer_Lab:     "",
		Seller_Pharmacy:  "",
		Current_Owner:    "",
	}
	medJSON, err := json.Marshal(medicament)
	if err != nil {
		return err
	}
	fmt.Printf("creating medicament with Serial_Number: " + medicament.Serial_Number)
	err = ctx.GetStub().PutState(medicament.Serial_Number, medJSON)
	if err != nil {
		return err
	}

	return nil
}

// function to register the dispatch of a medicament
func (s *SmartContract) DispatchMedicament(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {

	//OBTENGO EL ASSET Y COMPRUEBO QUE EXISTE
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, medicament.Expiration_Year, medicament.Expiration_Month)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//SI TODO VA BIEN:

	//SE SETEA EL FUTURE OWNER

	//EL STATUS CAMBIA
	status, err := s.getNewStatus(ctx, medicament.Status, "DispatchMedicament")
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

// function to register the receive of a medicament in a pharmacy
func (s *SmartContract) ReceiveMedicament(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {

	//OBTENGO EL ASSET Y COMPRUEBO QUE EXISTE
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, medicament.Expiration_Year, medicament.Expiration_Month)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//REVISAR QUE EL FUTURE OWNER DEL FARMACO CORRESPONDE A LA FARMACIA QUE EJECUTA ESTA FUNCION

	//EL STATUS CAMBIA
	status, err := s.getNewStatus(ctx, medicament.Status, "ReceiveMedicament")
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

// function to register the dispense of a medicament in a pharmacy
func (s *SmartContract) DispenseMedicament(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {

	//OBTENGO EL ASSET Y COMPRUEBO QUE EXISTE
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, medicament.Expiration_Year, medicament.Expiration_Month)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//REVISAR QUE EL CURRENT OWNER DEL FARMACO CORRESPONDE A LA FARMACIA QUE EJECUTA ESTA FUNCION

	//EL STATUS CAMBIA
	status, err := s.getNewStatus(ctx, medicament.Status, "DispenseMedicament")
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

// function that returs the new medicament status after going through a tracking point
func (s *SmartContract) getNewStatus(ctx contractapi.TransactionContextInterface, _currentStatus int, _function string) (int, error) {

	if (_function == "DispatchMedicament" && _currentStatus == 1) ||
		(_function == "ReceiveMedicament" && _currentStatus == 2) ||
		(_function == "DispenseMedicament" && _currentStatus == 3) {
		newStatus := _currentStatus + 1
		return newStatus, nil
	} else {
		return _currentStatus, fmt.Errorf("Status can not be modified")
	}
}

// function that checks if the arguments passed to a function are correct
func (s *SmartContract) areArgumentsCorrect(ctx contractapi.TransactionContextInterface, function string, _allArgs []string) (bool, error) {

	if function == "RegisterMedicament" {
		if len(_allArgs) != 6 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Medicament Name must be provided")
		}
		if len(_allArgs[1]) == 0 {
			return false, fmt.Errorf("Product Code must be provided")
		}
		if len(_allArgs[2]) == 0 {
			return false, fmt.Errorf("Serial Number must be provided")
		}
		if len(_allArgs[3]) == 0 {
			return false, fmt.Errorf("Lot Number must be provided")
		}
		if len(_allArgs[4]) == 0 {
			return false, fmt.Errorf("Expiration Year must be provided")
		}
		if len(_allArgs[5]) == 0 {
			return false, fmt.Errorf("Expiration Month must be provided")
		}
		return true, nil
	} else if function == "DispatchMedicament" {
		if len(_allArgs) != 1 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Serial Number must be provided")
		}
		return true, nil
	} else if function == "ReceiveMedicament" {
		if len(_allArgs) != 1 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Serial Number must be provided")
		}
		return true, nil
	} else if function == "DispenseMedicament" {
		if len(_allArgs) != 1 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Serial Number must be provided")
		}
		return true, nil
	}
	return false, fmt.Errorf("Incorrect function")
}

// function that checks if a user belongs to an entity
func (s *SmartContract) isUserInEntity(ctx contractapi.TransactionContextInterface, _userID string, _entity_users []string) (bool, error) {

	for i := 0; i < len(_entity_users); i++ {
		if _entity_users[i] == _userID {
			return true, nil
		}
	}
	return false, fmt.Errorf("User not valid")
}

// function that checks if a function is accessible by a determinate user
func (s *SmartContract) isFunctionAccessible(ctx contractapi.TransactionContextInterface, _function string, _userID string, _entityID string) (bool, error) {
	entity, err := s.ReadEntity(ctx, _entityID)

	if err != nil {
		return false, err
	}
	user, err := s.ReadUser(ctx, _userID)
	if err != nil {
		return false, err
	}

	users := entity.Entity_Users
	isValidUser, err := s.isUserInEntity(ctx, user.User_ID, users)
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
		return false, fmt.Errorf("Invalid user to call method '" + _function + "'")

	}
	return false, fmt.Errorf("User not registered")
}

// function that checks if a medicament is expired
func (s *SmartContract) isExpired(ctx contractapi.TransactionContextInterface, _Expiration_Year int, _Expiration_Month int) (bool, error) {

	currentDate := time.Now()
	if _Expiration_Year < int(currentDate.Year()) {
		return true, nil
	} else if _Expiration_Year == int(currentDate.Year()) {
		if _Expiration_Month < int(currentDate.Month()) {
			return true, nil
		}
	}
	return false, nil
}

// function that returns a registered medicament given a serial number
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

// function that returns a registered user given a user id
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

// function that returns a registered entity given an entity id
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

// function that returns all registered medicaments in the system
func (t *SmartContract) GetAllMedicaments(ctx contractapi.TransactionContextInterface) ([]*Medicament, error) { //

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
}

// function that returns all registered users in the system
func (t *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*Entity_User, error) {

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
}
