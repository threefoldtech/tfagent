package pkg

import "time"

type (
	// Message sent between 2 digital twins over the network
	Message struct {
		// Sender is the digital twin ID that sent the message
		Sender uint64
		// Receiver is the receiving digital twin ID
		RemoteID uint64
		// Expiration of the message, the network will attempt to deliver the
		// message untill this time has passed
		TTL time.Time
		// Topic of the message
		Topic string
		// Payload is the actual data
		Payload []byte
	}
)
