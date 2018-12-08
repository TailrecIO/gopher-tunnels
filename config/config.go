package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Stage = "stage"
	AccountId = "accountId"
	Region = "region"
	QueueType = "queueType"
	DynamoDbTable = "dynamoDbTable"
	BaseQueueEndpoint = "baseQueueEndpoint"
	BaseApiEndpoint = "baseApiEndpoint"

	Mode = "mode"
	TargetHost = "targetHost"
	TargetPort = "targetPort"

)

var stage string
func GetStage() string {
	if stage == "" {
		stage = os.Getenv(Stage)
		log.Printf("env.Stage: %v\n", stage)
	}
	return stage
}

var awsAccount string
func GetAwsAccount() string {
	if awsAccount == "" {
		awsAccount = os.Getenv(AccountId)
		log.Printf("env.AccountId: %v\n", awsAccount)
	}
	return awsAccount
}

var awsRegion string
func GetAwsRegion() string {
	if awsRegion == "" {
		awsRegion = os.Getenv(Region)
		log.Printf("env.Region: %v\n", awsRegion)
	}
	return awsRegion
}

var queueType string
func GetQueueType() string {
	if queueType == "" {
		queueType = os.Getenv(QueueType)
		log.Printf("env.QueueType: %v\n", queueType)
	}
	return queueType
}

var awsDynamoDbTable string
func GetAwsDynamoDbTable() string {
	if awsDynamoDbTable == "" {
		awsDynamoDbTable = os.Getenv(DynamoDbTable)
		log.Printf("env.DynamoDbTable: %v\n", awsDynamoDbTable)
	}
	return awsDynamoDbTable
}

var baseQueueEndpoint string
func GetBaseQueueEndpoint() string {
	if baseQueueEndpoint == "" {
		baseQueueEndpoint = os.Getenv(BaseQueueEndpoint)
		if strings.HasSuffix(baseQueueEndpoint, "/") {
			baseQueueEndpoint = strings.TrimSuffix(baseQueueEndpoint, "/")
		}
		log.Printf("env.BaseQueueEndpoint: %v\n", baseQueueEndpoint)
	}
	return baseQueueEndpoint
}

var baseApiEndpoint string
func GetBaseApiEndpoint() string {
	if baseApiEndpoint == "" {
		baseApiEndpoint = os.Getenv(BaseApiEndpoint)
		if strings.HasSuffix(baseApiEndpoint, "/") {
			baseApiEndpoint = strings.TrimSuffix(baseApiEndpoint, "/")
		}
		log.Printf("env.BaseApiEndpoint: %v\n", baseApiEndpoint)
	}
	return baseApiEndpoint
}

var mode string
func GetMode() string {
	if mode == "" {
		mode = os.Getenv(Mode)
		log.Printf("env.Mode: %v\n", mode)
	}
	return mode
}

var targetHost string
func GetTargetHost() string {
	if targetHost == "" {
		targetHost = os.Getenv(TargetHost)
		log.Printf("env.TargetHost: %v\n", targetHost)
	}
	return targetHost
}

var targetPort int
func GetTargetPort() int {
	if targetPort == 0 {
		var err error
		targetPort, err = strconv.Atoi(os.Getenv(TargetPort))
		if err != nil {
			log.Fatalf("Invalid %v: %v", TargetPort, os.Getenv(TargetPort))
		}
		log.Printf("env.TargetPort: %v\n", targetPort)
	}
	return targetPort
}

