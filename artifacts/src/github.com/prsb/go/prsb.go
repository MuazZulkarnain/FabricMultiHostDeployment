package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Token :  Define the token structure, with 4 properties.  Structure tags are used by encoding/json library
type Token struct {
	Amount         float64 `json:"amount"`
	Owner          string  `json:"owner"`
	Source         string  `json:"source"`
	ConversionRate float64 `json:"conversion_rate"`
	PastOperation  string  `json:"past_operation"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("prsb_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %s", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryToken":
		return s.queryToken(APIstub, args)
	case "queryTokenByTxID":
		return s.queryTokenByTxID(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createToken":
		return s.createToken(APIstub, args)
	case "updateTokenVolume":
		return s.updateTokenVolume(APIstub, args)
	case "changeTokenOwner":
		return s.changeTokenOwner(APIstub, args)
	case "retireToken":
		return s.retireToken(APIstub, args)
	case "queryAllTokens":
		return s.queryAllTokens(APIstub)
	case "queryTokenHistory":
		return s.queryTokenHistory(APIstub, args)
	case "queryTokenHistoryByTxID":
		return s.queryTokenHistoryByTxID(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

	// return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	tokens := []Token{
		{Amount: 2.31, Owner: "PRSB-A", Source: "PRSB-A", ConversionRate: 0.6689},
		{Amount: 19.87, Owner: "PRSB-B", Source: "PRSB-B", ConversionRate: 0.6689},
		{Amount: 4.11, Owner: "PRSB-C", Source: "PRSB-C", ConversionRate: 0.6689},
		{Amount: 7.49, Owner: "PRSB-D", Source: "PRSB-D", ConversionRate: 0.6689},
	}

	i := 0
	for i < len(tokens) {
		tokenAsBytes, _ := json.Marshal(tokens[i])
		APIstub.PutState("TOKEN"+strconv.Itoa(i), tokenAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createToken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// Check the number of arguments
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// Parse the amount argument
	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Invalid amount: " + err.Error())
	}

	// Parse the conversion rate argument
	conversion_rate, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return shim.Error("Invalid conversion rate: " + err.Error())
	}

	// Check if the asset already exists
	tokenAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if tokenAsBytes != nil {
		return shim.Error("Asset already exists: " + args[0])
	}

	// Create the token object
	var token = Token{Amount: amount, Owner: args[2], Source: args[3], ConversionRate: conversion_rate, PastOperation: args[5]}

	// Marshal the token object to bytes
	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error("Failed to marshal token: " + err.Error())
	}

	// Put the token in the ledger
	err = APIstub.PutState(args[0], tokenAsBytes)
	if err != nil {
		return shim.Error("Failed to put asset: " + err.Error())
	}

	// Create the composite key for the owner
	indexName := "owner~key"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{token.Owner, args[0]})
	if err != nil {
		return shim.Error("Failed to create composite key: " + err.Error())
	}

	// Put the index in the ledger
	err = APIstub.PutState(colorNameIndexKey, []byte{0x00})
	if err != nil {
		return shim.Error("Failed to put index: " + err.Error())
	}

	// Get the transaction ID (txID)
	txID := APIstub.GetTxID()

	// Build the response payload
	responsePayload := struct {
		Token Token  `json:"token"`
		TxID  string `json:"txID"`
	}{
		Token: token,
		TxID:  txID,
	}

	responsePayloadAsBytes, err := json.Marshal(responsePayload)
	if err != nil {
		return shim.Error("Failed to marshal response payload: " + err.Error())
	}

	// Return the response payload as success response
	return shim.Success(responsePayloadAsBytes)
}

func (s *SmartContract) updateTokenVolume(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	tokenAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get token: " + err.Error())
	}

	token := Token{}
	err = json.Unmarshal(tokenAsBytes, &token)
	if err != nil {
		return shim.Error("Failed to unmarshal token: " + err.Error())
	}

	newVolume, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Invalid volume: " + err.Error())
	}

	token.PastOperation = args[2]
	token.Amount = newVolume
	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error("Failed to marshal token: " + err.Error())
	}

	err = APIstub.PutState(args[0], tokenAsBytes)
	if err != nil {
		return shim.Error("Failed to update token: " + err.Error())
	}

	// Get the transaction ID (txID)
	txID := APIstub.GetTxID()

	// Build the response payload
	responsePayload := struct {
		Token Token  `json:"token"`
		TxID  string `json:"txID"`
	}{
		Token: token,
		TxID:  txID,
	}

	responsePayloadAsBytes, err := json.Marshal(responsePayload)
	if err != nil {
		return shim.Error("Failed to marshal response: " + err.Error())
	}

	return shim.Success(responsePayloadAsBytes)
}

func (s *SmartContract) changeTokenOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	tokenAsBytes, _ := APIstub.GetState(args[0])
	token := Token{}

	json.Unmarshal(tokenAsBytes, &token)
	token.Owner = args[1]

	tokenAsBytes, _ = json.Marshal(token)
	APIstub.PutState(args[0], tokenAsBytes)

	// Get the transaction ID (txID)
	txID := APIstub.GetTxID()

	// Build the response payload
	responsePayload := struct {
		Token Token  `json:"token"`
		TxID  string `json:"txID"`
	}{
		Token: token,
		TxID:  txID,
	}

	responsePayloadAsBytes, _ := json.Marshal(responsePayload)

	return shim.Success(responsePayloadAsBytes)
}

func (s *SmartContract) retireToken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err := APIstub.DelState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to retire token with ID %s: %s", args[0], err.Error()))
	}

	// Get the transaction ID (txID)
	txID := APIstub.GetTxID()

	// Build the response payload
	responsePayload := struct {
		TokenID string `json:"tokenID"`
		TxID    string `json:"txID"`
	}{
		TokenID: args[0],
		TxID:    txID,
	}

	responsePayloadAsBytes, _ := json.Marshal(responsePayload)

	return shim.Success(responsePayloadAsBytes)
}

func (s *SmartContract) queryToken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tokenAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(tokenAsBytes)
}

func (s *SmartContract) queryTokenByTxID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	txID := args[0]

	// Get the history of all transactions with the given txID as the key
	resultsIterator, err := APIstub.GetHistoryForKey(txID)
	if err != nil {
		return shim.Error("Failed to get transaction history: " + err.Error())
	}
	defer resultsIterator.Close()

	// Iterate through the transaction history and look for the token with the given txID
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Failed to get next transaction: " + err.Error())
		}

		// Unmarshal the transaction value to a Token object
		var token Token
		err = json.Unmarshal(queryResponse.Value, &token)
		if err != nil {
			return shim.Error("Failed to unmarshal transaction value to Token: " + err.Error())
		}

		// Check if the token has the given txID
		if queryResponse.TxId == txID {
			// Build the response payload
			responsePayload := struct {
				Token Token  `json:"token"`
				TxID  string `json:"txID"`
			}{
				Token: token,
				TxID:  queryResponse.TxId,
			}

			responsePayloadAsBytes, err := json.Marshal(responsePayload)
			if err != nil {
				return shim.Error("Failed to marshal response payload: " + err.Error())
			}

			// Return the response payload as success response
			return shim.Success(responsePayloadAsBytes)
		}
	}

	// If the function has not yet returned, then the token was not found with the given txID
	return shim.Error("Token not found with txID: " + txID)
}

func (s *SmartContract) queryAllTokens(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "TOKEN0"
	endKey := "TOKEN99999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTokens:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (t *SmartContract) queryTokenHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tokenName := args[0]

	resultsIterator, err := stub.GetHistoryForKey(tokenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryTokenHistoryByTxID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Check the number of arguments
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Get the history iterator
	resultsIterator, err := APIstub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error("Failed to get history for asset: " + err.Error())
	}
	defer resultsIterator.Close()

	// Convert the history to a slice of Token objects
	var history []Token
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Failed to get next history entry: " + err.Error())
		}

		// Unmarshal the transaction value to a Token object
		var token Token
		err = json.Unmarshal(response.Value, &token)
		if err != nil {
			return shim.Error("Failed to unmarshal transaction value: " + err.Error())
		}

		// Add the Token object to the history slice
		history = append(history, token)
	}

	// Convert the history slice to JSON bytes
	historyAsBytes, err := json.Marshal(history)
	if err != nil {
		return shim.Error("Failed to marshal history: " + err.Error())
	}

	// Return the history as a success response
	return shim.Success(historyAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
