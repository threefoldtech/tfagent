package pkg

type authenticatedConn struct {
	dtid uint64

	s *Server
}

func newAuthenticatedConn(dtid uint64, s *Server) *authenticatedConn {
	return &authenticatedConn{
		dtid: dtid,
		s: s,
	}
}

func (conn *authenticatedConn) HELLO() error {
	return nil
}
func (conn *authenticatedConn) AUTH() error {
	return nil
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


