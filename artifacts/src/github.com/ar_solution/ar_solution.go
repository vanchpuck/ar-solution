package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Document struct {
	DocNum       string `json:"docNum"`
	InDocNum     string `json:"inDocNum"`
	Submitter    string `json:"submitter"`
	DocType      string `json:"docType"`
	Date         string `json:"date"`
	InDocDate    string `json:"inDocDate"`
	Sender       string `json:"sender"`
	Recepient    string `json:"recepient"`
	AmountDebit  string `json:"amountDebit"`
	AmountCredit string `json:"amountCredit"`
	Description  string `json:"description"`
	Status       string `json:"status"`
}

type Selector interface {
	SelectorStr() string
}

type SelectorBuilder interface {
	DocType(DocumentType) SelectorBuilder
	Status(string) SelectorBuilder
	Submitter(string) SelectorBuilder
	Sender(string) SelectorBuilder
	Recepient(string) SelectorBuilder
	Build() Selector
}

type selectorBuilder struct {
	docType   DocumentType
	status    string
	submitter string
	sender    string
	recepient string
}

func (sb *selectorBuilder) DocType(docType DocumentType) SelectorBuilder {
	sb.docType = docType
	return sb
}

func (sb *selectorBuilder) Status(status string) SelectorBuilder {
	sb.status = status
	return sb
}

func (sb *selectorBuilder) Submitter(submitter string) SelectorBuilder {
	sb.submitter = submitter
	return sb
}

func (sb *selectorBuilder) Sender(sender string) SelectorBuilder {
	sb.sender = sender
	return sb
}

func (sb *selectorBuilder) Recepient(recepient string) SelectorBuilder {
	sb.recepient = recepient
	return sb
}

func (sb *selectorBuilder) Build() Selector {
	return &selector{
		docType:   sb.docType,
		status:    sb.status,
		submitter: sb.submitter,
		sender:    sb.sender,
		recepient: sb.recepient,
	}
}

type selector struct {
	docType   DocumentType
	status    string
	submitter string
	sender    string
	recepient string
}

func (s *selector) SelectorStr() string {
	var fields []string
	if len(s.docType) > 0 {
		fields = append(fields, "\"docType\":\""+string(s.docType)+"\"")
	}
	if len(s.status) > 0 {
		fields = append(fields, "\"status\":\""+s.status+"\"")
	}
	if len(s.submitter) > 0 {
		fields = append(fields, "\"submitter\":\""+s.submitter+"\"")
	}
	if len(s.sender) > 0 {
		fields = append(fields, "\"sender\":\""+s.sender+"\"")
	}
	if len(s.recepient) > 0 {
		fields = append(fields, "\"recepient\":\""+s.recepient+"\"")
	}
	return "{\"selector\":{" + strings.Join(fields, ",") + "}}"
}

type DocumentType string

const (
	Purchase         DocumentType = "PURCHASE"
	Sale             DocumentType = "SALE"
	Expense          DocumentType = "EXPENSE"
	Admission        DocumentType = "ADMISSION"
	AdmissionCorrect DocumentType = "ADMISSION_CORRECT"
	SaleCorrect      DocumentType = "SALE_CORRECT"
)

