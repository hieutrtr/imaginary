package kafka

import "encoding/json"

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
