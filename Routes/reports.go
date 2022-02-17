package Routes

import (
	"DarwinScubaDiveBackend/Utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strconv"
)

type ReportData struct {
	Ventas int
	Facturado int
	Recuperado int
	Recuperar int
	Detalle []DetalleProveedor
}
type DetalleProveedor struct {
	Proveedor string
	DetallePasajero []DetallePasajero
	Total int
}
type DetallePasajero struct {
	Id string
	FViaje string
	Ruta string
	Referencia string
	Precio int
}
func Report( w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	finicio, err := strconv.Atoi(params["Inicio"])
	ffinal, err := strconv.Atoi(params["Final"])
	var paxs []Pax
	var paxsCobrar []Pax
	var reportData ReportData
	var detalleProveedores []DetalleProveedor
	m := make(map[string]int)
	cur, curErr := Utils.ReservasCollection.Find(context.TODO(), bson.M{})
	if curErr != nil {
		log.Fatal(curErr)
	}
	for cur.Next(context.TODO()){
		var pax Pax
		err := cur.Decode(&pax)
		if err != nil {
			log.Fatal(err)
		}
		fviaje, err := strconv.Atoi(pax.FViaje)
 		if err != nil {
			log.Fatal(err)
		}
		if fviaje >= finicio && fviaje <= ffinal {
			reportData.Facturado += pax.Precio
			if pax.Pagado == true {
				reportData.Recuperado += pax.Precio
			}else {
				reportData.Recuperar += pax.Precio
				m[pax.Proveedor] += 1
				paxsCobrar = append(paxsCobrar, pax)
			}
			paxs = append(paxs, pax)
		}
	}
	reportData.Ventas = len(paxs)

	for i,_ := range m {
		var detalleProveedor DetalleProveedor
		detalleProveedor.Proveedor = i
		var detallePasajeros []DetallePasajero
		for _,v := range paxsCobrar {
			if v.Proveedor == i {
				var detallePasajero DetallePasajero
				detalleProveedor.Total += v.Precio
				detallePasajero.Id = v.Id
				detallePasajero.FViaje = v.FViaje
				detallePasajero.Ruta = v.Ruta
				detallePasajero.Precio = v.Precio
				detallePasajero.Referencia = v.Referencia
				detallePasajeros = append(detallePasajeros, detallePasajero)
				detalleProveedor.DetallePasajero = detallePasajeros
			}
		}
		detalleProveedores = append(detalleProveedores, detalleProveedor)
	}
	reportData.Detalle = detalleProveedores
	data, err := json.Marshal(reportData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w,string(data))
}
func UpdateFromReport(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id := params["Id"]
	fmt.Println("HASTA QUI")
	fmt.Println("BIEN")
	result, err := Utils.ReservasCollection.UpdateOne(context.TODO(),bson.D{{"_id", id}},bson.D{{"$set", bson.D{{"Pagado", true}}}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.MatchedCount)
	fmt.Println(result.ModifiedCount)
}