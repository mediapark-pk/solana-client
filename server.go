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
	myRouter.HandleFunc("/pKToAddress", pKToAddress).Methods("POST")
	myRouter.HandleFunc("/getBalance", getBalance).Methods("POST")
	myRouter.HandleFunc("/createAndSendTransaction", createAndSendTransaction).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
