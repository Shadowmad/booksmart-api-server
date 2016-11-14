package main

import (
	jsonLib "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/jmoiron/jsonq"
)

const (
	abeKey = "9318d400-d266-47b6-8888-679ca7a4c8ef"
	abeApi = "http://search2.abebooks.com/search"
)

var abeCond = [2]string{"newonly", "usedonly"}

func abeResponseInit(isbn *string, channel chan<- []phpResponseStruct) {
	for i := range abeCond {
		go abeResponse(isbn, channel, &abeCond[i])
	}
}

func abeResponse(isbn *string, channel chan<- []phpResponseStruct, cond *string) {
	req, err := http.NewRequest("GET", abeApi, nil)

	if err != nil {
		panic("Request did not build...")
	}
	//defer req.Body.Close()
	/**
	  @Purpose Sets up url string
	*/
	query := req.URL.Query()
	query.Add("clientkey", abeKey)
	query.Add("sortorder", "17")
	query.Add("isbn", *isbn)
	query.Add("destinationcountry", "usa")
	query.Add("outputsize", "micro")
	query.Add("maxresults", "1")
	query.Add("bookcondition", *cond)
	req.URL.RawQuery = query.Encode()
	/**
	  @Purpose Start fetching
	*/
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic("Client could not fetch data...")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	json, err := xj.Convert(strings.NewReader(string(body)))
	if err != nil {
		panic("That's embarrassing...")
	}
	jsonstring := json.String()
	data := map[string]interface{}{}
	dec := jsonLib.NewDecoder(strings.NewReader(jsonstring))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	parseAbeResponse(jq, channel, cond)
}

func parseAbeResponse(jq *jsonq.JsonQuery, channel chan<- []phpResponseStruct, cond *string) {
	offers, err := jq.Interface("searchResults", "Book")

	if(offers != nil){
		itemCond := ""
		buildResp := new(phpResponseStruct)
		if err != nil {
			fmt.Println(err)
		}
		if *cond == "newonly" {
			itemCond = "New"
		} else if *cond == "usedonly" {
			itemCond = "Used"
		}
		buildResp.Condition = itemCond
		buildResp.Merchant = "abebooks"
		buildResp.MerchantImage = ""
		buildResp.TypeOf = "buy"
		buildResp.Price = offers.(map[string]interface{})["listingPrice"].(string)
		buildResp.Shipping = offers.(map[string]interface{})["firstBookShipCost"].(string)
		buildResp.TotalPrice = strconv.FormatFloat(getTotalPrice(buildResp), 'f', 2, 64)
		buildResp.LinkToBuy = "http://" + offers.(map[string]interface{})["listingUrl"].(string)

		finalResp := new(phpHeaderName)
		finalResp.HeaderName = append(finalResp.HeaderName, *buildResp)
		channel <- finalResp.HeaderName
	}else{
		close(channel)
	}
}
