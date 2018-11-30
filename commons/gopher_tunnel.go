package commons

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tailrecio/gopher-tunnels/config"
	"log"
	"sync"
)

var keyPairMu sync.Mutex
var keyPair *KeyPair

var qSessionMu sync.Mutex
var qSession *sqs.SQS

var responseQueueMu sync.Mutex
var responseQueueNameMap = make(map[string]*string)

func GetKeyPair() *KeyPair {
	if keyPair == nil {
		keyPairMu.Lock()
		defer keyPairMu.Unlock()
		if keyPair == nil {
			keyPair = NewKeyPair()
		}
	}
	return keyPair
}

func getQueueSession() *sqs.SQS {
	if qSession == nil {
		qSessionMu.Lock()
		defer qSessionMu.Unlock()
		if qSession == nil {
			qSession = sqs.New(NewAwsSession())
		}
	}
	return qSession
}

func GetQueueUrl(queueName *string) *string {
	if queueName == nil {
		panic("queueName must not be empty!")
	}
	var queueUrl string
	if config.GetBaseQueueEndpoint() != "" {
		queueUrl = fmt.Sprintf("%v/%v", config.GetBaseQueueEndpoint(), *queueName)
	} else {
		queueUrl = fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v", config.GetAwsRegion(), config.GetAwsAccount(), *queueName)
	}
	return &queueUrl
}

func SendRequest(gopher *Gopher, request *WebhookRequest) error {
	requestQueueUrl := GetQueueUrl(gopher.RequestQueueName)
	var jsonBytes []byte
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	jsonStr := string(jsonBytes)

	cipher := Encrypt(gopher.GetPublicKey(), GetKeyPair().PrivateKey, &jsonStr)
	requestEnvelope := SealedRequest{
		EncodedPublicKey: GetKeyPair().GetHexEncodedPublicKey(),
		Cipher:           cipher,
	}

	_, err = sendMessage(requestQueueUrl, &requestEnvelope)

	return err
}

func ReadRequests(gopher *Gopher) ([]*WebhookRequest, error) {
	requestQueueUrl := GetQueueUrl(gopher.RequestQueueName)
	messages, err := readMessages(requestQueueUrl, 10) // TODO: should be configurable
	if err != nil {
		return nil, err
	}

	requests := make([]*WebhookRequest, len(messages))

	for i, message := range messages {
		var sealedRequest SealedRequest
		var request *WebhookRequest
		err = json.Unmarshal([]byte(*message.Body), &sealedRequest)
		if err != nil {
			request = ErrorRequest("unmarshal queue's message body error: " + err.Error())
		} else {
			var plainText *string
			plainText, err = Decrypt(sealedRequest.GetPublicKey(), GetKeyPair().PrivateKey, sealedRequest.Cipher)
			if err != nil {
				request = ErrorRequest(err.Error())
			} else {
				var r WebhookRequest
				err = json.Unmarshal([]byte(*plainText), &r)
				request = &r
				if err != nil {
					request = ErrorRequest("unmarshal decrypted message error: " + err.Error())
				}
			}
		}
		if request.Context == nil {
			request.Context = &WebhookRequestContext{}
		}
		request.Context.MessageId = message.MessageId
		request.Context.ReceiptHandle = message.ReceiptHandle
		requests[i] = request
	}

	return requests, nil
}

func SendResponse(gopher *Gopher, response *WebhookResponse) error {
	log.Printf("Sending a response back to Gopher: %v\n", *gopher.Id)
	// Delete request queue's item since the message will not be deleted automatically.
	requestQueueUrl := GetQueueUrl(gopher.RequestQueueName)
	deleteInput := sqs.DeleteMessageInput{
		QueueUrl:      requestQueueUrl,
		ReceiptHandle: response.Context.RequestReceiptHandle,
	}
	log.Printf("Deleting the request message: `%v` from the queue: %v\n", *deleteInput.ReceiptHandle, *requestQueueUrl)
	_, err := getQueueSession().DeleteMessage(&deleteInput)

	if err != nil {
		log.Printf("Failed to delete the message: `%v` from the queue: %v\n", *deleteInput.ReceiptHandle, *requestQueueUrl)
		return err
	}

	if gopher.Mode == ModeSync {
		responseQueueUrl := GetQueueUrl(response.Context.ResponseQueueName)
		_, err = sendMessage(responseQueueUrl, response)
		return err
	}

	return nil
}