// func (qf QueryFilter) String() string {
// 	return fmt.Sprintf("\"%v\":\"(%v)\"", qf.Key, qf.Value)
// }

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

	if function == "newPurchaseDoc" {
		return t.newPurchaseDoc(stub, args)
	}
	if function == "newSaleDoc" {
		return t.newSaleDoc(stub, args)
	}
	if function == "newExpenseDoc" {
		return t.newExpenseDoc(stub, args)
	}
	if function == "newAdmissionDoc" {
		return t.newAdmissionDoc(stub, args)
	}
	if function == "newAdmissionCorrectDoc" {
		return t.newAdmissionCorrectDoc(stub, args)
	}
	if function == "newSaleCorrectDoc" {
		return t.newSaleCorrectDoc(stub, args)
	}
	if function == "query" {
		return t.query(stub, args)
	}
	if function == "getAllPurchaseDocs" {
		return t.getAllPurchaseDocs(stub, args)
	}
	if function == "getPurchaseDocsBySender" {
		return t.getPurchaseDocsBySender(stub, args)
	}
	if function == "getPurchaseDocsByRecepient" {
		return t.getPurchaseDocsByRecepient(stub, args)
	}
	if function == "getPurchaseDocsBySubmitter" {
		return t.getPurchaseDocsBySubmitter(stub, args)
	}
	if function == "getAllExpenseDocs" {
		return t.getAllExpenseDocs(stub, args)
	}
	if function == "getExpenseDocsBySender" {
		return t.getExpenseDocsBySender(stub, args)
	}
	if function == "getExpenseDocsByRecepient" {
		return t.getExpenseDocsByRecepient(stub, args)
	}
	if function == "getExpenseDocsBySubmitter" {
		return t.getExpenseDocsBySubmitter(stub, args)
	}
	if function == "getAllSaleDocs" {
		return t.getAllSaleDocs(stub, args)
	}
	if function == "getSaleDocsBySender" {
		return t.getSaleDocsBySender(stub, args)
	}
	if function == "getSaleDocsByRecepient" {
		return t.getSaleDocsByRecepient(stub, args)
	}
	if function == "getSaleDocsBySubmitter" {
		return t.getSaleDocsBySubmitter(stub, args)
	}
	if function == "getAllAdmissionDocs" {
		return t.getAllAdmissionDocs(stub, args)
	}
	if function == "getAdmissionDocsBySender" {
		return t.getAdmissionDocsBySender(stub, args)
	}
	if function == "getAdmissionDocsByRecepient" {
		return t.getAdmissionDocsByRecepient(stub, args)
	}
	if function == "getAdmissionDocsBySubmitter" {
		return t.getAdmissionDocsBySubmitter(stub, args)
	}
	if function == "getAllAdmissionCorrectDocs" {
		return t.getAllAdmissionCorrectDocs(stub, args)
	}
	if function == "getAdmissionCorrectDocsBySender" {
		return t.getAdmissionCorrectDocsBySender(stub, args)
	}
	if function == "getAdmissionCorrectDocsByRecepient" {
		return t.getAdmissionCorrectDocsByRecepient(stub, args)
	}
	if function == "getAdmissionCorrectDocsBySubmitter" {
		return t.getAdmissionCorrectDocsBySubmitter(stub, args)
	}
	if function == "getAllSaleCorrectDocs" {
		return t.getAllSaleCorrectDocs(stub, args)
	}
	if function == "getSaleCorrectDocsBySender" {
		return t.getSaleCorrectDocsBySender(stub, args)
	}
	if function == "getSaleCorrectDocsByRecepient" {
		return t.getSaleCorrectDocsByRecepient(stub, args)
	}
	if function == "getSaleCorrectDocsBySubmitter" {
		return t.getSaleCorrectDocsBySubmitter(stub, args)
	}
	if function == "cancelDoc" {
		return t.cancelDoc(stub, args)
	}

	// const errMsg = "Unknown action, must be one of 'doPayment' but got: %v", args[0]
	var errMsg = fmt.Sprintf("Unknown action. Got: %v", args[0])
	logger.Errorf(errMsg)
	return shim.Error(errMsg)
}

func (cc *Chaincode) newPurchaseDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return checkArgsAndCreateDoc(stub, args, Purchase)
}

func (s *Chaincode) newSaleDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return checkArgsAndCreateDoc(stub, args, Sale)
}

func (s *Chaincode) newExpenseDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return checkArgsAndCreateDoc(stub, args, Expense)
}

func (s *Chaincode) newAdmissionDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return checkArgsAndCreateDoc(stub, args, Admission)
}

func (s *Chaincode) newAdmissionCorrectDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return checkArgsAndCreateDoc(stub, args, AdmissionCorrect)
}

func (s *Chaincode) newSaleCorrectDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return checkArgsAndCreateDoc(stub, args, SaleCorrect)
}

// func (s *Chaincode) doPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	if len(args) != 6 {
// 		return shim.Error("Incorrect number of arguments. Expecting 8")
// 	}
// 	return addDoc(stub, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], "PAYMENT")
// }

// func (s *Chaincode) doSupply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	if len(args) != 6 {
// 		return shim.Error("Incorrect number of arguments. Expecting 8")
// 	}
// 	return addDoc(stub, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], "SUPPLY")
// }

func checkArgsAndCreateDoc(stub shim.ChaincodeStubInterface, args []string, docType DocumentType) pb.Response {
	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments.")
	}
	return newDoc(stub, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], docType)
}

// TODO amount should be int
func newDoc(stub shim.ChaincodeStubInterface, docId string, docNum string, submitter string, docDate string,
	sender string, recepient string, amountDebit string, amountCredit string, description string, inDocNum string, inDocDate string, docType DocumentType) pb.Response {

	// TODO Check whether document already exist
	var doc = Document{
		DocNum:       docNum,
		InDocNum:     inDocNum,
		InDocDate:    inDocDate,
		Submitter:    submitter,
		Date:         docDate,
		Sender:       sender,
		Recepient:    recepient,
		AmountDebit:  amountDebit,
		AmountCredit: amountCredit,
		DocType:      string(docType),
		Description:  description,
		Status:       "COMMITTED"}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(submitter+"-"+docId, docBytes)

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

func (cc *Chaincode) getAllSaleDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	return cc.getDocs(stub, Sale, "", "", "")
}

