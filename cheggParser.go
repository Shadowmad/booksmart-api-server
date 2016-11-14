package main

import (
	jsonLib "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/jsonq"
)

const (
	CheggKey  = "735427de266ea862e387e493f6c2527a"
	CheggPass = "7187436"
	cheggApi  = "http://api.chegg.com/rent.svc"
)

func cheggResponse(isbn *string, channel chan<- []phpResponseStruct) {

	req, err := http.NewRequest("GET", cheggApi, nil)

	if err != nil {
		panic("Request did not build...")
	}
	//defer req.Body.Close()
	/**
	  @Purpose Sets up url string
	*/
	query := req.URL.Query()
	query.Add("KEY", CheggKey)
	query.Add("PW", CheggPass)
	query.Add("isbn", *isbn)
	query.Add("R", "JSON")
	query.Add("with_pids", "1")
	query.Add("page", "1")
	query.Add("results_per_page", "50")
	query.Add("V", "4.0")
	req.URL.RawQuery = query.Encode()
	/**
	  @Purpose Start fetching
	*/
	client := &http.Client{}
	fmt.Println(req)
	response, err := client.Do(req)
	if err != nil {
		panic("Client could not fetch data...")
	}
	body, errRespBody := ioutil.ReadAll(response.Body)
	if errRespBody != nil {
		panic("Cannot read body data from response...")
	}
	//Turn JSON response into Golang map for futher processing
	if body != nil {
		jsonstring := string(body)
		data := map[string]interface{}{}
		dec := jsonLib.NewDecoder(strings.NewReader(jsonstring))
		dec.Decode(&data)
		jq := jsonq.NewQuery(data)
		//process map
		compileCheggResponse(jq, channel, isbn)
	}
}

func compileCheggResponse(jsonResp *jsonq.JsonQuery, channel chan<- []phpResponseStruct, isbn *string) {
	var rentDone = false
	shipingPrice := getShippingPrice(jsonResp) //Will use global
	finalResp := new(phpHeaderName)

	for index := 0; index < 2; index++ {
		buildResp := new(phpResponseStruct)
		buildResp.Merchant = "chegg"
		buildResp.MerchantImage = ""
		buildResp.Shipping = shipingPrice
		buildResp.LinkToBuy = "http://chggtrx.com/click.track?CID=267582&AFID=412909&ADID=1088031&SID=&isbn_ean=" + *isbn

		//Rent
		rent, _ := jsonResp.Bool("Data", "Items", "0", "Renting")
		if rent && !rentDone {
			val, err := jsonResp.ArrayOfObjects("Data", "Items", "0", "Terms")
			if err != nil {
				panic(err)
			}
			rentDone = true
			buildResp.Condition = "Rental"
			buildResp.Price = compileRentStruct(&val)
			buildResp.TotalPrice = strconv.FormatFloat(getTotalPrice(buildResp), 'f', 2, 64)
			buildResp.TypeOf = "rent"
			finalResp.HeaderName = append(finalResp.HeaderName, *buildResp)
			continue
		}

		//Sales
		sale, err := jsonResp.ArrayOfObjects("Data", "Items", "0", "SellPrices")
		if err != nil {
			panic(err)
		}

		buildResp.TypeOf = "sale"
		//Get condition
		cond, _ := sale[0]["desc"].(string)
		buildResp.Condition = cond
		buildResp.Price, _ = sale[0]["price"].(string)
		buildResp.TotalPrice = strconv.FormatFloat(getTotalPrice(buildResp), 'f', 2, 64)
		finalResp.HeaderName = append(finalResp.HeaderName, *buildResp)
	}
	/*push, err := jsonLib.Marshal(finalResp.HeaderName)
	if err != nil {
		panic(err)
	}*/
	channel <- finalResp.HeaderName
}

//Calculat total prices
func getTotalPrice(interimResp *phpResponseStruct) (total float64) {
	price, _ := strconv.ParseFloat(interimResp.Price, 64)
	ship, _ := strconv.ParseFloat(interimResp.Shipping, 64)
	total = price + ship
	return total
}

//Get shipping prices
func getShippingPrice(shippings *jsonq.JsonQuery) (cost string) {
	//Shipping
	shipping, _ := (shippings).ArrayOfObjects("Data", "Items", "0", "ShippingPrices")
	for i := range shipping {
		if shipping[i]["method_name"] == "Ground Shipping" {
			cost = shipping[i]["cost_first"].(string)
		}
	}
	return cost
}

//Renting struct build
func compileRentStruct(rents *[]map[string]interface{}) (rentCost string) {
	for i := range *rents {
		if (*rents)[i]["term"] == "SEMESTER" {
			rentCost = (*rents)[i]["price"].(string)
		}
	}
	return rentCost
}
