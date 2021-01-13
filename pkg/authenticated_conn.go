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
		s:    s,
	}
}

func (conn *authenticatedConn) Auth(_ uint64, _ []byte) error {
	return errAlreadyAuthenticated
}

func (conn *authenticatedConn) LPush() error {
	return nil
}

func (conn *authenticatedConn) LPop() error {
	return nil
}

func (conn *authenticatedConn) LLen() error {
	return nil
}

func (conn *authenticatedConn) LRange() error {
	return nil
}
