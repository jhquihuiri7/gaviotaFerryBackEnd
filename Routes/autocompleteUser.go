package Routes

import (
	"DarwinScubaDiveBackend/Utils"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type AutocompleteUser struct {
	Referencia string `json:"Referencia, omitempty"`
	Cedula string `bson:"Cedula,omitempty"`
	Status string `json:"Status, omitempty"`
	Nacionalidad string `json:"Nacionalidad, omitempty"`
	Edad int `bson:"Edad,omitempty"`
	DateRegister int `bson:"DateRegister,omitempty"`
}

func GetAutocompleteUser(w http.ResponseWriter, r *http.Request) {
	cur, currErr := Utils.UsersCollection.Find(context.TODO(),bson.D{})
	if currErr != nil {
		fmt.Println(currErr)
	}
	var autocompleteUsers []AutocompleteUser
	for cur.Next(context.TODO()) {
		var autocompleteUser AutocompleteUser
		err := cur.Decode(&autocompleteUser)
		if err != nil {
			fmt.Println(err)
		}
		autocompleteUsers = append(autocompleteUsers, autocompleteUser)
	}
	data,err := json.Marshal(autocompleteUsers)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(data))
}
