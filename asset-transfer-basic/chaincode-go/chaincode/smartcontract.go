package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Entity_User struct {
	Name     string `json:"Name"`
	ID       string `json:"UserID"`
	Email    string `json:"Email"`
	Type     string `json:"UserType"`
	Address  string `json:"Address"`
	Password string `json:"Password"`
}

type Entity struct {
	Name         string        `json:"Name"`
	ID           string        `json:"UserID"`
	Type         string        `json:"UserType"`
	Entity_Users []Entity_User `json:"Entity_Users"`
}

type Medicament struct {
	Expiration_Month int    `json:"Expiration_Month"`
	Expiration_Year  int    `json:"Expiration_Year"`
	Lot_Number       string `json:"LotNumber"`
	Name             string `json:"Name"`
	Product_Code     int    `json:"Product_Code"`
	Serial_Number    string `json:"Serial_Number"`
	Status           int    `json:"Status"` // 1: creado | 2: despachado de lab | 3: recibido por farmacia | 4: dispensado | 5: indispensable por motivo que sea
	Entity_Producer  Entity `json:"Entity_Producer"`
	Entity_Owner     Entity `json:"Entity_Owner"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	//añadir creacion de usuarios, uno para
	entity_Users := []Entity_User{
		{Name: "lab", ID: "lab", Email: "lab@pg.com", Type: "lab", Address: "bangalore", Password: "adminpw"},
		{Name: "pharma", ID: "pharma", Email: "pharma@pg.com", Type: "pharma", Address: "bangalore", Password: "adminpw"},
	}

	for _, entity_User := range entity_Users {
		assetJSON, err := json.Marshal(entity_User)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(entity_User.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	medicaments := []Medicament{
		{Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352687", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352688", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Name: "Ibuprofeno", Product_Code: 8470008722513, Serial_Number: "6874352689", Lot_Number: "L201JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352687", Lot_Number: "L101JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352688", Lot_Number: "L101JX32", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
		{Name: "Paracetamol", Product_Code: 8470006723459, Serial_Number: "7874352689", Lot_Number: "L201JX30", Expiration_Year: 2024, Expiration_Month: 04, Status: 1},
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

func (s *SmartContract) isExpired(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {
	medicament, err := s.ReadAsset(ctx, _Serial_Number)
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
	/* asset, err := s.ReadAsset(ctx, _Serial_Number)
	if err != nil {
		return false
	}
	if asset != nil {
		return false
	} */

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("An expired medicament can not be registered")
	}

	//SI TODO VA BIEN, MEDICAMENTO ES REGISTRADO
	medicament := Medicament{
		Name:             _Name,
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
	medicament, err := s.ReadAsset(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("An expired medicament can not be registered")
	}

	//SI TODO VA BIEN:

	//SE SETEA EL FUTURE OWNER

	//EL STATUS CAMBIA
	medicament.Status = 2

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
	medicament, err := s.ReadAsset(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("An expired medicament can not be registered")
	}

	//REVISAR QUE EL FUTURE OWNER DEL FARMACO CORRESPONDE A LA FARMACIA QUE EJECUTA ESTA FUNCION

	//SI TODO VA BIEN, EL STATUS CAMBIA
	medicament.Status = 3

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
	medicament, err := s.ReadAsset(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}

	//REVISAR QUE NO ESTÉ CADUCADO
	expired, err := s.isExpired(ctx, _Serial_Number)
	if err != nil {
		return false, err
	}
	if expired {
		return false, fmt.Errorf("An expired medicament can not be registered")
	}

	//REVISAR QUE EL CURRENT OWNER DEL FARMACO CORRESPONDE A LA FARMACIA QUE EJECUTA ESTA FUNCION

	//SI TODO VA BIEN, EL STATUS CAMBIA
	medicament.Status = 4

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
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, _Serial_Number string) (*Medicament, error) {
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

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, _Serial_Number string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(_Serial_Number)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (t *SmartContract) GetAllMedicaments(ctx contractapi.TransactionContextInterface) ([]*Medicament, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
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
