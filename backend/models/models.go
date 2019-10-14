package models

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
)

//Payload is the incoming payload
type Payload struct {
	Name         string
	LastModified int
	Size         int
	Type         string
}

//Tx The unsigned TX to send to the ethereum node
type Tx struct {
	JSONRpc string
	Method  string
	Params  Param
	ID      int
}

//SignedTx is the signed TX
type SignedTx struct {
	JSONRpc string
	Method  string
	Params  []string
	ID      int
}

//Param is the param in the txpayload
type Param struct {
	Nonce    uint64
	From     string
	Gas      uint64
	GasPrice int64
	Value    int64
	Data     string
}

//ParsedResponse is the response when a payload is sent to the server. This will be used for all responses
type ParsedResponse struct {
	ID      int
	JSONRPC string
	Result  string
}

//Sign the transaction
func (incomingTx *Tx) Sign() *SignedTx {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	chainID := big.NewInt(3) // ropsten
	var amount = big.NewInt(incomingTx.Params.Value)
	var gasPrice = big.NewInt(incomingTx.Params.GasPrice)
	var bytesto [20]byte
	_bytesto, _ := hex.DecodeString("")

	senderPrivKey, _ := crypto.HexToECDSA(os.Getenv("PRIV_KEY"))
	copy(bytesto[:], _bytesto)
	to := common.Address([20]byte(bytesto))
	signer := types.NewEIP155Signer(chainID)

	transaction := types.NewTransaction(incomingTx.Params.Nonce, to, amount, incomingTx.Params.Gas, gasPrice, []byte(incomingTx.Params.Data))
	signedTx, _ := types.SignTx(transaction, signer, senderPrivKey)
	ts := types.Transactions{signedTx}
	rawTx := fmt.Sprintf("%x", ts.GetRlp(0))
	rawTx = "0x" + rawTx
	var params []string
	params = append(params, rawTx)
	rawTxPayload := &SignedTx{
		JSONRpc: "2.0",
		Method:  incomingTx.Method,
		Params:  params,
		ID:      incomingTx.ID,
	}
	return rawTxPayload
	// json_tx, _ := signed_tx.MarshalJSON()
	// _ = json.Unmarshal(json_tx, parsedTx)
	// pparsedTx.From = from
	// fmt.Println("data", parsed_tx.Data)
	// return parsed_tx, nil
}
