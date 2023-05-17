package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	ID         string `json:"ID"`
	Owner      string `json:"Owner"`
	Status     string `json:"Status"`
	Type       string `json:"Type"`
	Department string `json:"Department"`
	Code       string `json:"Code"`
	Value      int    `json:"Value"`
	Date       string `json:"Date"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "Склад:М2", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "М2", Value: 60, Date: "02.04.2023 09:01"},
		{ID: "Склад:М3", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "М3", Value: 130, Date: "02.04.2023 09:04"},
		{ID: "Склад:М2а", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "М2а", Value: 20, Date: "02.04.2023 09:08"},
		{ID: "Склад:М4", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "М4", Value: 80, Date: "02.04.2023 09:09"},
		{ID: "Склад:М8", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "М8", Value: 190, Date: "02.04.2023 09:15"},
		{ID: "Склад:А2", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "А2", Value: 340, Date: "02.04.2023 09:32"},
		{ID: "Склад:А19", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "А19", Value: 540, Date: "02.04.2023 09:45"},
		{ID: "Склад:Бр1", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "Бр1", Value: 10, Date: "02.04.2023 09:52"},
		{ID: "Склад:Л8", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "Л8", Value: 35, Date: "02.04.2023 09:55"},
		{ID: "Склад:Ц7", Owner: "Склад", Status: "Готов к продаже", Type: "Хранение", Department: "Склад", Code: "Ц7", Value: 40, Date: "02.04.2023 09:58"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, Owner string, Status string, Type string, Department string, Code string, Value int, Date string) error {
	id := Owner + ":" + Code
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	if Type == "Сделка" {
		if Status != "Первичная договоренность" {
			return fmt.Errorf("статус актива %s не соответсвует правилам", id)
		}
	}

	asset := Asset{
		ID:         id,
		Owner:      Owner,
		Status:     Status,
		Type:       Type,
		Department: Department,
		Code:       Code,
		Value:      Value,
		Date:       Date,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, Owner string, Status string, Type string, Department string, Code string, Value int, Date string) error {
	id := Owner + ":" + Code
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:         id,
		Owner:      Owner,
		Status:     Status,
		Type:       Type,
		Department: Department,
		Code:       Code,
		Value:      Value,
		Date:       Date,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
