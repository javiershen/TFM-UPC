package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

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
	User_Name    string   `json:"User_Name"`
	User_ID      string   `json:"User_ID"`
	Email        string   `json:"Email"`
	Rol          string   `json:"Rol"`
	Address      string   `json:"Address"`
	Password     string   `json:"Password"`
	Sessions_Log []string `json:"Sessions_Log"`
}

type Counter struct {
	Count int `json:"Count"`
}

type MedicamentDates struct {
	DispatchDate     string `json:"DispatchDate"`
	DispenseDate     string `json:"DispenseDate"`
	ReceiveDate      string `json:"ReceiveDate"`
	RegistrationDate string `json:"RegistrationDate"`
}

type Medicament struct {
	Current_Owner    string          `json:"Current_Owner"`
	Dates            MedicamentDates `json:"Dates"`
	Expiration_Month int             `json:"Expiration_Month"`
	Expiration_Year  int             `json:"Expiration_Year"`
	Lot_Number       string          `json:"Lot_Number"`
	Medicament_Name  string          `json:"Medicament_Name"`
	Producer_Lab     string          `json:"Producer_Lab"`
	Product_Code     int             `json:"Product_Code"`
	Seller_Pharmacy  string          `json:"Seller_Pharmacy"`
	Serial_Number    string          `json:"Serial_Number"`
	Status           int             `json:"Status"` // 1: created | 2: dispatched from lab | 3: received by farmacy | 4: dispensed by farmacy
}

type Session struct {
	EntityID       string `json:"EntityID"`
	GenerationDate string `json:"GenerationDate"`
	SessionID      string `json:"SessionID"`
	Status         int    `json:"Status"` // 1: active | 0: inactive
}

// invoke function to call tracking points functions
func (s *SmartContract) Invoke(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetStringArgs()
	correctArgs, err := s.areArgumentsCorrect(ctx, args[1:])
	if err != nil {
		return err
	}
	if correctArgs {
		function := args[1]
		userID := args[2]
		entityID := args[3]

		activeSessions, err := s.GetActiveSessions(ctx, entityID, userID)
		if err != nil {
			return err
		}
		if len(activeSessions) == 1 {
			session := activeSessions[0]
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
					if function == "RegisterMedicament" {
						return s.RegisterMedicament(ctx, entityID, allArgs)
					} else if function == "DispatchMedicament" {
						return s.DispatchMedicament(ctx, entityID, allArgs)
					} else if function == "ReceiveMedicament" {
						return s.ReceiveMedicament(ctx, entityID, allArgs)
					} else if function == "DispenseMedicament" {
						return s.DispenseMedicament(ctx, entityID, allArgs)
					}
					return fmt.Errorf("Invalid function")
				}

				return fmt.Errorf("Incorrect Args")
			} else {
				return fmt.Errorf("Session expired. Log in again, please")
			}
		} else if len(activeSessions) == 0 {
			return fmt.Errorf("No active session for that user. Please, log in")

		} else {
			s.LogOut(ctx, entityID, userID)
			return fmt.Errorf("No active session for that user. Please, log in")
		}
	}
	return fmt.Errorf("Function inaccessible")
}

// function used to log in with a user of a given entity
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
			sessions, err := s.GetActiveSessions(ctx, _entityID, user.User_ID)
			if err != nil {
				return err
			}
			if len(sessions) > 0 {
				s.LogOut(ctx, _entityID, user.User_ID)
			}
			sessionID := s.GenerateSessionID(ctx)

			s.CreateSession(ctx, sessionID, _entityID, user.User_ID)
			s.StoreSession(ctx, sessionID, user.User_ID)

			return nil
		} else {
			return fmt.Errorf("Wrong credentials")
		}

	}
	return fmt.Errorf("Invalid user")
}

