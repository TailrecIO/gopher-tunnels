package commons

import (
	"github.com/aws/aws-sdk-go/aws"
	s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/tailrecio/gopher-tunnels/config"
)

func NewAwsSession() *s.Session {

	sess, err := s.NewSession(&aws.Config{
		Region: aws.String(config.GetAwsRegion())},
	)
	if err != nil {
		panic(err.Error())
	}
	return sess
}
