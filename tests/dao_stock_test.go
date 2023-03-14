package tests

import (
	"testing"
	"time"

	stock_dao "github.com/n0033/stock-game/common/dao/stock"
)

var dao_stock = stock_dao.NewDAOStock()

func TestFind300Latest(t *testing.T) {
	datapoints := dao_stock.Find300LatestByCode(APPLE_CODE)
	if len(datapoints) == 0 || len(datapoints) > 300 {
		t.Fail()
	}
}

func TestLatest(t *testing.T) {
	datapoint, err := dao_stock.FindLatestByCode(APPLE_CODE)
	if err != nil {
		t.Fail()
	}

	if datapoint.Timestamp.Before(time.Now().AddDate(0, 0, -4)) {
		t.Fail()
	}
}
