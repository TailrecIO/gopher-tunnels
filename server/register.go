package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tailrecio/gopher-tunnels/commons"
)

var registerHandler = func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error

	var whRegister commons.WebhookRegister
	err = json.Unmarshal([]byte(req.Body), &whRegister)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, err
	}
	var gopher *commons.Gopher
	gopher, err = commons.NewGopher(whRegister.EncodedPublicKey, whRegister.Mode)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}
	_, err = commons.CreateRequestQueue(gopher)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	var outputBytes []byte
	outputBytes, err = json.Marshal(gopher)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(outputBytes)}, nil
}

func main() {

	lambda.Start(registerHandler)
}
