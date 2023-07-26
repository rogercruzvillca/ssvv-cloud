package secretmanager

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rogercruzvillca/ssvv-cloud/awsgo"
	"github.com/rogercruzvillca/ssvv-cloud/models"
)

func GetSecret(secretName string) (models.Secret, error) {
	var dataSecret models.Secret
	fmt.Println("> Obteniendo Secreto " + secretName)

	svc := secretsmanager.NewFromConfig(awsgo.Cfg)
	key, err := svc.GetSecretValue(awsgo.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		fmt.Println(err.Error())
		return dataSecret, err
	}
	json.Unmarshal([]byte(*key.SecretString), &dataSecret)
	fmt.Println("> Lectura de Secret OK " + secretName)
	return dataSecret, nil
}
