package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/createKeyPair", createKeyPair)
	myRouter.HandleFunc("/hello", hello)
	myRouter.HandleFunc("/pKToAddress", pKToAddress).Methods("POST")
	myRouter.HandleFunc("/getBalance", getBalance).Methods("POST")
	myRouter.HandleFunc("/createAndSendTransaction", createAndSendTransaction).Methods("POST")
	myRouter.HandleFunc("/getAtaAddress", getATA).Methods("POST")
	myRouter.HandleFunc("/transferTokens", transferTokens).Methods("POST")
	myRouter.HandleFunc("/createAtaAndTransferTokens", createATAAndtransferTokens).Methods("POST")
	log.Fatal(http.ListenAndServe(":12345", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