// InitLedger adds a base set of medicaments, entities and entity users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	//creation of lab users
	lab_Users := []Entity_User{
		{User_Name: "userLab", User_ID: "userLab", Email: "lab@pg.com", Rol: "user", Address: "bangalore", Password: "psw", Sessions_Log: []string{}},
		{User_Name: "adminLab", User_ID: "adminLab", Email: "lab@pg.com", Rol: "admin", Address: "bangalore", Password: "psw", Sessions_Log: []string{}},
	}

	lab_UsersID := []string{}
	for _, lab_User := range lab_Users {

		lab_UsersID = append(lab_UsersID, lab_User.User_ID)
		lab_UserJSON, err := json.Marshal(lab_User)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(lab_User.User_ID, lab_UserJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	//creation of pharmacy users
	pharmacy_Users := []Entity_User{
		{User_Name: "userPharmacy", User_ID: "userPharmacy", Email: "pharmacy@pg.com", Rol: "user", Address: "bangalore", Password: "psw", Sessions_Log: []string{}},
		{User_Name: "adminPharmacy", User_ID: "adminPharmacy", Email: "pharmacy@pg.com", Rol: "admin", Address: "bangalore", Password: "psw", Sessions_Log: []string{}},
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
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	//creation of lab and pharmacy entities
	entities := []Entity{
		{Entity_Name: "lab1", Entity_ID: "lab1", Type: "lab", Entity_Users: lab_UsersID},
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

	//creation of initial medicaments
	_RegisterDate, err := s.GetTxTimestamp(ctx)
	if err != nil {
		return err
	}
	MedDates := MedicamentDates{
		DispatchDate:     "",
		DispenseDate:     "",
		ReceiveDate:      "",
		RegistrationDate: _RegisterDate,
	}

	medicaments := []Medicament{
		{Medicament_Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352687", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1, Producer_Lab: "lab1", Seller_Pharmacy: "", Current_Owner: "lab1", Dates: MedDates},
		{Medicament_Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352687", Lot_Number: "L101JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1, Producer_Lab: "lab1", Seller_Pharmacy: "", Current_Owner: "lab1", Dates: MedDates},
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

	//start counter
	counter := Counter{Count: 1}
	counterJSON, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState("SessionCounter", counterJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// function to register a medicament
func (s *SmartContract) RegisterMedicament(ctx contractapi.TransactionContextInterface, _entityID string, args []string) error {

	//check that medicament does not exist
	_Serial_Number := args[2]
	exists, err := s.MedicamentExists(ctx, _Serial_Number)
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

	//check that medicament is not expired
	expired, err := s.isExpired(ctx, _Expiration_Year, _Expiration_Month)
	if err != nil {
		return err
	}
	if expired {
		return fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//initialization and update of dates
	emptyDates := MedicamentDates{
		DispatchDate:     "",
		DispenseDate:     "",
		ReceiveDate:      "",
		RegistrationDate: "",
	}
	MedDates, err := s.UpdateDates(ctx, "RegisterMedicament", emptyDates)
	if err != nil {
		return err
	}

	//creation of medicament
	medicament := Medicament{
		Medicament_Name:  args[0],
		Product_Code:     _Product_Code,
		Serial_Number:    _Serial_Number,
		Lot_Number:       args[3],
		Expiration_Year:  _Expiration_Year,
		Expiration_Month: _Expiration_Month,
		Status:           1,
		Producer_Lab:     _entityID,
		Seller_Pharmacy:  "",
		Current_Owner:    _entityID,
		Dates:            MedDates,
	}

	//update medicament on ledger
	medJSON, err := json.Marshal(medicament)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(medicament.Serial_Number, medJSON)
	if err != nil {
		return err
	}

	return nil
}

// function to register the dispatch of a medicament
func (s *SmartContract) DispatchMedicament(ctx contractapi.TransactionContextInterface, _entityID string, args []string) error {

	//obtain medicament to dispatch
	_Serial_Number := args[1]
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return err
	}

	//Check current owner
	if medicament.Current_Owner != _entityID {
		return fmt.Errorf("This medicament does not belong to this lab")
	}

	//Check future owner
	receiverEntity, err := s.ReadEntity(ctx, args[0])
	if err != nil {
		return err
	}
	if receiverEntity.Type != "pharmacy" {
		return fmt.Errorf("The recipient of the medicament must be a pharmacy")
	}

	//check that medicament is not expired
	expired, err := s.isExpired(ctx, medicament.Expiration_Year, medicament.Expiration_Month)
	if err != nil {
		return err
	}
	if expired {
		return fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//set pharmacy ID that will receive the medicament
	medicament.Seller_Pharmacy = receiverEntity.Entity_ID

	//modify medicament status
	status, err := s.getNewStatus(ctx, medicament.Status, "DispatchMedicament")
	if err != nil {
		return err
	}

	medicament.Status = status

	//update dates
	MedDates, err := s.UpdateDates(ctx, "DispatchMedicament", medicament.Dates)
	if err != nil {
		return err
	}
	medicament.Dates = MedDates

	//update medicament on ledger
	assetJSON, err := json.Marshal(medicament)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(_Serial_Number, assetJSON)
	if err != nil {
		return err
	}
	return nil
}

// function to register the receive of a medicament in a pharmacy
func (s *SmartContract) ReceiveMedicament(ctx contractapi.TransactionContextInterface, _entityID string, args []string) error {
	//obtain medicament that is received
	Serial_Number := args[0]
	medicament, err := s.ReadMedicament(ctx, Serial_Number)
	if err != nil {
		return err
	}

	//check that the medicament is being received by the correct entity
	if medicament.Seller_Pharmacy != _entityID {
		return fmt.Errorf("This medicament is not addressed to this pharmacy")
	}

	//check that medicament is not expired
	expired, err := s.isExpired(ctx, medicament.Expiration_Year, medicament.Expiration_Month)
	if err != nil {
		return err
	}
	if expired {
		return fmt.Errorf("This medicament is expired, it can not be processed")
	}

	medicament.Current_Owner = _entityID

	//modify medicament status
	status, err := s.getNewStatus(ctx, medicament.Status, "ReceiveMedicament")
	if err != nil {
		return err
	}

	medicament.Status = status

	//update dates
	MedDates, err := s.UpdateDates(ctx, "ReceiveMedicament", medicament.Dates)
	if err != nil {
		return err
	}
	medicament.Dates = MedDates

	//update medicament on ledger
	assetJSON, err := json.Marshal(medicament)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(Serial_Number, assetJSON)
	if err != nil {
		return err
	}
	return nil
}

// function to register the dispense of a medicament in a pharmacy
func (s *SmartContract) DispenseMedicament(ctx contractapi.TransactionContextInterface, _entityID string, args []string) error {
	//obtain medicament that is dispensed
	_Serial_Number := args[0]
	medicament, err := s.ReadMedicament(ctx, _Serial_Number)
	if err != nil {
		return err
	}

	//check current owner
	if medicament.Current_Owner != _entityID {
		return fmt.Errorf("This medicament does not belong to this pharmacy")
	}

	//check dispenser entity
	if medicament.Seller_Pharmacy != _entityID {
		return fmt.Errorf("This medicament is not addressed to this pharmacy")
	}

	//check that medicament is not expired
	expired, err := s.isExpired(ctx, medicament.Expiration_Year, medicament.Expiration_Month)
	if err != nil {
		return err
	}
	if expired {
		return fmt.Errorf("This medicament is expired, it can not be processed")
	}

	//modify medicament status
	status, err := s.getNewStatus(ctx, medicament.Status, "DispenseMedicament")
	if err != nil {
		return err
	}

	medicament.Status = status

	//update dates
	MedDates, err := s.UpdateDates(ctx, "DispenseMedicament", medicament.Dates)
	if err != nil {
		return err
	}
	medicament.Dates = MedDates

	//update medicament on ledger
	assetJSON, err := json.Marshal(medicament)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(_Serial_Number, assetJSON)
	if err != nil {
		return err
	}
	return nil
}

// function that validates if the password of a given user corresponds to the passed password
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

// function that creates a new session for a logged user
func (s *SmartContract) CreateSession(ctx contractapi.TransactionContextInterface, _sessionID string, _entityID string, _userID string) error {
	actualDateStr, err := s.GetTxTimestamp(ctx)
	if err != nil {
		return err
	}
	session := Session{
		EntityID:       _entityID,
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

// function that stores a session of a user in the user log
func (s *SmartContract) StoreSession(ctx contractapi.TransactionContextInterface, _sessionID string, _userID string) error {
	user, err := s.ReadUser(ctx, _userID)
	if err != nil {
		return err
	}
	user.Sessions_Log = append(user.Sessions_Log, _sessionID)
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(_userID, userJSON)
	if err != nil {
		return err
	}
	return nil
}

// function that checks if a session is expired and modifies its state in that case
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

// function that returns the current counter for sessions
func (s *SmartContract) GetCounter(ctx contractapi.TransactionContextInterface) int {
	counterJSON, _ := ctx.GetStub().GetState("SessionCounter")

	var counter Counter
	json.Unmarshal(counterJSON, &counter)

	return counter.Count
}

// function that increments the counter for sessions
func (s *SmartContract) IncrementCounter(ctx contractapi.TransactionContextInterface, _counter int) {
	_count := _counter + 1
	counter := Counter{Count: _count}
	counterJSON, _ := json.Marshal(counter)

	_ = ctx.GetStub().PutState("SessionCounter", counterJSON)
}

// function that generates a unique session ID
func (s *SmartContract) GenerateSessionID(ctx contractapi.TransactionContextInterface) string {

	counter := s.GetCounter(ctx)
	s.IncrementCounter(ctx, counter)

	return strconv.Itoa(counter)
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
		if len(_allArgs) != 2 {
			return false, fmt.Errorf("Incorrect number of arguments")
		}
		if len(_allArgs[0]) == 0 {
			return false, fmt.Errorf("Recipient entity ID must be provided")
		}
		if len(_allArgs[0]) == 1 {
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
		if _function == "LogIn" {
			return true, nil
		}
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
		return nil, fmt.Errorf("the medicament %s does not exist", _Serial_Number)
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
		return nil, fmt.Errorf("the entity %s does not exist", _Entity_ID)
	}

	var entity Entity
	err = json.Unmarshal(entityJSON, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// function that returns all the active sessions of a user in an entity, should always be 1 or 0
func (s *SmartContract) GetActiveSessions(ctx contractapi.TransactionContextInterface, _Entity_ID string, _User_ID string) ([]*Session, error) {
	user, err := s.ReadUser(ctx, _User_ID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		userLog := user.Sessions_Log
		var userSessions []*Session
		for _, userSessionID := range userLog {
			session, err := s.ReadSession(ctx, userSessionID)
			if err != nil {
				return nil, err
			}
			if session.EntityID == _Entity_ID && session.Status == 1 {
				userSessions = append(userSessions, session)
			}
		}
		return userSessions, nil
	} else {
		return nil, fmt.Errorf("Invalid User")
	}
}

// function used to log out from all the sessions of a user in an entity
func (s *SmartContract) LogOut(ctx contractapi.TransactionContextInterface, _Entity_ID string, _User_ID string) error {
	user, err := s.ReadUser(ctx, _User_ID)
	if err != nil {
		return err
	}
	if user != nil {
		userLog := user.Sessions_Log

		for _, userSessionID := range userLog {
			session, err := s.ReadSession(ctx, userSessionID)
			if err != nil {
				return err
			}
			if session.EntityID == _Entity_ID && session.Status == 1 {
				session.Status = 0

				assetJSON, err := json.Marshal(session)
				if err != nil {
					return err
				}

				err = ctx.GetStub().PutState(userSessionID, assetJSON)
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return fmt.Errorf("Invalid User")
	}
}

// function used to obtain all the sessions of a user in an entity
func (s *SmartContract) GetSessions(ctx contractapi.TransactionContextInterface, _Entity_ID string, _User_ID string) ([]*Session, error) {
	user, err := s.ReadUser(ctx, _User_ID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		userLog := user.Sessions_Log
		var userSessions []*Session
		for _, userSessionID := range userLog {
			session, err := s.ReadSession(ctx, userSessionID)
			if err != nil {
				return nil, err
			}
			if session.EntityID == _Entity_ID {
				userSessions = append(userSessions, session)
			}
		}
		return userSessions, nil
	} else {
		return nil, fmt.Errorf("Invalid User")
	}
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

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) MedicamentExists(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(_Serial_Number)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// function that returns the Tx time stamp
func (s *SmartContract) GetTxTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	txTimeAsPtr, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", err
	}
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()

	return timeStr, nil
}

// function that updates the registered dates of a medicament depending on the tracking point
func (s *SmartContract) UpdateDates(ctx contractapi.TransactionContextInterface, function string, _currentDates MedicamentDates) (MedicamentDates, error) {
	TxDate, err := s.GetTxTimestamp(ctx)
	if err != nil {
		return _currentDates, err
	}
	if function == "RegisterMedicament" {
		_currentDates.RegistrationDate = TxDate
	} else if function == "DispatchMedicament" {
		_currentDates.DispatchDate = TxDate
	} else if function == "ReceiveMedicament" {
		_currentDates.ReceiveDate = TxDate
	} else if function == "DispenseMedicament" {
		_currentDates.DispenseDate = TxDate
	} else {
		return _currentDates, fmt.Errorf("Undefined function, medicament dates have not been updated")
	}
	return _currentDates, nil
}

// function that gets a registered medicament given a serial number and returns it after checking the current owner of the medicament
func (s *SmartContract) GetMedicament(ctx contractapi.TransactionContextInterface, _userID string, _entityID string, _Serial_Number string) (*Medicament, error) {
	activeSessions, err := s.GetActiveSessions(ctx, _entityID, _userID)
	if err != nil {
		return nil, err
	}
	if len(activeSessions) == 1 {
		session := activeSessions[0]
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
			userInEntity, err := s.isUserInEntity(ctx, _userID, users)

			if err != nil {
				return nil, err
			}
			if !userInEntity {
				return nil, fmt.Errorf("Invalid user")
			}
			user, err := s.ReadUser(ctx, _userID)
			if err != nil {
				return nil, err
			}
			if user.Rol == "admin" {
				medicament, err := s.ReadMedicament(ctx, _Serial_Number)
				if err != nil {
					return nil, err
				}
				if medicament.Current_Owner != _entityID {
					return nil, fmt.Errorf("Can not access to medicament without being the owner")
				}
				return medicament, nil
			} else {
				return nil, fmt.Errorf("Can not access to a medicament without being the admin")
			}
		} else {
			return nil, fmt.Errorf("Session expired. Log in again, please")
		}
	} else if len(activeSessions) == 0 {
		return nil, fmt.Errorf("No active session for that user. Please, log in")

	} else {
		s.LogOut(ctx, _entityID, _userID)
		return nil, fmt.Errorf("No active session for that user. Please, log in")
	}
}

// function that returns all registered medicaments owner by an entity in the system
func (s *SmartContract) GetAllMedicaments(ctx contractapi.TransactionContextInterface, _userID string, _entityID string) ([]*Medicament, error) { //
	activeSessions, err := s.GetActiveSessions(ctx, _entityID, _userID)
	if err != nil {
		return nil, err
	}
	if len(activeSessions) == 1 {
		session := activeSessions[0]
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
			userInEntity, err := s.isUserInEntity(ctx, _userID, users)
			if err != nil {
				return nil, err
			}
			if !userInEntity {
				return nil, fmt.Errorf("Invalid user")
			}
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
					if medicament.Current_Owner == _entityID {
						medicaments = append(medicaments, &medicament)
					}
				}

				return medicaments, nil
			} else {
				return nil, fmt.Errorf("Can not access to medicaments without being the admin")
			}
		} else {
			return nil, fmt.Errorf("Session expired. Log in again, please")
		}
	} else if len(activeSessions) == 0 {
		return nil, fmt.Errorf("No active session for that user. Please, log in")

	} else {
		s.LogOut(ctx, _entityID, _userID)
		return nil, fmt.Errorf("No active session for that user. Please, log in")
	}
}

// function that returns all registered users of an entity in the system
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface, _userID string, _entityID string) ([]*Entity_User, error) {
	activeSessions, err := s.GetActiveSessions(ctx, _entityID, _userID)
	if err != nil {
		return nil, err
	}
	if len(activeSessions) == 1 {
		session := activeSessions[0]
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
			userInEntity, err := s.isUserInEntity(ctx, _userID, users)
			if err != nil {
				return nil, err
			}
			if !userInEntity {
				return nil, fmt.Errorf("Invalid user")
			}

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
	} else if len(activeSessions) == 0 {
		return nil, fmt.Errorf("No active session for that user. Please, log in")

	} else {
		s.LogOut(ctx, _entityID, _userID)
		return nil, fmt.Errorf("No active session for that user. Please, log in")
	}
}
