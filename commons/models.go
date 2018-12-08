package commons

import (
	"encoding/hex"
)

const (
	ModeSync  = "sync"
	ModeAsync = "async"

	QueueTypeFifo = "fifo"
	QueueTypeStandard = "standard"
)

func decodeKey(encodedKey string) *[32]byte {
	var dest = new([32]byte)
	_, err := hex.Decode(dest[:], []byte(encodedKey))
	if err != nil {
		panic(err)
	}
	return dest
}

type SealedRequest struct {
	EncodedPublicKey *string `json:"public_key"` // HEX encoded of server's public key
	Cipher           *string `json:"cipher"`     // Base64 encoded cipher
}

func (r *SealedRequest) GetPublicKey() *[32]byte {
	return decodeKey(*r.EncodedPublicKey)
}

type WebhookRegister struct {
	EncodedPublicKey *string `json:"public_key"`
	Mode             string  `json:"mode"`
}

func (r *WebhookRegister) GetPublicKey() *[32]byte {
	return decodeKey(*r.EncodedPublicKey)
}

type WebhookRequestContext struct {
	ResponseQueueName *string `json:"res_queue_name"`
	Error             *string `json:"error"` // non-empty of this attribute indicates the error
	MessageId         *string `json:"message_id"`
	ReceiptHandle     *string `json:"receipt_handle"` // receipt handle of the request queue's message
}

type WebhookRequest struct {
	Context     *WebhookRequestContext `json:"context"`
	Path        *string                `json:"path"`
	QueryParams map[string]string      `json:"query_params"`
	Method      *string                `json:"method"`
	Headers     map[string]string      `json:"headers"` // HTTP headers coming from webhook
	Body        *string                `json:"body"`
}

func ErrorRequest(error string) *WebhookRequest {
	return &WebhookRequest{Context: &WebhookRequestContext{Error: &error}}
}

type WebhookResponseContext struct {
	ResponseQueueName    *string `json:"res_queue_name"`
	RequestMessageId     *string `json:"req_message_id"`
	RequestReceiptHandle *string `json:"req_receipt_handle"` // receipt handle of the request queue's message
}

type WebhookResponse struct {
	Context    *WebhookResponseContext `json:"context"`
	StatusCode int                     `json:"status_code"`
	Headers    map[string]string       `json:"headers"` // HTTP headers emitting to webhook
	Body       *string                 `json:"body"`
}

func ErrorResponse(err error, context *WebhookResponseContext) *WebhookResponse {
	msg := err.Error()
	return &WebhookResponse{
		StatusCode: 500,
		Body:       &msg,
		Context:    context,
	}
}

// We should process responses in batch fashion
type WebhookResponses struct {
	ResponseQueueName *string
	Messages          []*WebhookResponse `json:"messages"`
}

type Gopher struct {
	Id               *string `json:"id"`
	EncodedPublicKey *string `json:"encoded_public_key"` // Hex encoded Client's public key
	Mode             string  `json:"mode"`               // sync or async
	RequestQueueName *string `json:"req_queue_name"`     // a request queue name generated from UUID
}

func (g *Gopher) GetPublicKey() *[32]byte {
	return decodeKey(*g.EncodedPublicKey)
}
