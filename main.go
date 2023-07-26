package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/rogercruzvillca/ssvv-cloud/awsgo"
	"github.com/rogercruzvillca/ssvv-cloud/db"
	"github.com/rogercruzvillca/ssvv-cloud/handlers"
	"github.com/rogercruzvillca/ssvv-cloud/models"
	"github.com/rogercruzvillca/ssvv-cloud/secretmanager"
)

func main() {
	lambda.Start(EjecutarLambda)
}

func EjecutarLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var response *events.APIGatewayProxyResponse
	awsgo.InicializarAWS()
	if !ValidarParametros() {
		response = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "error",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}
		return response, nil
	}
	secretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		response = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error al validar Secret",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}
		return response, nil
	}
	path := strings.Replace(request.PathParameters["ssvvcloud"], os.Getenv("UrlPrefix"), "", -1)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), secretModel.UserName)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), secretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), secretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), secretModel.DataBase)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), secretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	//Check connect to DB
	err = db.ConectionDB(awsgo.Ctx)
	if err != nil {
		response = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la conexion con la DB " + err.Error(),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}
		return response, nil
	}

	responseAPI := handlers.Manejadores(awsgo.Ctx, request)
	if responseAPI.Response == nil {
		response = &events.APIGatewayProxyResponse{
			StatusCode: responseAPI.Status,
			Body:       responseAPI.Message,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}
		return response, nil
	}
	return responseAPI.Response, nil
}

func ValidarParametros() bool {
	_, parametro := os.LookupEnv("SecretName")
	if !parametro {
		return parametro
	}
	_, parametro = os.LookupEnv("BucketName")
	if !parametro {
		return parametro
	}
	_, parametro = os.LookupEnv("UrlPrefix")
	if !parametro {
		return parametro
	}
	return parametro
}
