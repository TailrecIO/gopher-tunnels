package commons

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/tailrecio/gopher-tunnels/config"
	"os"
)

const AnonymousCredentials = "AnonymousCredentials"

func NewAwsSession() *s.Session {
	awsConfig := aws.Config{
		Region: aws.String(config.GetAwsRegion()),
		CredentialsChainVerboseErrors: aws.Bool(true),
	}
	if os.Getenv(AnonymousCredentials) == "true" {
		awsConfig.Credentials = credentials.AnonymousCredentials
	}
	sess, err := s.NewSession(&awsConfig)
	if err != nil {
		panic(err.Error())
	}
	return sess
}
