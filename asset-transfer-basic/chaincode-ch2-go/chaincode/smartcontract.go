package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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
	Medicament_Name string `json:"Medicament_Name"`
	Product_Code    int    `json:"Product_Code"`
}
type Prescription struct {
	DispensationDate       string `json:"DispensationDate"`
	GenerationDate         string `json:"GenerationDate"`
	Expiration_Month       int    `json:"Expiration_Month"`
	Expiration_Year        int    `json:"Expiration_Year"`
	PatientID              string `json:"PatientID"`
	Pharmacy_EntityID      string `json:"Pharmacy_EntityID"`
	Pharmacy_UserID        string `json:"Pharmacy_UserID"`
	Prescripted_Medicament int    `json:"Prescripted_Medicament"`
	Sanitary_EntityID      string `json:"Sanitary_EntityID"`
	Sanitary_UserID        string `json:"Sanitary_UserID"`
	Status                 int    `json:"Status"` // 1: creado | 2: utilizada en una pharma

}

type Session struct {
	EntityID       string `json:"EntityID"`
	UserID         string `json:"UserID"`
	GenerationDate string `json:"GenerationDate"`
	SessionID      string `json:"SessionID"`
	Status         int    `json:"Status"` // 1: activa | 0: inactiva
}

// invoke function to call tracking points functions
func (s *SmartContract) Invoke(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()
	correctArgs, err := s.areArgumentsCorrect(ctx, args[1:]) //first argument is Invoke function, we don't need it
	if err != nil {
		return err
	}
	if correctArgs {
		function := args[1]
		userID := args[2]
		entityID := args[3]
		sessionID := s.GenerateSessionID(ctx, entityID, userID)
		session, err := s.ReadSession(ctx, sessionID)
		if err != nil {
			return err
		}
		if session == nil {
			return fmt.Errorf("No active session for that user. Please, log in")
		}
		isexpired, err := s.IsSessionExpired(ctx, session)
		if err != nil {
			return err
		}
		if !isexpired {

			accessible, err := s.isFunctionAccessible(ctx, function, userID, entityID)
			if err != nil {
				return err
			}
			if accessible {
				allArgs := args[4:]
				if function == "GeneratePrescription" {
					return s.GeneratePrescription(ctx, entityID, userID, allArgs)
				} else if function == "ConsumePrescription" {
					return s.ConsumePrescription(ctx, entityID, userID, allArgs)
				}
				return fmt.Errorf("Invalid function")
			}
			return fmt.Errorf("Incorrect Args")
		} else {
			return fmt.Errorf("Session expired. Log in again, please")
		}
	}
	return fmt.Errorf("Function inaccessible")
}

func (s *SmartContract) IsPswCorrect(ctx contractapi.TransactionContextInterface, _UserID string, _psw string) (bool, error) {
	user, err := s.ReadUser(ctx, _UserID)
	if err != nil {
		return false, err
	}
	if user.Password != _psw {
		return false, fmt.Errorf("Wrong credentials")
	}
	return true, nil
}