func ReadResponse(responseQueueName *string) (*WebhookResponse, error) {
	var err error
	responseQueueUrl := responseQueueNameMap[*responseQueueName]
	if responseQueueUrl == nil {
		log.Printf("Response queue not found for name: %v\n", *responseQueueName)
		responseQueueMu.Lock()
		defer responseQueueMu.Unlock()
		if responseQueueNameMap[*responseQueueName] == nil {
			log.Printf("Fetching a response queue URL for name: %v\n", *responseQueueName)
			getQueueUrlInput := sqs.GetQueueUrlInput{
				QueueName: responseQueueName,
			}
			var getQueueUrlOutput *sqs.GetQueueUrlOutput
			getQueueUrlOutput, err = getQueueSession().GetQueueUrl(&getQueueUrlInput)
			if err != nil {
				// queue doesn't exist!
				log.Printf("Response queue: `%v` doesn't exist. Creating a new one...\n", *responseQueueName)
				responseQueueUrl, err = createQueue(responseQueueName, nil)
				if err != nil {
					return nil, err
				}
			} else {
				responseQueueUrl = getQueueUrlOutput.QueueUrl
			}
			responseQueueNameMap[*responseQueueName] = responseQueueUrl
		}
	}

	var messages []*sqs.Message
	messages, err = readMessages(responseQueueUrl, 1)
	if err != nil {
		return nil, err
	}
	if len(messages) > 0 {
		message := messages[0]
		var response WebhookResponse
		err = json.Unmarshal([]byte(*message.Body), &response)
		if err != nil {
			return nil, err
		}
		deleteMessageInput := sqs.DeleteMessageInput{
			QueueUrl:      responseQueueUrl,
			ReceiptHandle: message.ReceiptHandle,
		}
		_, err = getQueueSession().DeleteMessage(&deleteMessageInput)

		return &response, err
	}

	return nil, nil
}

func readMessages(queueUrl *string, maxNumberOfMessages int64) ([]*sqs.Message, error) {
	log.Printf("Reading messages from queue: %v\n", *queueUrl)
	input := sqs.ReceiveMessageInput{
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: aws.Int64(maxNumberOfMessages),
		WaitTimeSeconds:     aws.Int64(10), // TODO: should be configurable
	}
	output, err := getQueueSession().ReceiveMessage(&input)
	if err != nil {
		return nil, err
	}
	return output.Messages, err
}

func sendMessage(queueUrl *string, message interface{}) (*sqs.SendMessageOutput, error) {
	log.Printf("Sending a message to queue: %v\n", *queueUrl)
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	// TODO: hmmm is it bad to set MessageGroupId as a constant?
	input := sqs.SendMessageInput{
		QueueUrl:       queueUrl,
		MessageBody:    aws.String(string(jsonBytes)),
		MessageGroupId: aws.String("1"), // https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/using-messagegroupid-property.html
	}
	return getQueueSession().SendMessage(&input)
}

// Create a request queue and response queue (if applicable)
// The policy for a request queue is public accessible.
// The queue will be protected by two mechanism
// 1) Queue Name (UUID is virtually impossible to guess)
// 2) The queue is valid for a single tunnel only (will be destroyed when the session is invalidated)
// 3) The data is encrypted using ECDH
// 4) The lifespan of data is very short
func CreateRequestQueue(gopher *Gopher) (*string, error) {
	queueName := gopher.RequestQueueName
	resource := fmt.Sprintf("arn:aws:sqs:%v:%v:%v", config.GetAwsRegion(), config.GetAwsAccount(), *queueName)
	policy := fmt.Sprintf(`
{
  "Version": "2012-10-17",
  "Id": "%v/SQSDefaultPolicy",
  "Statement": [
    {
      "Sid": "Gopher_Tunnels_Request",
      "Effect": "Allow",
      "Principal": "*",
      "Action": [
        "SQS:GetQueueUrl",
        "SQS:ReceiveMessage"
      ],
      "Resource": "%v"
    }
  ]
}
`, resource, resource)
	return createQueue(queueName, &policy)
}

func CreateResponseQueue() (*string, error) {
	queueName := MakeQueueName("out")
	return createQueue(&queueName, nil)
}

func createQueue(queueName *string, policy *string) (*string, error) {
	attributes := map[string]*string{
		"MessageRetentionPeriod":    aws.String("60"), // 60 seconds // TODO: should be configurable
		"FifoQueue":                 aws.String("true"),
		"ContentBasedDeduplication": aws.String("true"), // TODO: not sure if we want this to be configurable
	}
	if policy != nil {
		attributes["Policy"] = policy
	}
	input := sqs.CreateQueueInput{
		QueueName:  queueName,
		Attributes: attributes,
	}
	out, err := getQueueSession().CreateQueue(&input)
	if err != nil {
		return nil, err
	}
	return out.QueueUrl, err
}
