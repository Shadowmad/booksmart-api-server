package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
func returnCompiledBookSearh(w http.ResponseWriter, r *http.Request) {
	//isbn := os.Args[1]
	queryVar := mux.Vars(r)
	isbn := queryVar["isbn"] //"0517548232" //
	fmt.Println(isbn)
	compiledResp := new(phpHeaderName)
	resp_channel := make(chan []phpResponseStruct)
	go amazonResponse(&isbn, resp_channel)
	go cheggResponse(&isbn, resp_channel)
	go abeResponseInit(&isbn, resp_channel)
	for i := 1; i <= 4; i++ {
		for _, resp_value := range <-resp_channel {
			compiledResp.HeaderName = append(compiledResp.HeaderName, resp_value)
		}
	}
	fmt.Println("Endpoint Hit: returnAllArticles")

	json.NewEncoder(w).Encode(compiledResp.HeaderName)
}
