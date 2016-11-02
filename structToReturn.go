package main

type phpResponseStruct struct{
  Merchant      string `json:"merchant_name"`
  MerchantImage string `json:"merchant_image"`
  Condition     string `json:"condition"`
  Price         string `json:"price"`
  Shipping      string `json:"shipping"`
  TotalPrice    string `json:"total_price"`
  LinkToBuy     string `json:"link_to_buy"`
  TypeOf        string `json:"type_of"`
}
type phpHeaderName struct{
  HeaderName    []phpResponseStruct
}
