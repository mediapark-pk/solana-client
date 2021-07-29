package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/types"
)

// TransactionObject struct
type TransactionObject struct {
	SenderPrivateKey string  `json:"senderPK"`
	RecipientAddress string  `json:"recipientAddress"`
	Amount           float64 `json:"amount"`
}

// TransactionSuccess struct
type TransactionSuccess struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func createAndSendTransaction(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var transactionObj TransactionObject
	json.Unmarshal(reqBody, &transactionObj)

	if transactionObj.SenderPrivateKey == "" {
		fatal := Error{Status: false, Message: "Required Param , senderPK"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if transactionObj.RecipientAddress == "" {
		fatal := Error{Status: false, Message: "Required Param , recipientAddress"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if transactionObj.Amount == 0 {
		fatal := Error{Status: false, Message: "Required Param , amount"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	fmt.Println(transactionObj.Amount)
	if math.IsNaN(transactionObj.Amount) {
		fatal := Error{Status: false, Message: "amount is not a number"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	senderPK, err := Base64Decode([]byte(transactionObj.SenderPrivateKey))
	if err != nil {
		fatal := Error{Status: false, Message: string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	if math.IsNaN(float64(transactionObj.Amount)) {
		fatal := Error{Status: false, Message: "given amount is not a number"}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	amountToBeTransfer := math.Pow(10, 9)
	amountToBeTransfer = transactionObj.Amount * amountToBeTransfer
	c := client.NewClient(client.MainnetRPCEndpoint)
	sender := types.AccountFromPrivateKeyBytes(senderPK)
	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fatal := Error{Status: false, Message: string("get recent block hash error, error " + string(err.Error()))}
		json.NewEncoder(w).Encode(fatal)
	}
	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.Transfer(
				sender.PublicKey, // from
				common.PublicKeyFromString(transactionObj.RecipientAddress), // to
				uint64(amountToBeTransfer),
			),
		},
		Signers:         []types.Account{sender},
		FeePayer:        sender.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		fatal := Error{Status: false, Message: string("generate tx error, error " + string(err.Error()))}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	txhash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fatal := Error{Status: false, Message: string("send raw tx error, error " + string(err.Error()))}
		json.NewEncoder(w).Encode(fatal)
		return
	}

	txSuccess := TransactionSuccess{Status: true, Message: txhash}
	json.NewEncoder(w).Encode(txSuccess)
}
