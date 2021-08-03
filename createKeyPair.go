package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

// KeyPair struct
type KeyPair struct {
	Status         bool   `json:"status"`
	AccountAddress string `json:"address"`
	PrivateKey     []byte `json:"pk"`
}

// PK struct
type PK struct {
	PrivateKey string `json:"pk"`
}

// ATAParams struct
type ATAParams struct {
	Address string `json:"address"`
	Mint    string `json:"mint"`
}

// ATAResponse struct
type ATAResponse struct {
	Status     bool   `json:"status"`
	AtaAddress string `json:"ataAddress"`
}

// Create fresh keyPair
func createKeyPair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	account := types.NewAccount()
	//Creating a fresh KeyPair
	keypair := KeyPair{Status: true, AccountAddress: account.PublicKey.ToBase58(), PrivateKey: account.PrivateKey}
	json.NewEncoder(w).Encode(keypair)
}

//get Address by privateKey
func pKToAddress(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var privkeybytes PK
	json.Unmarshal(reqBody, &privkeybytes)
	w.Header().Set("Content-Type", "application/json")
	// Decoding base64 Private Key to Bytes
	privBytes, err := Base64Decode([]byte(privkeybytes.PrivateKey))
	if err != nil {
		fatal := Error{Status: false, Message: string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	//Getting Account Address by PrivateKey Bytes
	account := types.AccountFromPrivateKeyBytes(privBytes)
	keypair := KeyPair{Status: true, AccountAddress: account.PublicKey.ToBase58(), PrivateKey: account.PrivateKey}
	json.NewEncoder(w).Encode(keypair)
}

//Get ATA account
func getATA(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var params ATAParams
	json.Unmarshal(reqBody, &params)

	w.Header().Set("Content-Type", "application/json")

	// Basic Validation
	if params.Address == "" {
		fatal := Error{Status: false, Message: "Required Param , Address"}
		json.NewEncoder(w).Encode(fatal)
		return

	}
	if params.Mint == "" {
		fatal := Error{Status: false, Message: "Required Param , Mint Address"}
		json.NewEncoder(w).Encode(fatal)
		return

	}

	//Getting  ATA account
	ata, _, err := common.FindAssociatedTokenAddress(common.PublicKeyFromString(params.Address), common.PublicKeyFromString(params.Mint))
	if err != nil {
		fatal := Error{Status: false, Message: "find SenderAta error, err: " + string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
	}
	ataAddress := ata.ToBase58()

	ataRes := ATAResponse{Status: true, AtaAddress: ataAddress}
	json.NewEncoder(w).Encode(ataRes)
}

//Connection with Client Testing
func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fatal := Error{Status: true, Message: "hello solana "}
	json.NewEncoder(w).Encode(fatal)

}
