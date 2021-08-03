package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/portto/solana-go-sdk/assotokenprog"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

// TokenTransferParams struct
type TokenTransferParams struct {
	SenderPrivateKey string  `json:"senderPK"`
	RecipientAddress string  `json:"recipientAddress"`
	Amount           float64 `json:"amount"`
	Scale            float64 `json:"scale"`
	Mint             string  `json:"mint"`
}

// TokenTransferResponse struct
type TokenTransferResponse struct {
	Status       bool   `json:"status"`
	SenderAta    string `json:"senderAta"`
	RecipientAta string `json:"recipientAta"`
	Hash         string `json:"hash"`
}

//Token Transfer only if already ATA exist
func transferTokens(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var params TokenTransferParams
	json.Unmarshal(reqBody, &params)

	w.Header().Set("Content-Type", "application/json")

	// Basic Validation
	if params.SenderPrivateKey == "" {
		fatal := Error{Status: false, Message: "Required Param , senderPK"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.RecipientAddress == "" {
		fatal := Error{Status: false, Message: "Required Param , recipientAddress"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.Mint == "" {
		fatal := Error{Status: false, Message: "Required Param , mint"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.Amount == 0 {
		fatal := Error{Status: false, Message: "Required Param , amount"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.Scale == 0 {
		fatal := Error{Status: false, Message: "Required Param , Scale"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if math.IsNaN(params.Amount) {
		fatal := Error{Status: false, Message: "amount is not a number"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if math.IsNaN(params.Scale) {
		fatal := Error{Status: false, Message: "Scale is not a number"}
		json.NewEncoder(w).Encode(fatal)
		return

	}

	// Decoding base64 Private Key to Bytes
	senderPK, err := Base64Decode([]byte(params.SenderPrivateKey))
	if err != nil {
		fatal := Error{Status: false, Message: string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}

	//Calculating Amount to be transfer
	amountToBeTransfer := math.Pow(10, params.Scale)
	amountToBeTransfer = params.Amount * amountToBeTransfer

	c := client.NewClient(client.MainnetRPCEndpoint)

	//Getting Sender Account
	sender := types.AccountFromPrivateKeyBytes(senderPK)

	//Getting Sender ATA account
	senderAta, _, err := common.FindAssociatedTokenAddress(common.PublicKeyFromString(sender.PublicKey.ToBase58()), common.PublicKeyFromString(params.Mint))
	if err != nil {
		fatal := Error{Status: false, Message: "find SenderAta error, err: " + string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	senderAtaAddress := senderAta.ToBase58()

	//Getting Recipient ATA Account
	recipientATA, _, err := common.FindAssociatedTokenAddress(common.PublicKeyFromString(params.RecipientAddress), common.PublicKeyFromString(params.Mint))
	if err != nil {
		fatal := Error{Status: false, Message: "find RecipientATA error, err: " + string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	recipientATAAddress := recipientATA.ToBase58()

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fatal := Error{Status: false, Message: string("get recent block hash error, error " + string(err.Error()))}
		json.NewEncoder(w).Encode(fatal)
		return
	}

	// Create Raw Transaction for Tokens transfer
	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			tokenprog.TransferChecked(
				common.PublicKeyFromString(senderAtaAddress),
				common.PublicKeyFromString(recipientATAAddress),
				common.PublicKeyFromString(params.Mint),
				sender.PublicKey,
				[]common.PublicKey{},
				uint64(amountToBeTransfer),
				uint8(params.Scale),
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

	//Sending tx
	txhash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fatal := Error{Status: false, Message: string("send raw tx error, error " + string(err.Error()))}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	txSuccess := TokenTransferResponse{Status: true, SenderAta: senderAtaAddress, RecipientAta: recipientATAAddress, Hash: txhash}
	json.NewEncoder(w).Encode(txSuccess)

}

//Create Recipent ATA address and send tokens
func createATAAndtransferTokens(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var params TokenTransferParams
	json.Unmarshal(reqBody, &params)

	w.Header().Set("Content-Type", "application/json")

	// Basic Validation
	if params.SenderPrivateKey == "" {
		fatal := Error{Status: false, Message: "Required Param , senderPK"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.RecipientAddress == "" {
		fatal := Error{Status: false, Message: "Required Param , recipientAddress"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.Mint == "" {
		fatal := Error{Status: false, Message: "Required Param , mint"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.Amount == 0 {
		fatal := Error{Status: false, Message: "Required Param , amount"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.Scale == 0 {
		fatal := Error{Status: false, Message: "Required Param , Scale"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if math.IsNaN(params.Amount) {
		fatal := Error{Status: false, Message: "amount is not a number"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if math.IsNaN(params.Scale) {
		fatal := Error{Status: false, Message: "Scale is not a number"}
		json.NewEncoder(w).Encode(fatal)
		return

	}

	// Decoding base64 Private Key to Bytes
	senderPK, err := Base64Decode([]byte(params.SenderPrivateKey))
	if err != nil {
		fatal := Error{Status: false, Message: string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}

	//Calculating Amount to be transfer
	amountToBeTransfer := math.Pow(10, params.Scale)
	amountToBeTransfer = params.Amount * amountToBeTransfer

	c := client.NewClient(client.MainnetRPCEndpoint)

	//Getting Sender Account
	sender := types.AccountFromPrivateKeyBytes(senderPK)

	//Getting Sender ATA account
	senderAta, _, err := common.FindAssociatedTokenAddress(common.PublicKeyFromString(sender.PublicKey.ToBase58()), common.PublicKeyFromString(params.Mint))
	if err != nil {
		fatal := Error{Status: false, Message: "find SenderAta error, err: " + string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	senderAtaAddress := senderAta.ToBase58()

	//Getting Recipient ATA Account
	recipientATA, _, err := common.FindAssociatedTokenAddress(common.PublicKeyFromString(params.RecipientAddress), common.PublicKeyFromString(params.Mint))
	if err != nil {
		fatal := Error{Status: false, Message: "find RecipientATA error, err: " + string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	recipientATAAddress := recipientATA.ToBase58()

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		fatal := Error{Status: false, Message: string("get recent block hash error, error " + string(err.Error()))}
		json.NewEncoder(w).Encode(fatal)
		return
	}

	// Create Raw Transaction with two instructions one is for ATA account create for reciepent and second for token transfer
	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			assotokenprog.CreateAssociatedTokenAccount(
				sender.PublicKey,
				common.PublicKeyFromString(params.RecipientAddress),
				common.PublicKeyFromString(params.Mint),
			),
			tokenprog.TransferChecked(
				common.PublicKeyFromString(senderAtaAddress),
				common.PublicKeyFromString(recipientATAAddress),
				common.PublicKeyFromString(params.Mint),
				sender.PublicKey,
				[]common.PublicKey{},
				uint64(amountToBeTransfer),
				uint8(params.Scale),
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

	//Sending tx
	txhash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		fatal := Error{Status: false, Message: string("send raw tx error, error " + string(err.Error()))}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	txSuccess := TokenTransferResponse{Status: true, SenderAta: senderAtaAddress, RecipientAta: recipientATAAddress, Hash: txhash}
	json.NewEncoder(w).Encode(txSuccess)

}
