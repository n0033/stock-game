package api_connector

import models "github.com/ktylus/stock-game/common/models/resources/api_connector"

type APIHTTPConnectorInterface interface {
	GetBase() *models.APIHTTPConnectorBase
	Get() *models.APIHTTPConnectorResponse
}
