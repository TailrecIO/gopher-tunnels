package commons

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/xid"
	"github.com/tailrecio/gopher-tunnels/config"
	"sync"
)

var dbSessionMu sync.Mutex
var dbSession *dynamodb.DynamoDB

var gopherMu sync.Mutex
var gopherMap = make(map[string]*Gopher)

func getTableKey(id *string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{ "id": { S: id,}}
}

func getDbSession() *dynamodb.DynamoDB {
	if dbSession == nil {
		dbSessionMu.Lock()
		defer dbSessionMu.Unlock()
		if dbSession == nil {
			dbSession = dynamodb.New(NewAwsSession())
		}
	}
	return dbSession
}

func NewGopher(publicKey *string, mode string) (*Gopher, error) {
	if mode != ModeAsync && mode != ModeSync {
		return nil, errors.New("invalid mode")
	}
	id := xid.New().String()
	requestQueueName := MakeQueueName("in")
	client := Gopher{
		Id:               &id,
		EncodedPublicKey: publicKey,
		Mode:             mode,
		RequestQueueName: &requestQueueName,
	}
	item, err := dynamodbattribute.MarshalMap(client)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(config.GetAwsDynamoDbTable()),
		Item:      item,
	}
	_, err = getDbSession().PutItem(input)
	return &client, err
}

// Get gopher info from DynamoDB and cache the result by queueName.
func GetGopher(id *string) (*Gopher, error) {
	var gopher = gopherMap[*id]
	if gopher == nil {
		gopherMu.Lock()
		defer gopherMu.Unlock()
		if gopherMap[*id] == nil {
			input := &dynamodb.GetItemInput{
				TableName: aws.String(config.GetAwsDynamoDbTable()),
				Key:       getTableKey(id),
			}
			out, err := getDbSession().GetItem(input)
			dynamodbattribute.UnmarshalMap(out.Item, &gopher)
			if err == nil {
				// mutate map when err is nil
				gopherMap[*id] = gopher
			}
			return gopher, err
		}
	}
	return gopher, nil
}

func KillGopher(id *string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(config.GetAwsDynamoDbTable()),
		Key:       getTableKey(id),
	}
	_, err := getDbSession().DeleteItem(input)
	return err
}
