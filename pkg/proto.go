package pkg

import "time"

// Message being sent between peers
type Message struct {
	// Sender digital twin ID
	Sender uint64 `json:"sender"`
	// Receiver digital twin ID
	Receiver uint64 `json:"receiver"`
	// Topic of the message
	Topic string `json:"topic"`
	// TTL of the message, after this time, the message should be flushed from
	// receiving and sending buffers.
	TTL time.Time `json:"ttl"`
	// Payload of the message
	Payload []byte `json:"payload"`
}
