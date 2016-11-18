package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", welcome)
	myRouter.HandleFunc("/search/buy/{isbn}", returnCompiledBookSearh)
	myRouter.HandleFunc("/search/sell/{isbn}", returnCompiledSellBookSearh)

	log.Fatal(http.ListenAndServe("", myRouter))
}

func main() {

	handleRequests()

}
