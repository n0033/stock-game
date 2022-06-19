package alpha_vantage

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	models "github.com/ktylus/stock-game/common/models/resources/alpha_vantage"
	models_connector "github.com/ktylus/stock-game/common/models/resources/api_connector"
	"github.com/ktylus/stock-game/config"
)

func ToStockDatapoints(response models_connector.APIHTTPConnectorResponse) []models.StockDatapoint {
	result := make([]models.StockDatapoint, 0)
	payload := response.Payload
	_, err := payload["Error Message"]
	if err {
		return result
	}
	metadata := payload["Meta Data"].(map[string]interface{})
	datapoints := payload[fmt.Sprintf("Time Series (%s)", config.AV_DEFAULT_INTERVAL)].(map[string]interface{})

	symbol := metadata["2. Symbol"].(string)

	for key, value := range datapoints {
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:00", key, time.Local)
		if err != nil {
			log.Fatal(err)
		}
		volume, err := strconv.Atoi(value.(map[string]interface{})["5. volume"].(string))
		if err != nil {
			log.Fatal(err)
		}
		open, err := strconv.ParseFloat(value.(map[string]interface{})["1. open"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		high, err := strconv.ParseFloat(value.(map[string]interface{})["2. high"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		low, err := strconv.ParseFloat(value.(map[string]interface{})["3. low"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		close, err := strconv.ParseFloat(value.(map[string]interface{})["4. close"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}

		datapoint := models.StockDatapoint{
			Symbol:    symbol,
			Timestamp: timestamp,
			Volume:    volume,
			Values: models.StockDatapointValues{
				Open:  open,
				High:  high,
				Low:   low,
				Close: close,
			},
		}

		result = append(result, datapoint)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

func HistoricalToStockDatapoint(response models_connector.APIHTTPConnectorResponse) []models.StockDatapoint {
	result := make([]models.StockDatapoint, 0)
	payload := response.Payload

	for key, value := range payload {
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:00", key, time.Local)
		if err != nil {
			log.Fatal(err)
		}
		volume, err := strconv.Atoi(value.(map[string]interface{})["volume"].(string))
		if err != nil {
			log.Fatal(err)
		}
		open, err := strconv.ParseFloat(value.(map[string]interface{})["open"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		high, err := strconv.ParseFloat(value.(map[string]interface{})["high"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		low, err := strconv.ParseFloat(value.(map[string]interface{})["low"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		close, err := strconv.ParseFloat(value.(map[string]interface{})["close"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		symbol := value.(map[string]interface{})["symbol"].(string)

		datapoint := models.StockDatapoint{
			Symbol:    symbol,
			Timestamp: timestamp,
			Volume:    volume,
			Values: models.StockDatapointValues{
				Open:  open,
				High:  high,
				Low:   low,
				Close: close,
			},
		}

		result = append(result, datapoint)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

func CryptoToStockDatapoints(response models_connector.APIHTTPConnectorResponse) []models.StockDatapoint {
	result := make([]models.StockDatapoint, 0)
	payload := response.Payload

	metadata := payload["Meta Data"].(map[string]interface{})
	datapoints := payload[fmt.Sprintf("Time Series Crypto (%s)", config.AV_DEFAULT_INTERVAL)].(map[string]interface{})

	symbol := metadata["2. Digital Currency Code"].(string)

	for key, value := range datapoints {
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:00", key, time.Local)
		if err != nil {
			log.Fatal(err)
		}
		volume := int(value.(map[string]interface{})["5. volume"].(float64))
		open, err := strconv.ParseFloat(value.(map[string]interface{})["1. open"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		high, err := strconv.ParseFloat(value.(map[string]interface{})["2. high"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		low, err := strconv.ParseFloat(value.(map[string]interface{})["3. low"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		close, err := strconv.ParseFloat(value.(map[string]interface{})["4. close"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}

		datapoint := models.StockDatapoint{
			Symbol:    symbol,
			Timestamp: timestamp,
			Volume:    volume,
			Values: models.StockDatapointValues{
				Open:  open,
				High:  high,
				Low:   low,
				Close: close,
			},
		}

		result = append(result, datapoint)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

func ToStockDatapointsCrypto(response models_connector.APIHTTPConnectorResponse) []models.StockDatapointCrypto {
	result := make([]models.StockDatapointCrypto, 0)
	payload := response.Payload

	metadata := payload["Meta Data"].(map[string]interface{})
	datapoints := payload[fmt.Sprintf("Time Series Crypto (%s)", config.AV_DEFAULT_INTERVAL)].(map[string]interface{})

	symbol := metadata["2. Digital Currency Code"].(string)
	currency := metadata["4. Market Code"].(string)

	for key, value := range datapoints {
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:00", key, time.Local)
		if err != nil {
			log.Fatal(err)
		}
		volume := int(value.(map[string]interface{})["5. volume"].(float64))
		open, err := strconv.ParseFloat(value.(map[string]interface{})["1. open"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		high, err := strconv.ParseFloat(value.(map[string]interface{})["2. high"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		low, err := strconv.ParseFloat(value.(map[string]interface{})["3. low"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		close, err := strconv.ParseFloat(value.(map[string]interface{})["4. close"].(string), 64)
		if err != nil {
			log.Fatal(err)
		}

		datapoint := models.StockDatapointCrypto{
			Symbol:    symbol,
			Currency:  currency,
			Timestamp: timestamp,
			Volume:    volume,
			Values: models.StockDatapointValues{
				Open:  open,
				High:  high,
				Low:   low,
				Close: close,
			},
		}

		result = append(result, datapoint)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

func ToCompanyOverview(response models_connector.APIHTTPConnectorResponse) models.CompanyOverview {
	var overview models.CompanyOverview
	bytes, err := json.Marshal(response.Payload)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &overview)
	if err != nil {
		log.Fatal(err)
	}
	return overview
}

func ToCompanySymbolSearch(response models_connector.APIHTTPConnectorResponse) []models.CompanySymbolSearch {
	result := make([]models.CompanySymbolSearch, 0)
	payload := response.Payload
	matches := payload["bestMatches"].([]interface{})
	for _, match := range matches {
		var match_obj models.CompanySymbolSearch
		bytes, err := json.Marshal(match)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(bytes, &match_obj)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, match_obj)
	}

	return result
}

func ToCurrencyExchangeRate(response models_connector.APIHTTPConnectorResponse) models.CurrencyExchangeRate {
	payload := response.Payload
	data := payload["Realtime Currency Exchange Rate"].(map[string]interface{})

	from_currency := data["1. From_Currency Code"].(string)
	from_currency_name := data["2. From_Currency Name"].(string)
	to_currency := data["3. To_Currency Code"].(string)
	to_currency_name := data["4. To_Currency Name"].(string)
	exchange_rate, err := strconv.ParseFloat(data["5. Exchange Rate"].(string), 64)
	if err != nil {
		log.Fatal(err)
	}
	bid_price, err := strconv.ParseFloat(data["8. Bid Price"].(string), 64)
	if err != nil {
		log.Fatal(err)
	}
	ask_price, err := strconv.ParseFloat(data["9. Ask Price"].(string), 64)
	if err != nil {
		log.Fatal(err)
	}
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:01", data["6. Last Refreshed"].(string), time.Local)
	if err != nil {
		log.Fatal(err)
	}

	return models.CurrencyExchangeRate{
		Currency_from:      from_currency,
		Currency_from_name: from_currency_name,
		Currency_to:        to_currency,
		Currency_to_name:   to_currency_name,
		Exchange_rate:      exchange_rate,
		Bid_price:          bid_price,
		Ask_price:          ask_price,
		Timestamp:          timestamp,
	}
}
