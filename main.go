package main

import (
	"DarwinScubaDiveBackend/ConectDB"
	"DarwinScubaDiveBackend/Routes"
	"DarwinScubaDiveBackend/Utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
)
var client *mongo.Client
func init (){
	client = ConectDB.ConectDB()
	Utils.ReservasCollection = client.Database("GaviotaFerry").Collection("Reservas")
	Utils.UsersCollection = client.Database("GaviotaFerry").Collection("Usuarios")
}
func main(){
	router := mux.NewRouter()
	router.HandleFunc("/", Routes.Index)
	router.HandleFunc(
		"/addUser/{FViaje}/{Ruta}/{Referencia}/{Proveedor}/{Cedula}/{Telefono}/{Status}/"+
			"{Nacionalidad}/{Observacion}/{FReserva}/{Edad}/{Precio}/{Pagado}",
		Routes.AddUser).Methods("POST")
	router.HandleFunc("/deleteUser/{Id}", Routes.DeleteUser).Methods("POST")
	router.HandleFunc(
		"/updateUser/{Id}/{FViaje}/{Ruta}/{Referencia}/{Proveedor}/{Cedula}/{Telefono}/{Status}/"+
			"{Nacionalidad}/{Observacion}/{FReserva}/{Edad}/{Precio}/{Pagado}",
		Routes.UpdateUser).Methods("POST")
	router.HandleFunc("/getDailyData/{Time}/{FViaje}",Routes.GetDailyData).Methods("GET")
	//credentials := handlers.AllowCredentials()
	router.HandleFunc("/report/{Inicio}/{Final}", Routes.Report).Methods("POST")
	router.HandleFunc("/reportUpdate/{Id}", Routes.UpdateFromReport).Methods("POST")
	methods := handlers.AllowedMethods([]string{"POST"})
	origins := handlers.AllowedOrigins([]string{"*"})
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, handlers.CORS(methods, origins)(router))
}


