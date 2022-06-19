package alpha_vantage_connector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gocarina/gocsv"
	models_av "github.com/ktylus/stock-game/common/models/resources/alpha_vantage"
	models_connector "github.com/ktylus/stock-game/common/models/resources/api_connector"
	"github.com/ktylus/stock-game/config"
)

type AlphaVantageConnector struct {
	base     models_connector.APIHTTPConnectorBase
	api_keys []string
}

var connector *AlphaVantageConnector
var http_client *http.Client

func newAlphaVantageConnector() *AlphaVantageConnector {
	routes := make(map[string]string)

	routes["query"] = "/query"
	base := models_connector.APIHTTPConnectorBase{
		Base_uri: config.AV_BASE_URL,
		Routes:   routes,
	}
	connector := AlphaVantageConnector{
		base:     base,
		api_keys: []string{config.AV_API_KEY},
	}
	return &connector
}

func GetAlphaVantageConnector() *AlphaVantageConnector {
	if connector == nil {
		connector = newAlphaVantageConnector()
	}
	if http_client == nil {
		http_client = &http.Client{}
	}
	return connector
}

func (connector *AlphaVantageConnector) GetBase() models_connector.APIHTTPConnectorBase {
	return connector.base
}

func (connector *AlphaVantageConnector) parseJSONBody(payload []byte) map[string]interface{} {
	body_map := make(map[string]interface{})
	err := json.Unmarshal(payload, &body_map)
	if err != nil {
		log.Fatal(err)
	}
	return body_map
}

func (connector *AlphaVantageConnector) parseCSVBody(payload []byte) map[string]interface{} {
	result := make(map[string]interface{})
	var body_map []map[string]interface{}
	var datapoints []models_av.StockDatapointInCSV
	err := gocsv.UnmarshalBytes(payload, &datapoints)

	if err != nil {
		log.Fatal(err)
	}

	datapoints_bytes, _ := json.Marshal(datapoints)

	json.Unmarshal(datapoints_bytes, &body_map)

	for _, val := range body_map {
		result[val["time"].(string)] = val
	}
	return result
}

func (connector *AlphaVantageConnector) Get(url string, query map[string]string) models_connector.APIHTTPConnectorResponse {
	var response *http.Response
	var parsed map[string]interface{}

	for _, api_key := range connector.api_keys { // alpha vantage is free service with limited amount of requests
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}

		query_params := request.URL.Query()
		for key, val := range query {
			query_params.Add(key, val)
		}
		query_params.Add("apikey", api_key)
		request.URL.RawQuery = query_params.Encode()

		response, err = http_client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		body_bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		if response.Header.Get("Content-Type") == "application/json" {
			parsed = connector.parseJSONBody(body_bytes)
		}

		if response.Header.Get("Content-Type") == "application/x-download" {
			parsed = connector.parseCSVBody(body_bytes)
			_, exists := query["symbol"]

			if exists {
				for key := range parsed {
					parsed[key].(map[string]interface{})["symbol"] = query["symbol"]
				}
			}
		}

		_, exists := parsed["Note"] // exceeded daily of per-minute call limit
		if !exists {
			break
		}

	}

	_, exists := parsed["Note"]
	if exists {
		fmt.Println("Limit of calls to alpha vantage exceeded")
	}

	_, exists = parsed["Error Message"]
	if exists {
		fmt.Println("Invalid request")
	}

	connector_response := models_connector.APIHTTPConnectorResponse{
		Status_code: response.StatusCode,
		Payload:     parsed,
	}

	return connector_response
}

func (connector *AlphaVantageConnector) GetCurrentCompanyData(code string, response_size string) models_connector.APIHTTPConnectorResponse {
	url := connector.base.Base_uri + connector.base.Routes["query"]
	query := make(map[string]string)
	query["function"] = config.AV_FUNCTIONS["current_company"]
	query["symbol"] = code
	query["interval"] = config.AV_DEFAULT_INTERVAL
	query["datatype"] = config.AV_DEFAULT_DATATYPE
	query["outputsize"] = response_size
	return connector.Get(url, query)
}

func (connector *AlphaVantageConnector) GetHistoricalCompanyData(code string, year_num int, month_num int) models_connector.APIHTTPConnectorResponse {
	// year_num and month_num stand for number of month in the past starting from current date
	url := connector.base.Base_uri + connector.base.Routes["query"]
	query := make(map[string]string)
	query["function"] = config.AV_FUNCTIONS["historical_company"]
	query["datatype"] = config.AV_DEFAULT_DATATYPE
	query["symbol"] = code
	query["interval"] = config.AV_DEFAULT_INTERVAL
	query["slice"] = "year" + strconv.Itoa(year_num) + "month" + strconv.Itoa(month_num)
	return connector.Get(url, query)
}

func (connector *AlphaVantageConnector) SearchForCompany(term string) models_connector.APIHTTPConnectorResponse {
	url := connector.base.Base_uri + connector.base.Routes["query"]
	query := make(map[string]string)
	query["function"] = config.AV_FUNCTIONS["company_search"]
	query["datatype"] = config.AV_DEFAULT_DATATYPE
	query["keywords"] = term
	return connector.Get(url, query)
}

func (connector *AlphaVantageConnector) GetCompanyOverview(code string) models_connector.APIHTTPConnectorResponse {
	url := connector.base.Base_uri + connector.base.Routes["query"]
	query := make(map[string]string)
	query["function"] = config.AV_FUNCTIONS["company_overview"]
	query["symbol"] = code
	return connector.Get(url, query)
}

func (connector *AlphaVantageConnector) GetCryptoCurrencyExchangeRate(crypto_code string, code_to string) models_connector.APIHTTPConnectorResponse {
	// code_from - source currency code
	// code_to - currency code to get rate for
	url := connector.base.Base_uri + connector.base.Routes["query"]
	query := make(map[string]string)
	query["function"] = config.AV_FUNCTIONS["exchange_rate"]
	query["from_currency"] = crypto_code
	query["to_currency"] = code_to
	return connector.Get(url, query)
}

func (connector *AlphaVantageConnector) GetCurrentCryptoData(crypto_code string, currency_code string) models_connector.APIHTTPConnectorResponse {
	// code_from - source currency code
	// code_to - currency code to get rate for
	url := connector.base.Base_uri + connector.base.Routes["query"]
	query := make(map[string]string)
	query["function"] = config.AV_FUNCTIONS["current_crypto"]
	query["symbol"] = crypto_code
	query["market"] = currency_code
	query["interval"] = config.AV_DEFAULT_INTERVAL
	return connector.Get(url, query)
}
