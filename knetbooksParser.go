package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	knetToken  = "iOrXcqEJgTC1Cxq1"
	knetSecret = "WSk3ii4l9HTceu2jHPk4et6s4GLrdd8q"
	knetApi    = "https://api.shareasale.com/x.cfm"
	knetAction = "merchantStatus"
	knetAffil  = "1380442"
	knetVer    = "2.1"
	knetXML    = "1"
	knetFormat = "xml"
	BlockSize  = 64
)

func generateSignature() (sha string) {
	stringToSign := knetToken + ":" + generateCurrrentDateUTC() + ":" + knetAction + ":" + knetSecret
	hasher := sha256.New()
	hasher.Write([]byte(stringToSign))
	sha = strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))

	fmt.Println(sha)
	return sha
}
func generateCurrrentDateUTC() (timeR string) {
	timeR = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	fmt.Println(timeR)
	return timeR
}
func knetBooksResponse() {
	req, err := http.NewRequest("GET", "https://api.shareasale.com/x.cfm?action=merchantStatus&affiliateId=1380442&token=iOrXcqEJgTC1Cxq1&programStatus=onlineNotLowFunds&version=2.1", nil)
	if err != nil {
		panic("Request did not build...")
	}
	//defer req.Body.Close()
	/**
	  @Purpose Sets up url string
	*/
	// query := req.URL.Query()
	// query.Add("version", knetVer)
	// query.Add("action", knetAction)
	// query.Add("affiliateId", knetAffil)
	// query.Add("token", knetToken)
	// query.Add("keyword", "Economics")
	// query.Add("XMLFormat", knetXML)
	// query.Add("format", knetFormat)
	// query.Add("merchantId", "31586")
	// req.URL.RawQuery = query.Encode()

	req.Header.Set("x-ShareASale-Date", generateCurrrentDateUTC())
	req.Header.Set("x-ShareASale-Authentication", generateSignature())
	fmt.Println(req)
	/**
	  @Purpose Start fetching
	*/
	// client := &http.Client{}
	// response, err := client.Do(req)
	// if err != nil {
	// 	panic("Client could not fetch data...")
	// }
	// body, errRespBody := ioutil.ReadAll(response.Body)
	// if errRespBody != nil {
	// 	panic("Cannot read body data from response...")
	// }
	// //Turn JSON response into Golang map for futher processing
	// if body != nil {
	// 	jsonК, err := xj.Convert(strings.NewReader(string(body)))
	// 	if err != nil {
	// 		panic("That's embarrassing...")
	// 	}
	// 	jsonstring := jsonК.String()
	// 	data := map[string]interface{}{}
	// 	dec := json.NewDecoder(strings.NewReader(jsonstring))
	// 	dec.Decode(&data)
	// 	jq := jsonq.NewQuery(data)
	// 	fmt.Println(jq)
	// 	//process map
	// }
}
