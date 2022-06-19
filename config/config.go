package config

import (
	"fmt"
	"os"
	"strconv"
)

// DATABASE
var MONGODB_USERNAME string = os.Getenv("MONGODB_USERNAME")
var MONGODB_PASSWORD string = os.Getenv("MONGODB_PASSWORD")
var MONGODB_DATABASE string = os.Getenv("MONGODB_DATABASE")
var MONGODB_URI string = os.Getenv("MONGODB_URI")
var MONGODB_CONSTR string = fmt.Sprintf(MONGODB_URI, MONGODB_USERNAME, MONGODB_PASSWORD)

var COLLECTION = map[string]string{
	"assets":           "assets",
	"authorizations":   "authorizations",
	"companies":        "companies",
	"cryptocurrencies": "cryptocurrencies",
	"stock":            "stock",
	"transactions":     "transactions",
	"users":            "users",
}

// SECURITY
var AUTH_COOKIE_EXPIRY, _ = strconv.Atoi(os.Getenv("AUTH_COOKIE_EXPIRY"))
var AUTH_SECRET = os.Getenv("AUTH_SECRET")

// USER
var USER_DEFAULT_BALANCE, _ = strconv.ParseFloat(os.Getenv("USER_DEFAULT_BALANCE"), 8)

// ALPHA VANTAGE
var AV_API_KEY string = os.Getenv("ALPHA_VANTAGE_API_KEY")
var AV_SECOND_API_KEY string = os.Getenv("ALPHA_VANTAGE_SECOND_API_KEY")
var AV_THIRD_API_KEY string = os.Getenv("ALPHA_VANTAGE_THIRD_API_KEY")
var AV_FOURTH_API_KEY string = os.Getenv("ALPHA_VANTAGE_FOURTH_API_KEY")
var AV_FIFTH_API_KEY string = os.Getenv("ALPHA_VANTAGE_FIFTH_API_KEY")
var AV_SIXTH_API_KEY string = os.Getenv("ALPHA_VANTAGE_SIXTH_API_KEY")
var AV_BASE_URL string = os.Getenv("ALPHA_VANTAGE_BASE_URL")
var AV_DEFAULT_DATATYPE string = os.Getenv("ALPHA_VANTAGE_DEFAULT_DATATYPE")
var AV_DEFAULT_INTERVAL string = os.Getenv("ALPHA_VANTAGE_DEFAULT_INTERVAL")
var AV_FUNCTIONS = map[string]string{
	"current_company":    "TIME_SERIES_INTRADAY",
	"historical_company": "TIME_SERIES_INTRADAY_EXTENDED",
	"current_crypto":     "CRYPTO_INTRADAY",
	"company_overview":   "OVERVIEW",
	"company_search":     "SYMBOL_SEARCH",
	"exchange_rate":      "CURRENCY_EXCHANGE_RATE",
}
var AV_RESPONSE_SIZE = map[string]string{
	"default":  "compact",
	"extended": "full",
}
