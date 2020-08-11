package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PropertyTransfer contract to show the property transfer transactions
type PropertyTransferSmartContract struct {
    contractapi.Contract
}

// Property describes basic details
type Property struct {
	ID	string `json:"id"`
	Name	string `json:"name"`
	Area	int `json:"area"`
	OwnerName	string `json:"ownerName"`
	Value	int `json:"value"`
}

// This function helps to Add new property
func (pc *PropertyTransferSmartContract) AddProperty(ctx contractapi.TransactionContextInterface, id string,  name string, area int, ownerName string, value int) error {
    propertyJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return fmt.Errorf("Failed to read the data from world state", err)
    }
	
    if propertyJSON != nil {
		return fmt.Errorf("the property %s already exists", id)
    }
	
	prop := Property{
		ID:            id,
		Name:          name,
		Area:          area,
		OwnerName:     ownerName,
		Value: 		   value,
	}
	
	propertyBytes, err := json.Marshal(prop)	
	if err != nil {
		return err
	}

    return ctx.GetStub().PutState(id, propertyBytes)
}

// This function returns all the existing properties 
func (pc *PropertyTransferSmartContract) QueryAllProperties(ctx contractapi.TransactionContextInterface) ([]*Property, error) {
	propertyIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer propertyIterator.Close()

	var properties []*Property
	for propertyIterator.HasNext() {
		propertyResponse, err := propertyIterator.Next()
		if err != nil {
			return nil, err
		}

		var property *Property
		err = json.Unmarshal(propertyResponse.Value, &property)
		if err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}

	return properties, nil
}


// This function helps to query the property by Id
func (pc *PropertyTransferSmartContract) QueryPropertyById(ctx contractapi.TransactionContextInterface, id string) (*Property, error) {
    propertyJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("Failed to read the data from world state", err)
    }
	
    if propertyJSON == nil {
		return nil, fmt.Errorf("the property %s does not exist", id)
    }
	
	var property *Property
	err = json.Unmarshal(propertyJSON, &property)
	
	if err != nil {
		return nil, err
	}
	return property, nil
}


// This functions helps to transfer the ownserhip of the property
func (pc *PropertyTransferSmartContract) TransferProperty(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	property, err := pc.QueryPropertyById(ctx, id)
	if err != nil {
		return err
	}

	property.OwnerName = newOwner
	propertyJSON, err := json.Marshal(property)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, propertyJSON)
	
}

func main() {
    propTransferSmartContract := new(PropertyTransferSmartContract)

    cc, err := contractapi.NewChaincode(propTransferSmartContract)

    if err != nil {
        panic(err.Error())
    }

    if err := cc.Start(); err != nil {
        panic(err.Error())
    }
}
