package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tailrecio/gopher-tunnels/commons"
	"strings"
)

var responseQueueName = commons.MakeQueueName("out")

var webhookHandler = func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	id := req.PathParameters["id"]
	if id == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "id is missing"}, nil
	}

	gopher, err := commons.GetGopher(&id)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	context := commons.WebhookRequestContext{
		ResponseQueueName: &responseQueueName,
	}

	// s = "/webhook/{id}/path1/path2/"
	// parts[0] = "" (an empty string)
	// parts[1] = "webhook"
	// parts[2] = "{id}"
	// parts[3] = "path1/path2/" (the rest of path string without / prefix)
	parts := strings.SplitN(req.Path, "/", 4)
	proxyPath := "/" + parts[3]

	whReq := commons.WebhookRequest{
		Context:     &context,
		Path:        &proxyPath,
		QueryParams: req.QueryStringParameters,
		Method:      &req.HTTPMethod,
		Headers:     req.Headers,
		Body:        &req.Body,
	}

	err = commons.SendRequest(gopher, &whReq)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	if gopher.Mode == commons.ModeSync {
		var whRes *commons.WebhookResponse
		whRes, err = commons.ReadResponse(&responseQueueName)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: whRes.StatusCode,
			Headers:    whRes.Headers,
			Body:       *whRes.Body,
		}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "OK"}, nil
}

func main() {

	lambda.Start(webhookHandler)
}
