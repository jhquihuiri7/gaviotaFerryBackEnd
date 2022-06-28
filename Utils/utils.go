package Utils

import (
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)
var ReservasCollection *mongo.Collection
var UsersCollection *mongo.Collection
var VariablesCollection *mongo.Collection
func ProcessNA(value string) string {
	if value == "NA" || strings.ToUpper(value) == "BORRAR"{
		return ""
	}else {
		r := strings.NewReplacer("@", "", "*", "")
		return r.Replace(strings.ToUpper(value))
	}
}
func GetColleciton(collection *mongo.Collection) *mongo.Collection{
	return collection
}
