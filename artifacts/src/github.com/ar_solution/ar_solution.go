package main

import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Document struct {
	DocNum	string `json:"docNum"`
	Submitter	string `json:"submitter"`
	DocType  string `json:"docType"`
	Date string `json:"date"`
	Sender  string `json:"sender"`
	Recepient  string `json:"recepient"`
	Amount  string `json:"amount"`
}

var logger = shim.NewLogger("ar_solution")

// SimpleChaincode example simple Chaincode implementation
type Chaincode struct {
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	logger.Info("########### ar_solution Init ###########")
	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### ar_solution Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	
	if function == "doPayment" {
		return t.doPayment(stub, args)
	}
	if function == "query" {
		return t.query(stub, args)
	}

	// const errMsg = "Unknown action, must be one of 'doPayment' but got: %v", args[0]
	var errMsg = fmt.Sprintf("Unknown action, must be one of 'doPayment', 'query' but got: %v", args[0])
	logger.Errorf(errMsg)
	return shim.Error(errMsg)
}

func (s *Chaincode) doPayment(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var doc = Document{DocNum: args[0], Submitter: args[1], DocType: "PAYMENT", Date: args[2], Sender: args[3], Recepient: args[4], Amount: args[5]}

	docBytes, _ := json.Marshal(doc)
	APIstub.PutState(args[1] + "-" + args[0], docBytes)

	return shim.Success(nil)
}

func (t *Chaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting document id to query")
	}

	var docId string = args[0]

	DocBytes, err := stub.GetState(docId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + docId + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(DocBytes)

	// // Get the state from the ledger
	// DocBytes, err := stub.GetState(A)
	// if err != nil {
	// 	jsonResp := "{\"Error\":\"Failed to get state for " + docId + "\"}"
	// 	return shim.Error(jsonResp)
	// }

	// if DocBytes == nil {
	// 	jsonResp := "{\"Error\":\"Nil amount for " + docId + "\"}"
	// 	return shim.Error(jsonResp)
	// }

	// jsonResp := "{\"DocId\":\"" + docId + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	// logger.Infof("Query Response:%s\n", jsonResp)
	// return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
