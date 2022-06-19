package api_connector

type APIHTTPConnectorBase struct {
	Base_uri string
	Routes   map[string]string
}

type APIHTTPConnectorResponse struct {
	Status_code int                    `json:"status_code"`
	Payload     map[string]interface{} `json:"payload"`
}
