package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Document struct {
	DocNum    string `json:"docNum"`
	Submitter string `json:"submitter"`
	DocType   string `json:"docType"`
	Date      string `json:"date"`
	Sender    string `json:"sender"`
	Recepient string `json:"recepient"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
}

var logger = shim.NewLogger("ar_solution")

// SimpleChaincode example simple Chaincode implementation
type Chaincode struct {
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### ar_solution Init ###########")
	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### ar_solution Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "doPayment" {
		return t.doPayment(stub, args)
	}
	if function == "doSupply" {
		return t.doSupply(stub, args)
	}
	if function == "query" {
		return t.query(stub, args)
	}
	if function == "queryAllBySender" {
		return t.queryAllDocsBySender(stub, args)
	}
	if function == "queryAllByRecepient" {
		return t.queryAllDocsByRecepient(stub, args)
	}
	if function == "queryAllPaymentDocs" {
		return t.queryAllPaymentDocs(stub, args)
	}
	if function == "queryAllSupplyDocs" {
		return t.queryAllSupplyDocs(stub, args)
	}

	// const errMsg = "Unknown action, must be one of 'doPayment' but got: %v", args[0]
	var errMsg = fmt.Sprintf("Unknown action. Got: %v", args[0])
	logger.Errorf(errMsg)
	return shim.Error(errMsg)
}

func (s *Chaincode) doPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	return addDoc(stub, args[0], args[1], args[2], args[3], args[4], args[5], "PAYMENT")
}

func (s *Chaincode) doSupply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	return addDoc(stub, args[0], args[1], args[2], args[3], args[4], args[5], "SUPPLY")
}

// amount should be int
func addDoc(stub shim.ChaincodeStubInterface, docNum string, submitter string, docDate string,
	sender string, recepient string, amount string, docType string) pb.Response {

	// Check whether document already exist
	var doc = Document{
		DocNum:    docNum,
		Submitter: submitter,
		Date:      docDate,
		Sender:    sender,
		Recepient: recepient,
		Amount:    amount,
		DocType:   docType,
		Status:    "COMMITTED"}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(submitter+"-"+docNum, docBytes)

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

func (t *Chaincode) queryAllDocsBySender(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// sender := strings.ToLower(args[0])
	logger.Info("Argument: " + args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"sender\":\"%s\"}}", args[0])

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (t *Chaincode) queryAllDocsByRecepient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// sender := strings.ToLower(args[0])
	logger.Info("Argument: " + args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"recepient\":\"%s\"}}", args[0])

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (t *Chaincode) queryAllPaymentDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return getDocsByType(stub, "PAYMENT")
}

func (t *Chaincode) queryAllSupplyDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return getDocsByType(stub, "SUPPLY")
}

func getDocsByType(stub shim.ChaincodeStubInterface, docType string) pb.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", docType)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		logger.Info("################# Query response ####################")
		logger.Info(queryResponse)

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

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
