package ConectDB

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)
func ConectDB () *mongo.Client{

	//options := options.Client().ApplyURI("mongodb+srv://jhquihuiri7:Meta085216841@micluster.18zqq.mongodb.net/")
	options := options.Client().ApplyURI("mongodb+srv://doadmin:Z3d87ni4E91g05aX@logiciel-applab-dab57134.mongo.ondigitalocean.com/admin?authSource=admin&replicaSet=logiciel-applab&tls=true&tlsCAFile=ConectDB/ca-certificate.crt")
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
