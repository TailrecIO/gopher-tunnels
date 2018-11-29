package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tailrecio/gopher-tunnels/commons"
	"log"
)

var respondHandler = func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var err error

	id := req.PathParameters["id"]
	if id == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "id is missing"}, nil
	}

	var whResponse commons.WebhookResponse
	err = json.Unmarshal([]byte(req.Body), &whResponse)
	if err != nil {
		log.Printf("Couldn't unmarshal a response object due to: %v\n", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, err
	}

	var gopher *commons.Gopher
	gopher, err = commons.GetGopher(&id)
	if err != nil {
		log.Printf("Couldn't find Gopher: %v\n", id)
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "gopher not found"}, err
	}

	err = commons.SendResponse(gopher, &whResponse)

	if err != nil {
		log.Printf("Failed to send a response back to Gopher: `%v` due to %v\n", id, err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "internal server error"}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: ""}, err
}

func main() {

	lambda.Start(respondHandler)

}
