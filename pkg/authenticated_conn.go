package pkg

import (
	"time"

	"github.com/pkg/errors"
)

type authenticatedConn struct {
	dtid uint64

	s *Server
}

const defaultMsgTTL = time.Hour

var errAlreadyAuthenticated = errors.New("already authenticated")
var errNoMessage = errors.New("no message for the given filter")

func newAuthenticatedConn(dtid uint64, s *Server) *authenticatedConn {
	return &authenticatedConn{
		dtid: dtid,
		s:    s,
	}
}

// Auth implements connection
func (conn *authenticatedConn) Auth(_ uint64, _ []byte) error {
	return errAlreadyAuthenticated
}

// LPush implements connection
func (conn *authenticatedConn) LPush(dtid uint64, subject string, payload []byte) error {
	msg := Message{
		Sender:   conn.dtid,
		Receiver: dtid,
		Topic:    subject,
		TTL:      time.Now().Add(defaultMsgTTL),
		Payload:  payload,
	}

	return errors.Wrap(conn.s.node.Send(msg), "could not send message")
}

// LPop implements connection
func (conn *authenticatedConn) LPop(dtid uint64, subject string) (Message, error) {
	conn.s.node.recvQLock.Lock()
	defer conn.s.node.recvQLock.Lock()

	var idx int
	var msg Message
	err := errNoMessage
	for i, m := range conn.s.node.recvQ {
		if m.Receiver == conn.dtid && (dtid == 0 || m.Sender == dtid) && (subject == "" || m.Topic == subject) {
			idx = i
			m = msg
			err = nil
			break
		}
	}

	if err == nil {
		conn.s.node.recvQ = append(conn.s.node.recvQ[:idx], conn.s.node.recvQ[idx+1:]...)
	}

	return msg, err
}

// LLen implements connection
func (conn *authenticatedConn) LLen(dtid uint64, subject string) (uint64, error) {
	conn.s.node.recvQLock.Lock()
	defer conn.s.node.recvQLock.Lock()

	var count uint64
	for _, m := range conn.s.node.recvQ {
		if m.Receiver == conn.dtid && (dtid == 0 || m.Sender == dtid) && (subject == "" || m.Topic == subject) {
			count++
		}
	}

	return count, nil
}

// LRange implements connection
func (conn *authenticatedConn) LRange(dtid uint64, subject string, start int, end int) ([]Message, error) {
	conn.s.node.recvQLock.Lock()
	defer conn.s.node.recvQLock.Lock()

	foundIDs := make([]int, 0, (end-start)+1)
	messages := make([]Message, 0, (end-start)+1)
	var idx int

	for i, m := range conn.s.node.recvQ {
		if m.Receiver == conn.dtid && (dtid == 0 || m.Sender == dtid) && (subject == "" || m.Topic == subject) {
			if idx >= start && idx <= end {
				messages = append(messages, m)
				foundIDs = append(foundIDs, i)
			}
			idx++
		}
	}

	// TODO: should LRANGE delete elements?
	// iterate over the found ID's backwards - otherwise the id's will be invalid
	// after the first removal
	// for i := len(foundIDs) - 1; i > 0; i-- {
	// 	conn.s.node.recvQ = append(conn.s.node.recvQ[:foundIDs[i]], conn.s.node.recvQ[foundIDs[i]+1:]...)
	// }

	return messages, nil
}
