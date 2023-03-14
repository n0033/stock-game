package search

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"

	company_dao "github.com/n0033/stock-game/common/dao/company"
	crypto_dao "github.com/n0033/stock-game/common/dao/cryptocurrency"
	models_company "github.com/n0033/stock-game/common/models/mongo/company"
	models_av "github.com/n0033/stock-game/common/models/resources/alpha_vantage"
	models_search "github.com/n0033/stock-game/common/models/resources/search"
	av_adapter "github.com/n0033/stock-game/services/adapter/alpha_vantage"
	av_connector "github.com/n0033/stock-game/services/api_connector/connectors/alpha_vantage"
)

func CompanySearch(c *fiber.Ctx) error {
	results := make([]map[string]interface{}, 0)
	var search_request models_search.SearchRequest
	av_connector := av_connector.GetAlphaVantageConnector()

	if err := c.BodyParser(&search_request); err != nil {
		log.Fatal(err)
	}

	dao_company := company_dao.NewDAOCompany()
	dao_crypto := crypto_dao.NewDAOCryptocurrency()

	db_companies := dao_company.Search(search_request.Term)
	db_crypto := dao_crypto.Search(search_request.Term)

	var companies []models_av.CompanySymbolSearch

	if len(db_companies) > 0 || len(db_crypto) > 0 {
		for _, company := range db_companies {
			result := make(map[string]interface{})

			search_result := models_search.SearchResult{
				Symbol: company.Code,
				Name:   company.Name,
			}
			bytes, _ := json.Marshal(search_result)
			json.Unmarshal(bytes, &result)
			results = append(results, result)
		}

		for _, crypto := range db_crypto {
			result := make(map[string]interface{})

			search_result := models_search.SearchResult{
				Symbol: crypto.Code,
				Name:   crypto.Name,
			}
			bytes, _ := json.Marshal(search_result)
			json.Unmarshal(bytes, &result)
			results = append(results, result)
		}
		return c.JSON(results)
	}

	if len(db_companies) == 0 && len(db_crypto) == 0 {
		connector_response := av_connector.SearchForCompany(search_request.Term)
		companies = av_adapter.ToCompanySymbolSearch(connector_response)
		for _, value := range companies {
			_, err := dao_company.FindByCode(value.Symbol)
			if err != nil {
				dao_company.Create(models_company.CompanyInCreate{
					CompanyBase: models_company.CompanyBase{
						Code: value.Symbol,
						Name: value.Name},
				})
			}
		}
	}

	for _, value := range companies {
		result := make(map[string]interface{})
		search_result := models_search.SearchResult{
			Symbol: value.Symbol,
			Name:   value.Name,
		}
		bytes, _ := json.Marshal(search_result)
		json.Unmarshal(bytes, &result)
		results = append(results, result)
	}
	return c.JSON(results)
}
