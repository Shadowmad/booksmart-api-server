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
  BB_cond_good  string `json:"bb_cond_good"`
  BB_cond_new   string `json:"bb_cond_new"`
  BB_price_new  string `json:"bb_price_new"`
  BB_price_good string `json:"bb_price_good"`
}
type phpHeaderName struct{
  HeaderName    []phpResponseStruct
}

type phpResponseSellStruct struct{
  Merchant      string `json:"merchant_name"`
  MerchantImage string `json:"merchant_image"`
  Condition     string `json:"condition"`
  Price         string `json:"price"`
  LinkToSell    string `json:"link_to_sell"`
}
type phpSellHeaderName struct{
  HeaderName    []phpResponseSellStruct
}
