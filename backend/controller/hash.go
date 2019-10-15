package controller

import (
	"backend/models"
	"backend/utils"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

var (
	ethRPC string = "https://ropsten.infura.io/v3/2a3f078d3755444b8777a0204e5f694a"
)

//Payload for other functions
type Payload struct {
	JSONRPC string
	Method  string
	Params  []interface{}
	ID      int
}

//UploadHash will upload the hash fropm the frontend
var UploadHash = func(writer http.ResponseWriter, request *http.Request) {
	payload := &models.Payload{}
	err := json.NewDecoder(request.Body).Decode(payload)
	if err != nil {
		utils.Respond(writer, utils.Message(false, "Error while decoding request body, ensure that it is a string"))
		return
	}
	hash := getHash(payload, writer)
	//Check if the hash exists
	exists, hashErr := checkHashInDb(hash)
	if hashErr != nil {
		fmt.Println(hashErr)
		return
	}
	if !exists {
		//hash doesnt exist
		tx := txBuilder(hash, "eth_sendRawTransaction", writer)
		signedTx := tx.Sign()
		sendSigned(hash, signedTx, writer)
		return
	}
	utils.Respond(writer, utils.Message(false, "Document already signed"))
	// utils.Respond(writer, resp)
}

func checkHashInDb(hash string) (bool, error) {
	hashStruct := &models.HashStruct{}
	err := models.GetDB().Table("hash_structs").Where("hash = ?", hash).First(hashStruct).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//hash doesnt exist
			return false, nil
		}
		return true, err
	}

	return true, nil

}
func sendSigned(hash string, incoming *models.SignedTx, writer http.ResponseWriter) {
	requestBody, err := json.Marshal(incoming)
	if err != nil {
		fmt.Println(err)
		utils.Respond(writer, utils.Message(false, "Failed to decode request payload to geth node"))
	}
	response, responseErr := http.Post(ethRPC, "application/json", bytes.NewBuffer(requestBody))
	if responseErr != nil {
		fmt.Println(responseErr)
	}
	//close payload to prevent leakages
	defer response.Body.Close()

	// read the response
	body, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		fmt.Println(bodyErr)
	}
	//parse the response
	var parsedResponse = new(models.ParsedResponse)

	parsedErr := json.Unmarshal(body, &parsedResponse)
	if parsedErr != nil {
		fmt.Print(parsedErr)
	}
	txHash := parsedResponse.Result
	if txHash == "" {
		utils.Respond(writer, utils.Message(false, "transaction failed"))
		return
	} else {
		resp := make(map[string]interface{})
		resp["transactionHash"] = parsedResponse.Result
		//sucesfully sent so we can sve to DB
		hashStruct := &models.HashStruct{Hash: hash}
		models.GetDB().Create(hashStruct)
		fmt.Println("CREATED")
		utils.Respond(writer, resp)
		return
	}

}

func getPendingTransaction(address string) uint64 {
	var params = []interface{}{address, "latest"}
	nonceDummy := &models.NonceStruct{}
	err := models.GetDB().Table("nonce_structs").Where("address = ?", address).First(nonceDummy).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// params = append(params, "0x0")
			payload := &Payload{
				JSONRPC: "2.0",
				Method:  "eth_getTransactionCount",
				Params:  params,
				ID:      1,
			}
			requestBody, err := json.Marshal(payload)
			// fmt.Println(requestBody)
			if err != nil {
				fmt.Println(err)
			}
			response, responseErr := http.Post(ethRPC, "application/json", bytes.NewBuffer(requestBody))
			if responseErr != nil {
				fmt.Println(responseErr)
			}
			//close payload to prevent leakages
			defer response.Body.Close()

			// read the response
			body, bodyErr := ioutil.ReadAll(response.Body)
			if bodyErr != nil {
				fmt.Println("FAILED TO IO PARSE")
			}
			// parse the response
			var parsedResponse = new(models.ParsedResponse)

			parsedErr := json.Unmarshal(body, &parsedResponse)
			if parsedErr != nil {
				fmt.Print(parsedErr)
			}
			nonceDummy.Address = address
			nonceDummy.Nonce = hextodec(parsedResponse.Result)
			models.GetDB().Create(nonceDummy)
			fmt.Println("Propogating transaction with nonce %d\n", hextodec(parsedResponse.Result))
			return hextodec(parsedResponse.Result)
		}
	}
	//update the nonce
	nonceDummy.Nonce = nonceDummy.Nonce + uint64(1)
	models.GetDB().Model(nonceDummy).Update("nonce", nonceDummy.Nonce)
	fmt.Println("Propogating transaction with nonce %d\n", nonceDummy.Nonce)
	return nonceDummy.Nonce
}
func getHash(payload *models.Payload, writer http.ResponseWriter) string {
	h := sha256.New()
	h.Write([]byte(payload.Name))
	h.Write([]byte(payload.Type))
	h.Write([]byte(strconv.FormatInt(int64(payload.LastModified), 10)))
	h.Write([]byte(strconv.FormatInt(int64(payload.Size), 10)))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash
}

func txBuilder(hash string, method string, writer http.ResponseWriter) *models.Tx {
	nonce := getPendingTransaction("0xE9C0614F054FAd022e989034c00b136E507e162b")
	params := &models.Param{
		Nonce:    nonce,
		From:     "0xE9C0614F054FAd022e989034c00b136E507e162b",
		Gas:      100000,
		GasPrice: 50000000000,
		Value:    0,
		Data:     hash,
	}
	tx := &models.Tx{
		JSONRpc: "2.0",
		Method:  method,
		Params:  *params,
		ID:      1,
	}
	return tx
}

func hextodec(hex string) uint64 {
	ru := fmt.Sprint(hex[2:])
	x, _ := strconv.ParseInt(ru, 16, 64)
	return uint64(x)
}
