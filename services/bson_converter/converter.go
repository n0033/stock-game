package converter

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func StructToBsonD[T any](source T) bson.D {
	doc, err := bson.Marshal(source)
	if err != nil {
		message := fmt.Sprintf("Cannot marshal object of type `%T`.", source)
		log.Fatal(message)
	}
	var bson_doc bson.D
	err = bson.Unmarshal(doc, &bson_doc)
	if err != nil {
		message := fmt.Sprintf("Cannot unmarshal object of type `%T`.", source)
		log.Fatal(message)
	}
	return bson_doc
}

func StructToBsonM[T any](source T) bson.M {
	doc, err := bson.Marshal(source)
	if err != nil {
		message := fmt.Sprintf("Cannot marshal object of type `%T`.", source)
		log.Fatal(message)
	}
	var bson_doc bson.M
	err = bson.Unmarshal(doc, &bson_doc)
	if err != nil {
		message := fmt.Sprintf("Cannot unmarshal object of type `%T`.", source)
		log.Fatal(message)
	}
	return bson_doc
}

func BsonDToStruct[T any](source bson.D) T {
	var obj T
	doc_bytes, err := bson.Marshal(source)
	if err != nil {
		message := fmt.Sprintf("Cannot marshal object of type `%T`.", source)
		log.Fatal(message)
	}
	bson.Unmarshal(doc_bytes, &obj)
	if err != nil {
		message := fmt.Sprintf("Cannot unmarshal object of type `%T`.", source)
		log.Fatal(message)
	}
	return obj
}
