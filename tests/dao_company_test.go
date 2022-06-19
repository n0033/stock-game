package tests

import (
	"testing"

	company_dao "github.com/ktylus/stock-game/common/dao/company"
	models_company "github.com/ktylus/stock-game/common/models/mongo/company"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CREATED_COMPANY_ID *primitive.ObjectID
var COMPANY_CODE string = "XYZ"
var dao_company = company_dao.NewDAOCompany()

func TestCreateCompany(t *testing.T) {
	company_create := models_company.CompanyInCreate{
		CompanyBase: models_company.CompanyBase{
			Name: COMPANY_CODE,
			Code: COMPANY_CODE,
		},
	}

	db_company, err := dao_company.Create(company_create)
	if err != nil {
		t.Fail()
	}

	CREATED_COMPANY_ID = &db_company.ID
}

func TestCreateExistingCompany(t *testing.T) {
	company_create := models_company.CompanyInCreate{
		CompanyBase: models_company.CompanyBase{
			Name: COMPANY_CODE,
			Code: COMPANY_CODE,
		},
	}
	_, err := dao_company.Create(company_create)
	if err == nil {
		t.Fail()
	}
}

func TestUpdateCompany(t *testing.T) {
	company_update := models_company.CompanyInUpdate{
		CompanyBase: models_company.CompanyBase{
			Name: "XYZ1",
			Code: COMPANY_CODE,
		},
	}

	_, err := dao_company.Update(*CREATED_COMPANY_ID, company_update)
	if err != nil {
		t.Fail()
	}
}

func TestFindOneCompany(t *testing.T) {
	_, err := dao_company.FindOne(*CREATED_COMPANY_ID)
	if err != nil {
		t.Fail()
	}
}

func TestFindByCodeCompany(t *testing.T) {
	_, err := dao_company.FindByCode(COMPANY_CODE)
	if err != nil {
		t.Fail()
	}
}

func TestDeleteCompany(t *testing.T) {
	_, err := dao_company.Delete(*CREATED_COMPANY_ID)
	if err != nil {
		t.Fail()
	}
}
