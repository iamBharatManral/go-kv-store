package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iamBharatManral/go-kvstore/store"
	"github.com/iamBharatManral/go-kvstore/transaction"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var logger transaction.TransactionLogger

func initializeTransactionLog() error {
	var err error
	//logger, err = transaction.NewFileTransactionLogger("transaction.log")
	logger, err = transaction.NewPostgresTransactionLogger(transaction.PostgresDBParams{
		Host:     "localhost",
		DBName:   "kvs",
		User:     "test",
		Password: "hunter2",
	})
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errors := logger.ReadEvents()
	e, ok := transaction.Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case transaction.EventDelete:
				err = store.Delete(e.Key)
			case transaction.EventPut:
				err = store.Put(e.Key, e.Value)
			}
		}
	}
	logger.Run()
	return err
}

func main() {
	err := initializeTransactionLog()
	if err != nil {
		log.Fatal(err.Error())
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", PutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", GetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", DeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
