package aws_utils_go

// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SecretsManager is our wrapper around AWS SecretsManager
type SecretsManager struct {
	Svc *secretsmanager.SecretsManager
}

// GetSecretsManagerFromEnv constructs SecretsManager using environment variables
func GetSecretsManagerFromEnv() SecretsManager {
	region := core_utils_go.NonEmptyEnv("AWS_REGION")
	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
                                  aws.NewConfig().WithRegion(region))
	return SecretsManager{
		Svc: svc,
	}
}

// ReadSecretString reads the secret
func (s SecretsManager)ReadSecretString(secretName string) (secretString string, err error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	var result *secretsmanager.GetSecretValueOutput
	result, err = s.Svc.GetSecretValue(input)
	if err == nil {
		if result.SecretString != nil {
			secretString = *result.SecretString
		} else {
			err = errors.New("Couldn't find secret " + secretName + " in region " + core_utils_go.NonEmptyEnv("AWS_REGION"))
		}
	}
	return
}
