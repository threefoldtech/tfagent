package pkg

import "github.com/pkg/errors"

type authenticatedConn struct {
	dtid uint64

	s *Server
}

var errAlreadyAuthenticated = errors.New("already authenticated")

func newAuthenticatedConn(dtid uint64, s *Server) *authenticatedConn {
	return &authenticatedConn{
		dtid: dtid,
		s: s,
	}
}


func (conn *authenticatedConn) AUTH() error {
	return errAlreadyAuthenticated
}

func (conn *authenticatedConn) LPUSH() error {
	return nil
}

func (conn *authenticatedConn) LPOP() error {
	return nil
}

func (conn *authenticatedConn) LLEN() error {
	return nil
}

func (conn *authenticatedConn) LRANGE() error {
	return nil
}

