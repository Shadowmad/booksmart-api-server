package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", welcome)
	myRouter.HandleFunc("/search/{service}/{isbn}", returnCompiledBookSearh)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func main() {

	handleRequests()

}
