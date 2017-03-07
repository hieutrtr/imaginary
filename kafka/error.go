package kafka

import (
	"encoding/json"
	"strings"
)

const (
	Unavailable uint8 = iota
	BadRequest
	NotAllowed
	Unsupported
	Unauthorized
	InternalError
	NotFound
)

var (
	ErrMessage  = NewError("Invailid or missing message", BadRequest)
	ErrProducer = NewError("Producer is not ready", InternalError)
)

// NewError create new error
func NewError(err string, code uint8) Error {
	err = strings.Replace(err, "\n", "", -1)
	return EventError{err, code}
}

// EventError Kafka error handling
type EventError struct {
	Message string `json:"message"`
	Code    uint8  `json:"code"`
}

func (e EventError) JSON() []byte {
	buf, _ := json.Marshal(e)
	return buf
}

func (e EventError) Error() string {
	return e.Message
}
