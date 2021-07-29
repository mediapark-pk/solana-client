package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/portto/solana-go-sdk/client"
)

// Address params
type Address struct {
	Address string `json:"address"`
}

// Balance response
type Balance struct {
	Status  bool   `json:"status"`
	Balance uint64 `json:"balance"`
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var address Address
	json.Unmarshal(reqBody, &address)
	c := client.NewClient(client.MainnetRPCEndpoint)
	balance, err := c.GetBalance(context.Background(), address.Address)
	if err != nil {
		fatal := Error{Status: false, Message: "get balance error " + string(err.Error())}
		json.NewEncoder(w).Encode(fatal)
		return
		log.Fatalln("get balance error", err)
	}
	accountBalance := Balance{Status: true, Balance: balance}
	json.NewEncoder(w).Encode(accountBalance)
}
