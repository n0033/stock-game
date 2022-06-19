package alpha_vantage

import "time"

type StockDatapointValues struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

type StockDatapoint struct {
	Symbol    string
	Timestamp time.Time
	Values    StockDatapointValues
	Volume    int
}

type StockDatapointInCSV struct {
	Time   string `csv:"time" json:"time"`
	Open   string `csv:"open" json:"open"`
	High   string `csv:"high" json:"high"`
	Low    string `csv:"low" json:"low"`
	Close  string `csv:"close" json:"close"`
	Volume string `csv:"volume" json:"volume"`
}

type StockDatapointCrypto struct {
	Symbol    string
	Currency  string
	Timestamp time.Time
	Values    StockDatapointValues
	Volume    int
}

type CompanyOverview struct {
	Symbol      string `json:"Symbol"`
	Asset_type  string `json:"AssetType"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Currency    string `json:"Currency"`
	Country     string `json:"Country"`
	Sector      string `json:"Sector"`
	Industry    string `json:"Industry"`
	Address     string `json:"Address"`
}

type CompanySymbolSearch struct {
	Symbol   string `json:"1. symbol"`
	Name     string `json:"2. name"`
	Currency string `json:"8. currency"`
}

type CurrencyExchangeRate struct {
	Currency_from      string
	Currency_from_name string
	Currency_to        string
	Currency_to_name   string
	Exchange_rate      float64
	Bid_price          float64
	Ask_price          float64
	Timestamp          time.Time
}
