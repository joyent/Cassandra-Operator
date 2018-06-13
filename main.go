package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/niravpatel27/cassandra-operator-workshop/account"
	"github.com/niravpatel27/cassandra-operator-workshop/cassandra"
)

type healthCheckResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type AccountRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	City  string `json:"city"`
}

func main() {
	fmt.Println("starting http server")
	cassandraSession := cassandra.Session
	defer cassandraSession.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", healthCheck)
	router.HandleFunc("/accounts", getAccounts).Methods("GET")
	router.HandleFunc("/account", createAccount).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, healthCheckResponse{Status: "OK", Code: 200})
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	a := account.Account{}
	accounts, err := a.GetAccounts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, accounts)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	var acctReq AccountRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&acctReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	fmt.Printf("%+v\n", acctReq)

	// generate new UUID
	var gocqlUUID gocql.UUID
	gocqlUUID = gocql.TimeUUID()

	acct := account.Account{ID: gocqlUUID, Name: acctReq.Name, Age: acctReq.Age, City: acctReq.City, Email: acctReq.Email}
	if err := acct.CreateAccount(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, acct)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
