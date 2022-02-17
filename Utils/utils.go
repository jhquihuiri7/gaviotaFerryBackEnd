package Utils

import (
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)
var ReservasCollection *mongo.Collection
var UsersCollection *mongo.Collection
func ProcessNA(value string) string {
	if value == "NA" || strings.ToUpper(value) == "BORRAR"{
		return ""
	}else {
		return strings.ReplaceAll(strings.ToUpper(value),"@","")
	}
}
func GetColleciton(collection *mongo.Collection) *mongo.Collection{
	return collection
}
