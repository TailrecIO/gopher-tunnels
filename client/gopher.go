package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tailrecio/gopher-tunnels/commons"
	"github.com/tailrecio/gopher-tunnels/config"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)


func serviceEndpoint(path string) string {
	return config.GetBaseApiEndpoint() + path
}

func register() (*commons.Gopher, error) {

	whRegister := commons.WebhookRegister{
		EncodedPublicKey: commons.GetKeyPair().GetHexEncodedPublicKey(),
		Mode:             config.GetMode(),
	}
	resp, err := resty.R().
		SetBody(whRegister).
		Post(serviceEndpoint("/register"))
	if err != nil {
		return nil, err
	}
	if resp.RawResponse.StatusCode != 200 {
		return nil, errors.New(string(resp.Body()))
	}

	var gopher commons.Gopher
	err = json.Unmarshal(resp.Body(), &gopher)
	if err != nil {
		return nil, err
	}
	return &gopher, nil
}

func executeRequest(gopher *commons.Gopher, req *commons.WebhookRequest) {
	fmt.Printf("Request: %v\n", *req)
	var res *commons.WebhookResponse
	if req.Context.Error != nil {
		responseCtx := commons.WebhookResponseContext{
			ResponseQueueName:    req.Context.ResponseQueueName,
			RequestMessageId:     req.Context.MessageId,
			RequestReceiptHandle: req.Context.ReceiptHandle,
		}
		res = &commons.WebhookResponse{
			Context:    &responseCtx,
			StatusCode: 500,
		}
	} else {
		res = forwardRequest(req)
	}
	fmt.Printf("Response: %v\n", *res)
	resp, err := resty.R().
		SetBody(res).
		Post(serviceEndpoint("/respond/" + *gopher.Id))
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	} else {
		if resp.RawResponse.StatusCode != 200 {
			fmt.Printf("Error: %v\n", string(resp.Body()))
		}
	}
}

func forwardRequest(req *commons.WebhookRequest) *commons.WebhookResponse {

	var proxyReq *http.Request
	// TODO: support HTTPS?
	// TODO: support custom host

	targetUrl := createUrlFromWebhookRequest(req)
	log.Printf("Forwarding a request to %v\n", targetUrl)

	responseCtx := commons.WebhookResponseContext{
		ResponseQueueName:    req.Context.ResponseQueueName,
		RequestMessageId:     req.Context.MessageId,
		RequestReceiptHandle: req.Context.ReceiptHandle,
	}
	fmt.Printf("Request: %v\n", req)
	fmt.Printf("Method: %v\n", *req.Method)
	fmt.Printf("Target URL: %v\n", targetUrl.String())
	fmt.Printf("Body: %v\n", string(*req.Body))
	proxyReq, err := http.NewRequest(*req.Method, targetUrl.String(), strings.NewReader(*req.Body))
	if err != nil {
		return commons.ErrorResponse(err, &responseCtx)
	}
	for header, value := range req.Headers {
		proxyReq.Header.Add(header, value)
	}

	client := &http.Client{}
	var proxyRes *http.Response
	var headerRes = make(map[string]string)
	proxyRes, err = client.Do(proxyReq)
	if err != nil {
		return commons.ErrorResponse(err, &responseCtx)
	}
	for header, value := range proxyRes.Header {
		headerRes[header] = value[0]
	}
	var resBytes []byte
	resBytes, err = ioutil.ReadAll(proxyRes.Body)
	if err != nil {
		return commons.ErrorResponse(err, &responseCtx)
	} else {
		resBody := string(resBytes)
		return &commons.WebhookResponse{
			Context:    &responseCtx,
			StatusCode: proxyRes.StatusCode,
			Headers:    headerRes,
			Body:       &resBody,
		}
	}
}

func listen(gopher *commons.Gopher) {

	for {
		log.Println("Waiting for incoming messages...")
		requests, err := commons.ReadRequests(gopher)
		if err != nil {
			panic(err)
		}
		log.Printf("Received %v messages\n", len(requests))
		for _, req := range requests {
			go executeRequest(gopher, req)
		}
	}
}

func loadConfiguration(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Fatalf("The config file: `%v` does not exist!", configFile)
	}
	clientConfig := make(map[string]string)
	ymlData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load the config file: `%v` due to `%v`", configFile, err.Error())
	}
	err = yaml.Unmarshal(ymlData, &clientConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML from a file due to `%v`", err.Error())
	}
	for k, v := range clientConfig {
		if os.Getenv(k) == "" {
			err = os.Setenv(k, v)
			if err != nil {
				log.Fatalf("Failed to set environment name: `%v` with value: `%v` due to `%v`", k, v, err.Error())
			}
		}
	}
}

func setCommandLineOptions() {
	var mode string
	var targetHost string
	var targetPort string
	flag.StringVar(&mode, "mode", "", "webhook mode: sync, async")
	flag.StringVar(&targetHost, "host", "", "target host")
	flag.StringVar(&targetPort, "port", "", "target port")
	flag.Parse()
	if mode != "" {
		os.Setenv(config.Mode, mode)
	}
	if targetHost != "" {
		os.Setenv(config.TargetHost, targetHost)
	}
	if targetPort != "" {
		os.Setenv(config.TargetPort, targetPort)
	}
}

func createQueryStringFromMap(paramMap map[string]string) *string {
	if paramMap == nil {
		return nil
	}
	queryString := ""
	if len(paramMap) > 0 {
		params := make([]string, len(paramMap))
		for k, v := range paramMap {
			param := url.QueryEscape(k) + "="
			if v != "" {
				param += url.QueryEscape(v)
			}
			params = append(params, param)
		}
		queryString = strings.Join(params, "&")
	}
	return &queryString
}

func createUrlFromWebhookRequest(req *commons.WebhookRequest) *url.URL {
	queryString := createQueryStringFromMap(req.QueryParams)
	targetUrl := url.URL{
		Scheme:   "http", // TODO: support HTTPS scheme
		Host:     fmt.Sprintf("%v:%v", config.GetTargetHost(), config.GetTargetPort()),
	}
	if req.Path != nil {
		targetUrl.Path = *req.Path
	}
	if queryString != nil {
		targetUrl.RawQuery = *queryString
	}
	return &targetUrl
}

func parseSimpleQuery(encodedQuery *string) (map[string]string, error) {
	values, err := url.ParseQuery(*encodedQuery)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			// take only the head
			m[k] = v[0]
		} else {
			m[k] = ""
		}
	}
	return m, nil
}

func main() {

	setCommandLineOptions()

	var configFile string
	args := os.Args[1:]
	if len(args) > 0 && len(args)%2 == 1 {
		// take the last one as a file argument when the number of args is odd
		configFile = args[len(args)-1]
	} else {
		configFile = "gopher.yml"
	}

	loadConfiguration(configFile)

		gopher, err := register()
	if err != nil {
		log.Fatalf("Couldn't register the service due to %v", err.Error())
	}

	log.Printf("Gopher ID: %v\n", *gopher.Id)

	listen(gopher)

}
