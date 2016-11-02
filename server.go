package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/search/{service}/{isbn}", returnCompiledBookSearh)

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {

	handleRequests()

}
