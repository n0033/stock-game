package tests

import (
	"fmt"
	"testing"

	"github.com/ktylus/stock-game/config"
	av_adapter "github.com/ktylus/stock-game/services/adapter/alpha_vantage"
	av_connector "github.com/ktylus/stock-game/services/api_connector/connectors/alpha_vantage"
)

var CONNECTOR = av_connector.GetAlphaVantageConnector()

func TestConvertToStockDatapoints(t *testing.T) {
	response := CONNECTOR.GetCurrentCompanyData(APPLE_CODE, config.AV_RESPONSE_SIZE["default"])
	av_datapoints := response.Payload[fmt.Sprintf("Time Series (%s)", config.AV_DEFAULT_INTERVAL)].(map[string]interface{})
	datapoints := av_adapter.ToStockDatapoints(response)
	if len(av_datapoints) != len(datapoints) {
		t.Fail()
	}
}

func TestConvertCompanyOverview(t *testing.T) {
	response := CONNECTOR.GetCompanyOverview(APPLE_CODE)
	av_adapter.ToCompanyOverview(response)
}