func (cc *Chaincode) getSaleDocsBySender(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Sale, "", args[0], "")
}

func (cc *Chaincode) getSaleDocsByRecepient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Sale, "", "", args[0])
}

func (cc *Chaincode) getSaleDocsBySubmitter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Sale, args[0], "", "")
}

func (cc *Chaincode) getAllPurchaseDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	return cc.getDocs(stub, Purchase, "", "", "")
}

func (cc *Chaincode) getPurchaseDocsBySender(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Purchase, "", args[0], "")
}

func (cc *Chaincode) getPurchaseDocsByRecepient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Purchase, "", "", args[0])
}

func (cc *Chaincode) getPurchaseDocsBySubmitter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Purchase, args[0], "", "")
}

func (cc *Chaincode) getAllExpenseDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	return cc.getDocs(stub, Expense, "", "", "")
}

func (cc *Chaincode) getExpenseDocsBySender(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Expense, "", args[0], "")
}

func (cc *Chaincode) getExpenseDocsByRecepient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Expense, "", "", args[0])
}

func (cc *Chaincode) getExpenseDocsBySubmitter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Expense, args[0], "", "")
}

func (cc *Chaincode) getAllAdmissionDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	return cc.getDocs(stub, Admission, "", "", "")
}

func (cc *Chaincode) getAdmissionDocsBySender(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Admission, "", args[0], "")
}

func (cc *Chaincode) getAdmissionDocsByRecepient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Admission, "", "", args[0])
}

func (cc *Chaincode) getAdmissionDocsBySubmitter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, Admission, args[0], "", "")
}

func (cc *Chaincode) getAllAdmissionCorrectDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	return cc.getDocs(stub, AdmissionCorrect, "", "", "")
}

func (cc *Chaincode) getAdmissionCorrectDocsBySender(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, AdmissionCorrect, "", args[0], "")
}

func (cc *Chaincode) getAdmissionCorrectDocsByRecepient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, AdmissionCorrect, "", "", args[0])
}

func (cc *Chaincode) getAdmissionCorrectDocsBySubmitter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, AdmissionCorrect, args[0], "", "")
}

func (cc *Chaincode) getAllSaleCorrectDocs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	return cc.getDocs(stub, SaleCorrect, "", "", "")
}

func (cc *Chaincode) getSaleCorrectDocsBySender(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, SaleCorrect, "", args[0], "")
}

func (cc *Chaincode) getSaleCorrectDocsByRecepient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, SaleCorrect, "", "", args[0])
}

func (cc *Chaincode) getSaleCorrectDocsBySubmitter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	return cc.getDocs(stub, SaleCorrect, args[0], "", "")
}

func (t *Chaincode) getDocs(stub shim.ChaincodeStubInterface, docType DocumentType, submitter string, sender string, recepient string) pb.Response {
	builder := &selectorBuilder{}
	if len(docType) > 0 {
		logger.Info("DocumentType: " + string(docType))
		builder.DocType(docType)
	}
	if len(submitter) > 0 {
		logger.Info("Submitter: " + submitter)
		builder.Submitter(submitter)
	}
	if len(sender) > 0 {
		logger.Info("Sender: " + sender)
		builder.Sender(sender)
	}
	if len(recepient) > 0 {
		logger.Info("Recepient: " + recepient)
		builder.Recepient(recepient)
	}
	return queryDatabase(stub, builder.Build().SelectorStr())
}

func getDocsByTypeAndSender(stub shim.ChaincodeStubInterface, docType DocumentType, sender string) pb.Response {
	return queryDatabase(stub, fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"sender\":\"%s\"}}", string(docType), sender))
}

func queryDatabase(stub shim.ChaincodeStubInterface, queryString string) pb.Response {
	// TODO queryString should be sanitized
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *Chaincode) cancelDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var submitter string = args[0]
	var docNumber string = args[1]
	var key string = key(submitter, docNumber)

	docBytes, _ := stub.GetState(key)
	doc := Document{}

	json.Unmarshal(docBytes, &doc)
	doc.Status = "CANCELED"

	docBytes, _ = json.Marshal(doc)
	stub.PutState(key, docBytes)

	return shim.Success(nil)
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

func key(submitter string, docNumber string) string {
	return submitter + "-" + docNumber
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
