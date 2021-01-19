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

func (conn *authenticatedConn) Auth(_ uint64, _ []byte) error {
	return errAlreadyAuthenticated
}

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
		conn.s.node.recvQ = append(conn.s.node.recvQ[:idx], conn.s.node.recvQ[idx + 1:]...)
	}

	return msg, err
}

func (conn *authenticatedConn) LLen() error {
	return nil
}

func (conn *authenticatedConn) LRange() error {
	return nil
}