func (s *SmartContract) CreateSession(ctx contractapi.TransactionContextInterface, _sessionID string, _entityID string, _userID string) error {
	actualDateStr, err := s.GetTxTimestamp(ctx)
	if err != nil {
		return err
	}
	session := Session{
		EntityID:       _entityID,
		UserID:         _userID,
		GenerationDate: actualDateStr,
		SessionID:      _sessionID,
		Status:         1,
	}
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(session.SessionID, sessionJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) IsSessionExpired(ctx contractapi.TransactionContextInterface, _session *Session) (bool, error) {
	if _session.Status == 1 {
		actualDateStr, err := s.GetTxTimestamp(ctx)
		sessionDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", _session.GenerationDate)
		if err != nil {
			return false, err
		}
		actualDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", actualDateStr)
		if err != nil {
			return false, err
		}
		if sessionDate.Sub(actualDate).Minutes() > 5 {
			_session.Status = 0
			assetJSON, err := json.Marshal(_session)
			if err != nil {
				return true, err
			}

			err = ctx.GetStub().PutState(_session.SessionID, assetJSON)
			if err != nil {
				return true, err
			}
			return true, fmt.Errorf("Session expired. Log in again, please")
		}
		return false, nil
	}
	return true, fmt.Errorf("Session expired. Log in again, please")
}

func (s *SmartContract) LogIn(ctx contractapi.TransactionContextInterface, _entityID string, _userName string, _psw string) error {
	entity, err := s.ReadEntity(ctx, _entityID)
	if err != nil {
		return err
	}

	user, err := s.ReadUser(ctx, _userName)
	if err != nil {
		return err
	}
	users := entity.Entity_Users

	isValidUser, err := s.isUserInEntity(ctx, user.User_ID, users)
	if err != nil {
		return err
	}
	if isValidUser {
		isPswCorrect, err := s.IsPswCorrect(ctx, user.User_ID, _psw)
		if err != nil {
			return err
		}
		if isPswCorrect {
			sessionID := s.GenerateSessionID(ctx, _entityID, user.User_ID)
			session, err := s.ReadSession(ctx, sessionID)
			if err != nil {
				return err
			}
			if session == nil {
				err := s.CreateSession(ctx, sessionID, _entityID, user.User_ID)
				if err != nil {
					return err
				}
			} else {
				isExpired, err := s.IsSessionExpired(ctx, session)
				if err != nil {
					return err
				}
				if isExpired {
					actualDateStr, err := s.GetTxTimestamp(ctx)
					if err != nil {
						return err
					}
					session.GenerationDate = actualDateStr
					session.Status = 1
					assetJSON, err := json.Marshal(session)
					if err != nil {
						return err
					}

					err = ctx.GetStub().PutState(session.SessionID, assetJSON)
					if err != nil {
						return err
					}
				}
			}
			return nil
		} else {
			return fmt.Errorf("Wrong credentials")
		}

	}
	return fmt.Errorf("Invalid user")
}

func (s *SmartContract) GenerateSessionID(ctx contractapi.TransactionContextInterface, _entityID string, _userID string) string {
	return _entityID + "-" + _userID
}

// InitLedger adds a base set of medicaments, entities and entity users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	sanitary_Users := []Entity_User{
		{User_Name: "sanitaryUser", User_ID: "sanitaryUser", Email: "sanitary@pg.com", Rol: "user", Address: "bangalore", Password: "adminpw"},
		{User_Name: "adminLab", User_ID: "adminLab", Email: "sanitary@pg.com", Rol: "admin", Address: "bangalore", Password: "adminpw"},
	}

	sanitary_UsersID := []string{}
	for _, sanitary_User := range sanitary_Users {

		sanitary_UsersID = append(sanitary_UsersID, sanitary_User.User_ID)
		sanitary_UserJSON, err := json.Marshal(sanitary_User)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(sanitary_User.User_ID, sanitary_UserJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	pharmacy_Users := []Entity_User{
		{User_Name: "pharmacyUser", User_ID: "pharmacyUser", Email: "pharmacy@pg.com", Rol: "user", Address: "bangalore", Password: "adminpw"},
		{User_Name: "pharmacyAdmin", User_ID: "pharmacyAdmin", Email: "pharmacy@pg.com", Rol: "admin", Address: "bangalore", Password: "adminpw"},
	}

	pharmacy_UsersID := []string{}
	for _, pharmacy_User := range pharmacy_Users {
		pharmacy_UsersID = append(pharmacy_UsersID, pharmacy_User.User_ID)
		pharmacy_UserJSON, err := json.Marshal(pharmacy_User)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(pharmacy_User.User_ID, pharmacy_UserJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %v", err)
		}
	}

	entities := []Entity{
		{Entity_Name: "hospital1", Entity_ID: "hospital1", Type: "sanitation", Entity_Users: sanitary_UsersID},
		{Entity_Name: "pharmacy2", Entity_ID: "pharmacy2", Type: "pharmacy", Entity_Users: pharmacy_UsersID},
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
		{Medicament_Name: "Ibuprofeno", Product_Code: 8470008722513},
		{Medicament_Name: "Paracetamol", Product_Code: 8470006723459},
	}

	for _, medicament := range medicaments {
		medJSON, err := json.Marshal(medicament)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(strconv.Itoa(medicament.Product_Code), medJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	_GenerationDate, err := s.GetTxTimestamp(ctx)
	if err != nil {
		return err
	}

	prescription := Prescription{
		DispensationDate:       "",
		GenerationDate:         _GenerationDate,
		Expiration_Month:       2023,
		Expiration_Year:        1,
		PatientID:              "1234567891ABCD",
		Pharmacy_EntityID:      "",
		Pharmacy_UserID:        "",
		Prescripted_Medicament: medicaments[0].Product_Code,
		Sanitary_UserID:        "hospital1",
		Sanitary_EntityID:      "sanitaryUser",
		Status:                 1,
	}
	assetJSON, err := json.Marshal(prescription)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(prescription.PatientID, assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// function to register a medicament
func (s *SmartContract) GeneratePrescription(ctx contractapi.TransactionContextInterface, _entityID string, _userID string, args []string) error {
	medicament_Code, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	_PatientID := args[1]
	//REVISAR QUE MEDICAMENTO EXISTA
	medicament, err := s.ReadMedicament(ctx, strconv.Itoa(medicament_Code))
	if err != nil {
		return err
	}
	if medicament != nil {
		//revisar que receta no esté caducada
		_Expiration_Year, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}
		_Expiration_Month, err := strconv.Atoi(args[3])
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

		_GenerationDate, err := s.GetTxTimestamp(ctx)
		if err != nil {
			return err
		}
		prescription := Prescription{
			DispensationDate:       "",
			GenerationDate:         _GenerationDate,
			Expiration_Year:        _Expiration_Year,
			Expiration_Month:       _Expiration_Month,
			PatientID:              _PatientID,
			Pharmacy_EntityID:      "",
			Pharmacy_UserID:        "",
			Prescripted_Medicament: medicament_Code,
			Sanitary_UserID:        _userID,
			Sanitary_EntityID:      _entityID,
			Status:                 1}

		//SI TODO VA BIEN, MEDICAMENTO ES REGISTRADO

		medJSON, err := json.Marshal(prescription)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(prescription.PatientID, medJSON)
		if err != nil {
			return err
		}

		return nil

	} else {
		return fmt.Errorf("Undefined Medicament")
	}
}

// function to register the receive of a medicament in a pharmacy
func (s *SmartContract) ConsumePrescription(ctx contractapi.TransactionContextInterface, _entityID string, _userID string, args []string) error {
	medicament_Code, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	//OBTENGO EL ASSET Y COMPRUEBO QUE EXISTE
	medicament, err := s.ReadMedicament(ctx, strconv.Itoa(medicament_Code))
	if err != nil {
		return err
	}
	if medicament != nil {
		_PatientID := args[1]
		prescription, err := s.ReadPrescription(ctx, _PatientID)
		if err != nil {
			return err
		}
		if prescription != nil && prescription.Prescripted_Medicament == medicament_Code {
			//REVISAR QUE NO ESTÉ CADUCADO
			expired, err := s.isExpired(ctx, prescription.Expiration_Year, prescription.Expiration_Month)
			if err != nil {
				return err
			}
			if expired {
				return fmt.Errorf("This medicament is expired, it can not be processed")
			}

			status, err := s.getNewStatus(ctx, prescription.Status, "ConsumePrescription")
			if err != nil {
				return err
			}

			_GenerationDate, err := s.GetTxTimestamp(ctx)
			if err != nil {
				return err
			}
			prescription.Status = status
			prescription.DispensationDate = _GenerationDate
			prescription.Pharmacy_EntityID = _entityID
			prescription.Pharmacy_UserID = _userID

			//ASSET ACTUALIZADO
			assetJSON, err := json.Marshal(prescription)
			if err != nil {
				return err
			}

			err = ctx.GetStub().PutState(strconv.Itoa(medicament_Code), assetJSON)
			if err != nil {
				return err
			}
			return nil

		} else {
			return fmt.Errorf("Invalid Prescription")
		}
	} else {
		return fmt.Errorf("Undefined Medicament")
	}

}

// function that returs the new medicament status after going through a tracking point
func (s *SmartContract) getNewStatus(ctx contractapi.TransactionContextInterface, _currentStatus int, _function string) (int, error) {

	if _function == "ConsumePrescription" && _currentStatus == 1 {
		newStatus := 0
		return newStatus, nil
	} else {
		return _currentStatus, fmt.Errorf("Status can not be modified")
	}
}

// function that checks if the arguments passed to a function are correct
func (s *SmartContract) areArgumentsCorrect(ctx contractapi.TransactionContextInterface, _Args []string) (bool, error) {
	if len(_Args) <= 3 {
		return false, fmt.Errorf("Incorrect number of arguments")
	}
	function := _Args[0]
	if len(function) == 0 {
		return false, fmt.Errorf("Function must be provided")
	}
	if len(_Args[1]) == 0 {
		return false, fmt.Errorf("User ID must be provided")
	}
	if len(_Args[2]) == 0 {
		return false, fmt.Errorf("Entity ID must be provided")
	}
	_allArgs := _Args[3:]
	if function == "GeneratePrescription" {
		if len(_allArgs) != 4 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Medicament Code must be provided")
		}
		if len(_allArgs[1]) == 0 {
			return false, fmt.Errorf("Patient ID must be provided")
		}
		if len(_allArgs[2]) == 0 {
			return false, fmt.Errorf("Expiration Year must be provided")
		}
		if len(_allArgs[3]) == 0 {
			return false, fmt.Errorf("Expiration Month must be provided")
		}
		return true, nil
	} else if function == "ConsumePrescription" {
		if len(_allArgs) != 2 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Medicament Code must be provided")
		}
		if len(_allArgs[1]) == 0 {
			return false, fmt.Errorf("Patient ID must be provided")
		}
		return true, nil
	} else if function == "LogIn" {
		if len(_allArgs) != 1 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Psw must be provided")
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
		if _function == "LogIn" {
			return true, nil
		}
		if entity.Type == "sanitation" {
			if _function == "GeneratePrescription" {
				return true, nil
			}

		} else if entity.Type == "pharmacy" {
			if _function == "ConsumePrescription" {
				return true, nil
			}
		} else {
			return false, fmt.Errorf("Invalid user to call method '" + _function + "'")
		}

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
func (s *SmartContract) ReadMedicament(ctx contractapi.TransactionContextInterface, _Product_Code string) (*Medicament, error) {
	assetJSON, err := ctx.GetStub().GetState(_Product_Code)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the medicament %s does not exist", _Product_Code)
	}

	var medicament Medicament
	err = json.Unmarshal(assetJSON, &medicament)
	if err != nil {
		return nil, err
	}

	return &medicament, nil
}

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, _PatientID string) (*Prescription, error) {
	assetJSON, err := ctx.GetStub().GetState(_PatientID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the prescription %s does not exist", _PatientID)
	}

	var prescription Prescription
	err = json.Unmarshal(assetJSON, &prescription)
	if err != nil {
		return nil, err
	}

	return &prescription, nil
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
		return nil, fmt.Errorf("the entity %s does not exist", _Entity_ID)
	}

	var entity Entity
	err = json.Unmarshal(entityJSON, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// function that returns a registered session given an session id
func (s *SmartContract) ReadSession(ctx contractapi.TransactionContextInterface, _SessionID string) (*Session, error) {
	sessionJSON, err := ctx.GetStub().GetState(_SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if sessionJSON == nil {
		return nil, nil
	}

	var session Session
	err = json.Unmarshal(sessionJSON, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SmartContract) GetTxTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	txTimeAsPtr, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", err
	}
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()

	return timeStr, nil
}

// function that gets a registered medicament given a serial number and returns it after checking the current owner of the medicament
func (s *SmartContract) GetMedicament(ctx contractapi.TransactionContextInterface, _userID string, _entityID string, _Product_Code string) (*Medicament, error) {
	sessionID := s.GenerateSessionID(ctx, _entityID, _userID)
	session, err := s.ReadSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("No active session for that user. Please, log in")
	}
	isexpired, err := s.IsSessionExpired(ctx, session)
	if err != nil {
		return nil, err
	}
	if !isexpired {
		user, err := s.ReadUser(ctx, _userID)
		if err != nil {
			return nil, err
		}
		if user.Rol == "admin" {
			medicament, err := s.ReadMedicament(ctx, _Product_Code)
			if err != nil {
				return nil, err
			}
			return medicament, nil
		} else {
			return nil, fmt.Errorf("Can not access to a medicament without being the admin")
		}
	} else {
		return nil, fmt.Errorf("Session expired. Log in again, please")
	}
}

// function that returns all registered medicaments owner by an entity in the system
func (s *SmartContract) GetAllMedicaments(ctx contractapi.TransactionContextInterface, _userID string, _entityID string) ([]*Medicament, error) { //
	sessionID := s.GenerateSessionID(ctx, _entityID, _userID)
	session, err := s.ReadSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("No active session for that user. Please, log in")
	}
	isexpired, err := s.IsSessionExpired(ctx, session)
	if err != nil {
		return nil, err
	}
	if !isexpired {
		user, err := s.ReadUser(ctx, _userID)
		if err != nil {
			return nil, err
		}
		if user.Rol == "admin" {

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
			}

			return medicaments, nil
		} else {
			return nil, fmt.Errorf("Can not access to medicaments without being the admin")
		}
	} else {
		return nil, fmt.Errorf("Session expired. Log in again, please")
	}
}

// function that returns all registered users of an entity in the system
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface, _userID string, _entityID string) ([]*Entity_User, error) {
	sessionID := s.GenerateSessionID(ctx, _entityID, _userID)
	session, err := s.ReadSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("No active session for that user. Please, log in")
	}
	isexpired, err := s.IsSessionExpired(ctx, session)
	if err != nil {
		return nil, err
	}
	if !isexpired {

		entity, err := s.ReadEntity(ctx, _entityID)
		if err != nil {
			return nil, err
		}

		users := entity.Entity_Users

		currentUser, err := s.ReadUser(ctx, _userID)
		if err != nil {
			return nil, err
		}
		if currentUser.Rol == "admin" {
			var usersEntity []*Entity_User
			for _, userEntityID := range users {
				user, err := s.ReadUser(ctx, userEntityID)
				if err != nil {
					return nil, err
				}
				usersEntity = append(usersEntity, user)
			}
			return usersEntity, nil
		} else {
			return nil, fmt.Errorf("Can not access to users without being the admin")
		}

	} else {
		return nil, fmt.Errorf("Session expired. Log in again, please")
	}
}