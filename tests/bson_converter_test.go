package tests

import (
	"testing"

	convert "github.com/n0033/stock-game/services/bson_converter"
	"go.mongodb.org/mongo-driver/bson"
)

type TestStruct struct {
	Field_one string  `json:"field_one" bson:"field_one"`
	Field_two float64 `json:"field_two" bson:"field_two"`
}

func TestConvertStructToBsonM(t *testing.T) {
	obj := TestStruct{
		Field_one: "test",
		Field_two: 99.9,
	}
	bson_obj := convert.StructToBsonM(obj)

	_, exists := bson_obj["field_one"]
	if !exists {
		t.Fail()
	}

	_, exists = bson_obj["field_two"]

	if !exists {
		t.Fail()
	}

}

func TestConvertBsonDToStruct(t *testing.T) {
	bson_obj := bson.D{
		{"field_one", "test"},
		{"field_two", 9.99},
	}
	obj := convert.BsonDToStruct[TestStruct](bson_obj)

	if obj.Field_one != "test" {
		t.Fail()
	}

	if obj.Field_two != 9.99 {
		t.Fail()
	}
}
