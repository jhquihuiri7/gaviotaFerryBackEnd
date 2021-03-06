package Routes

import (
	"DarwinScubaDiveBackend/Utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)
type Pax struct {
	Id string `bson:"_id,omitempty"`
	Ruta string`bson:"ruta,omitempty"`
	Referencia string `bson:"referencia,omitempty"`
	Proveedor string `bson:"proveedor,omitempty"`
	Cedula string `bson:"cedula,omitempty"`
	Telefono string `bson:"telefono,omitempty"`
	Status string `bson:"status,omitempty"`
	Nacionalidad string `bson:"nacionalidad,omitempty"`
	Observacion string `bson:"observacion,omitempty"`
	FReserva string `bson:"freserva,omitempty"`
	FViaje string `bson:"fviaje,omitempty"`
	Edad int `bson:"edad,omitempty"`
	Precio  int `bson:"precio,omitempty"`
	Pagado bool `bson:"pagado,omitempty"`
}
type PaxFreq struct {
	Id string `bson:"_id,omitempty"`
	Referencia string `bson:"referencia,omitempty"`
	Cedula string `bson:"cedula,omitempty"`
	Telefono string `bson:"telefono,omitempty"`
	Status string `bson:"status,omitempty"`
	Nacionalidad string `bson:"nacionalidad,omitempty"`
	Edad int `bson:"edad,omitempty"`
}

func UpdateCapitanName(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	replacement := bson.D{{"capitanName", params["capName"]}}
	Utils.VariablesCollection.FindOneAndReplace(context.TODO(), bson.D{{"_id", "05b4c21b-21dd-405c-9c93-7e35346a6603"}}, replacement)
}
func GetCapitanName(w http.ResponseWriter, r *http.Request){
	var capitanName Utils.CapitanName
	opts := options.FindOne().SetProjection(bson.D{{"_id", 0}})
	err := Utils.VariablesCollection.FindOne(context.TODO(), bson.D{{"_id", "05b4c21b-21dd-405c-9c93-7e35346a6603"}},opts).Decode(&capitanName)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(capitanName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w,string(data))
}
func AddUser(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	fviaje := params["FViaje"]
	ruta := params["Ruta"]
	referencia := params["Referencia"]
	proveedor := params["Proveedor"]
	cedula := params["Cedula"]
	telefono := params["Telefono"]
	status := params["Status"]
	nacionalidad := params["Nacionalidad"]
	observacion := params["Observacion"]
	freserva := params["FReserva"]
	edad, _ := strconv.Atoi(params["Edad"])
	precio, _ := strconv.Atoi(params["Precio"])
	pagado, _ := strconv.ParseBool(params["Pagado"])
	newUser := bson.D{
		{"_id", uuid.NewV4().String()},
		{"FViaje",fviaje},
		{"Ruta", ruta},
		{"Referencia", Utils.ProcessNA(referencia)},
		{"Proveedor", Utils.ProcessNA(proveedor)},
		{"Cedula", Utils.ProcessNA(cedula)},
		{"Telefono", Utils.ProcessNA(telefono)},
		{"Status", status},
		{"Nacionalidad", Utils.ProcessNA(nacionalidad)},
		{"Observacion", Utils.ProcessNA(observacion)},
		{"FReserva", freserva},
		{"Edad", edad},
		{"Precio", precio},
		{"Pagado", pagado},
	}
	if strings.Contains(referencia,"@") {
		newFreqUser := bson.D{
			{"_id", uuid.NewV4().String()},
			{"Referencia", Utils.ProcessNA(referencia)},
			{"Cedula", Utils.ProcessNA(cedula)},
			{"Telefono", Utils.ProcessNA(telefono)},
			{"Status", status},
			{"Nacionalidad", Utils.ProcessNA(nacionalidad)},
			{"Edad", edad},
			{"DateRegister", time.Now().Year()},
		}
		Utils.UsersCollection.InsertOne(context.TODO(),newFreqUser)
	}else if strings.Contains(referencia,"*") {
		fmt.Println("AQUI")
		newFreqUser := bson.D{
			{"Referencia", Utils.ProcessNA(referencia)},
			{"Cedula", Utils.ProcessNA(cedula)},
			{"Telefono", Utils.ProcessNA(telefono)},
			{"Status", status},
			{"Nacionalidad", Utils.ProcessNA(nacionalidad)},
			{"Edad", edad},
			{"DateRegister", time.Now().Year()},
		}
		var replacedDocument bson.M
		opts := options.FindOneAndReplace().SetUpsert(true)
		filter := bson.D{{"Cedula", Utils.ProcessNA(cedula)}}
		replacement := newFreqUser
		err := Utils.UsersCollection.FindOneAndReplace(context.TODO(), filter, replacement, opts).Decode(&replacedDocument)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			log.Fatal(err)
		}
		fmt.Printf("replaced document %v", replacedDocument)
	}

	Utils.ReservasCollection.InsertOne(context.TODO(),newUser)

}
func Index (w http.ResponseWriter, r *http.Request){
	var paxs []Pax
	cur, currErr := Utils.ReservasCollection.Find(context.TODO(), bson.D{})
	if currErr != nil {
		log.Fatal(currErr)
	}
	for cur.Next(context.TODO()) {
		var pax Pax
		err := cur.Decode(&pax)
		if err != nil {
			log.Fatal(err)
		}
		paxs = append(paxs, pax)
	}
	data, _ := json.Marshal(paxs)
	fmt.Fprintln(w,string(data))
}
func GetDailyData(w http.ResponseWriter, r *http.Request){
	var paxs []Pax
	var dailyData []Pax
	cur, currErr := Utils.ReservasCollection.Find(context.TODO(), bson.D{})
	if currErr != nil {
		log.Fatal(currErr)
	}
	for cur.Next(context.TODO()) {
		var pax Pax
		err := cur.Decode(&pax)
		if err != nil {
			log.Fatal(err)
		}
		paxs = append(paxs, pax)

	}
	params := mux.Vars(r)
	time := params["Time"]
	date:= params["FViaje"]

	for _, v := range paxs {
		if v.FViaje == date {
			if time == "Todo" {
				dailyData = append(dailyData, v)
			}else if time == "AM" {
				if v.Ruta == "SC-SX" {
					dailyData = append(dailyData, v)
				}
			}else if time == "PM" {
				if v.Ruta == "SX-SC" {
					dailyData = append(dailyData, v)
				}
			}
		}
	}
	data, _ := json.Marshal(dailyData)
	fmt.Fprintln(w, string(data))
}
func DeleteUser(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, _ := params["Id"]
	Utils.ReservasCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	fmt.Fprintln(w, "OK")
}
func UpdateUser (w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id := params["Id"]
	fviaje := params["FViaje"]
	//ruta := params["Ruta"]
	referencia := params["Referencia"]
	proveedor := params["Proveedor"]
	cedula := params["Cedula"]
	telefono := params["Telefono"]
	status := params["Status"]
	nacionalidad := params["Nacionalidad"]
	observacion := params["Observacion"]
	//freserva := params["FReserva"]
	edad, _ := strconv.Atoi(params["Edad"])
	precio, _ := strconv.Atoi(params["Precio"])
	pagado, _ := strconv.ParseBool(params["Pagado"])
	data := bson.D{
		{"Referencia", Utils.ProcessNA(referencia)},
		{"Proveedor", Utils.ProcessNA(proveedor)},
		{"Cedula", Utils.ProcessNA(cedula)},
		{"Telefono", Utils.ProcessNA(telefono)},
		{"Status", status},
		{"Nacionalidad", Utils.ProcessNA(nacionalidad)},
		{"Observacion", Utils.ProcessNA(observacion)},
		{"Edad", edad},
		{"Precio", precio},
		{"Pagado", pagado},
		{"FViaje", fviaje},
	}
	Utils.ReservasCollection.UpdateOne(context.TODO(),bson.M{"_id": id},bson.D{{"$set", data}})
}