package stock_datapoint

import (
	"time"

	models_av "github.com/ktylus/stock-game/common/models/resources/alpha_vantage"
	"github.com/ktylus/stock-game/common/models/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockDatapointValue struct {
	Open  float64 `bson:"open"`
	High  float64 `bson:"high"`
	Low   float64 `bson:"low"`
	Close float64 `bson:"close"`
}

type StockDatapointBase struct {
	Code      string              `bson:"code"`
	Timestamp time.Time           `bson:"timestamp"`
	Volume    float64             `bson:"volume"`
	Value     StockDatapointValue `bson:"value"`
	Last_used time.Time           `bson:"last_used"`
}
type StockDatapointInDB struct {
	StockDatapointBase `bson:",inline"`
	utils.DBModelMixin `bson:",inline"`
}

func (dp_value *StockDatapointInDB) GetValue() float64 {
	return dp_value.Value.Close
}

type StockDatapointInResponse struct {
	StockDatapointBase `bson:",inline" json:",inline"`
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
}

func (resp *StockDatapointInResponse) FromStockDatapointInDB(db_datapoint StockDatapointInDB) StockDatapointInResponse {
	return StockDatapointInResponse{
		StockDatapointBase: StockDatapointBase{
			Code:      db_datapoint.Code,
			Timestamp: db_datapoint.Timestamp,
			Volume:    db_datapoint.Volume,
			Value:     db_datapoint.Value,
			Last_used: time.Now().UTC(),
		},
		ID: db_datapoint.ID,
	}
}

type StockDatapointInUpdate struct {
	StockDatapointBase `bson:",inline"`
}

type StockDatapointInCreate struct {
	StockDatapointBase `bson:",inline"`
}

func (datapoint *StockDatapointInCreate) FromAVStockDatapoint(av_datapoint models_av.StockDatapoint) StockDatapointInCreate {
	return StockDatapointInCreate{
		StockDatapointBase: StockDatapointBase{
			Code:      av_datapoint.Symbol,
			Timestamp: av_datapoint.Timestamp,
			Volume:    float64(av_datapoint.Volume),
			Value: StockDatapointValue{
				Open:  av_datapoint.Values.Open,
				High:  av_datapoint.Values.High,
				Low:   av_datapoint.Values.Low,
				Close: av_datapoint.Values.Close,
			},
			Last_used: time.Now().UTC(),
		},
	}
}
