package db

import (
	"context"
	"fmt"

	"github.com/rogercruzvillca/ssvv-cloud/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDN *mongo.Client
var DataBaseName string

func ConectionDB(ctx context.Context) error {
	user := ctx.Value(models.Key("user")).(string)
	password := ctx.Value(models.Key("password")).(string)
	host := ctx.Value(models.Key("host")).(string)
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, password, host)

	var ClientOptions = options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(ctx, ClientOptions)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Conexion exitosa con DB")
	MongoDN = client
	DataBaseName = ctx.Value(models.Key("database")).(string)
	return nil
}

func VerificarConexionDB() bool {
	err := MongoDN.Ping(context.TODO(), nil)
	return err == nil
}
