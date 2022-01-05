package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
)

type Pax struct {
	Id, Ruta, Referencia, Proveedor, Cedula, Telefono, Status, Nacionalidad, Observacion, FReserva, FViaje string
	Edad, Precio  int
	Pagado bool
}
var Paxs []Pax
var client *mongo.Client
var collection *mongo.Collection
func init (){
	client = ConectDB()
	collection = client.Database("GaviotaFerry").Collection("Reservas")
}
func main(){
	//Paxs = append(Paxs, Pax{Id:uuid.NewV4().String()})
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc(
		"/addUser/{FViaje}/{Ruta}/{Referencia}/{Proveedor}/{Cedula}/{Telefono}/{Status}/"+
			"{Nacionalidad}/{Observacion}/{FReserva}/{Edad}/{Precio}/{Pagado}",
		addUser).Methods("POST")
	router.HandleFunc("/deleteUser/{Id}", deleteUser).Methods("POST")
	router.HandleFunc(
		"/updateUser/{Id}/{FViaje}/{Ruta}/{Referencia}/{Proveedor}/{Cedula}/{Telefono}/{Status}/"+
			"{Nacionalidad}/{Observacion}/{FReserva}/{Edad}/{Precio}/{Pagado}",
		updateUser).Methods("POST")
	router.HandleFunc("/getDailyData/{Time}/{FViaje}",getDailyData).Methods("GET")
	//credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST"})
	origins := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":8080", handlers.CORS(methods, origins)(router))
}
func index (w http.ResponseWriter, r *http.Request){
	data, _ := json.Marshal(Paxs)
	fmt.Fprintln(w,string(data))
}
func addUser(w http.ResponseWriter, r *http.Request){
	//{Ruta}/{Referencia}/{Proveedor}/{Cedula}/{Telefono}/{Status}/"+
	//			"{Nacionalidad}/{Observacion}/{FReserva}/{Edad}/{Precio}/{Pagado}
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
	}
	collection.InsertOne(context.TODO(),newUser)
	Paxs = append(Paxs,
		Pax{
		Id:uuid.NewV4().String(),
		FViaje: fviaje,
		Ruta: ruta,
		Referencia: referencia,
		Proveedor: proveedor,
		Cedula: cedula,
		Telefono: telefono,
		Status: status,
		Nacionalidad: nacionalidad,
		Observacion: observacion,
		FReserva: freserva,
		Edad: edad,
		Precio: precio,
		Pagado: pagado,
		})
}
func getDailyData (w http.ResponseWriter, r *http.Request){
	var dailyData []Pax
	params := mux.Vars(r)
	time := params["Time"]
	date:= params["FViaje"]

	for _, v := range Paxs {
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
func deleteUser (w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, _ := params["Id"]
	for i, v := range Paxs {
		if v.Id == id {
			Paxs = append(Paxs[:i], Paxs[i+1:]...)
		}
	}
	data, _ := json.Marshal(Paxs)
	fmt.Fprintln(w, string(data))
}
func updateUser (w http.ResponseWriter, r *http.Request){
	indexUser := 0
	params := mux.Vars(r)
	id := params["Id"]
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
	for i, v := range Paxs {
		if v.Id == id {
			indexUser = i
			break
		}
	}
	Paxs[indexUser] = Pax{
		Id:id,
		FViaje: fviaje,
		Ruta: ruta,
		Referencia: referencia,
		Proveedor: proveedor,
		Cedula: cedula,
		Telefono: telefono,
		Status: status,
		Nacionalidad: nacionalidad,
		Observacion: observacion,
		FReserva: freserva,
		Edad: edad,
		Precio: precio,
		Pagado: pagado,
	}

}
func ConectDB () *mongo.Client{
	//options := options.Client().ApplyURI("mongodb+srv://jhquihuiri7:Meta085216841@micluster.18zqq.mongodb.net/")
	options := options.Client().ApplyURI("mongodb+srv://doadmin:Z3d87ni4E91g05aX@logiciel-applab-dab57134.mongo.ondigitalocean.com/admin?authSource=admin&replicaSet=logiciel-applab&tls=true&tlsCAFile=ca-certificate.crt")
	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return  client
}