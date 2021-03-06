package main

import (
	jsonLib "encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	//"reflect"

	amazonproduct "github.com/Shadowmad/go-amazon-product-api"
	xj "github.com/basgys/goxml2json"
	"github.com/jmoiron/jsonq"
)

const (
	AWSAccessKeyId = "AKIAIUYTRIWRVKP6B6CA"
	AWSSecret      = "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"
	AssociateTag   = "booksmart044-20"
)

func amazonResponse(product_id *string, channel chan<- []phpResponseStruct, sellOrBuy string) {
	var api amazonproduct.AmazonProductAPI

	api.AccessKey = AWSAccessKeyId
	api.SecretKey = AWSSecret
	api.Host = "webservices.amazon.com"
	api.AssociateTag = AssociateTag
	api.Client = &http.Client{} // optional
	resp, err := api.ItemLookup(*product_id)
	if err == nil {
		json, err := xj.Convert(strings.NewReader(resp))
		if err != nil {
			panic("That's embarrassing...")
		}
		jsonstring := json.String()
		data := map[string]interface{}{}
		dec := jsonLib.NewDecoder(strings.NewReader(jsonstring))
		dec.Decode(&data)
		jq := jsonq.NewQuery(data)
		if sellOrBuy == "buy" {
			compileResponse(jq, channel)
		}
		if sellOrBuy == "sell" {
			compileSellResponse(jq, channel)
		}

	}
}
func compileSellResponse(jsonResp *jsonq.JsonQuery, channel chan<- []phpResponseStruct) {
	checkForEligibility, _ := jsonResp.String("ItemLookupResponse", "Items", "Item", "ItemAttributes", "IsEligibleForTradeIn")
	finalResp  := new(phpHeaderName)
	if checkForEligibility == "1" {
		conditions := [2]string{"new", "good"}
		amazonLink := ""
		merchantName := "amazon"
		merchantImage := ""
		sellPrice, err := jsonResp.String("ItemLookupResponse", "Items", "Item", "ItemAttributes", "TradeInValue", "Amount")
		sellPriceInt, _   := strconv.Atoi(sellPrice)
		if err != nil {
			fmt.Println(err)
		}
		sellBuild := new(phpResponseStruct)
		sellBuild.Merchant  = merchantName
		sellBuild.MerchantImage = merchantImage
		sellBuild.LinkToBuy 		= amazonLink
		sellBuild.Condition	= "complex"
		sellBuild.LinkToBuy = "http://www.amazon.com/s?index=textbooks-tradein&tag=" + AssociateTag + "&field-keywords="
		for i := range conditions {
			if(conditions[i] == "new"){
				sellBuild.BB_cond_new = conditions[i]
				sellBuild.BB_price_new = strconv.FormatFloat((float64(sellPriceInt) / 100), 'f', 2, 64)
			}
			if(conditions[i] == "good"){
				sellBuild.BB_cond_good = conditions[i]
				sellBuild.BB_price_good = strconv.FormatFloat((float64(sellPriceInt) / 100), 'f', 2, 64)
			}
		}
		finalResp.HeaderName = append(finalResp.HeaderName, *sellBuild)
	}
	channel <- finalResp.HeaderName
}
func compileResponse(jsonResp *jsonq.JsonQuery, channel chan<- []phpResponseStruct) {
	offers, _ := jsonResp.ArrayOfObjects("ItemLookupResponse", "Items", "Item", "Offers", "Offer")
	finalResp := new(phpHeaderName)
	for i := range offers {
		readyResp := new(phpResponseStruct)
		condition := offers[i]["OfferAttributes"].(map[string]interface{})["Condition"].(string)
		if condition == "Used" {
			temp, _ := jsonResp.String("ItemLookupResponse", "Items", "Item", "Offers", "MoreOffersUrl")
			readyResp.LinkToBuy = temp + "&condition=" + condition
		} else {
			temp, _ := jsonResp.String("ItemLookupResponse", "Items", "Item", "Offers", "MoreOffersUrl")
			readyResp.LinkToBuy = temp + "&condition=" + condition
		}
		readyResp.Merchant = "amazon"
		readyResp.MerchantImage = "" //to take from config files
		readyResp.Condition = offers[i]["OfferAttributes"].(map[string]interface{})["Condition"].(string)
		AmazonPrice, _ := strconv.ParseFloat(offers[i]["OfferListing"].(map[string]interface{})["Price"].(map[string]interface{})["Amount"].(string), 64)
		readyResp.Price = strconv.FormatFloat(AmazonPrice/100, 'f', 2, 64)
		readyResp.Shipping = "0.00"
		readyResp.TypeOf = "sale"

		itemAmountInt, _ := strconv.Atoi(offers[i]["OfferListing"].(map[string]interface{})["Price"].(map[string]interface{})["Amount"].(string))
		shippingAmountInt, _ := strconv.Atoi(readyResp.Shipping)
		total := strconv.FormatFloat(((float64(itemAmountInt) + float64(shippingAmountInt)) / 100), 'f', 2, 64)
		readyResp.TotalPrice = total

		finalResp.HeaderName = append(finalResp.HeaderName, *readyResp)
	}
	/*c, err := jsonLib.Marshal(finalResp.HeaderName)

	if err != nil {
		panic(err)
	}*/
	channel <- finalResp.HeaderName
}
