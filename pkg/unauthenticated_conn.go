package pkg

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

var (
	errNotAuthenticated       = errors.New("command requires authentication")
	errInvalidSignatureLength = errors.New("invalid signature length")
)

type unauthenticatedConn struct {
	s *Server
}

func newUnauthenticatedConn(s *Server) *unauthenticatedConn {
	return &unauthenticatedConn{
		s: s,
	}
}

// Auth implements connection
func (conn *unauthenticatedConn) Auth(dtid uint64, rawSig []byte) error {
	var sig [SignatureSize]byte
	switch len(rawSig) {
	case SignatureSize:
		copy(sig[:], rawSig)
	case SignatureSize * 2:
		data, err := hex.DecodeString(string(rawSig))
		if err != nil {
			return err
		}
		copy(sig[:], data)
	default:
		return errInvalidSignatureLength
	}

	pk, err := conn.s.ps.PublicKey(dtid)
	if err != nil {
		return errors.Wrap(err, "could not get public key")
	}

	if !signatureValid(pk, sig) {
		return errAuthorizationFailed
	}

	return nil
}

// LPush implements connection
func (conn *unauthenticatedConn) LPush(_ uint64, _ string, _ []byte) error {
	return errNotAuthenticated
}

// LPop implements connection
func (conn *unauthenticatedConn) LPop(_ uint64, _ string) (Message, error) {
	return Message{}, errNotAuthenticated
}

// LLen implements connection
func (conn *unauthenticatedConn) LLen(_ uint64, _ string) error {
	return errNotAuthenticated
}

// LRange implements connection
func (conn *unauthenticatedConn) LRange(_ uint64, _ string, _ int, _ int) ([]Message, error) {
	return nil, errNotAuthenticated
}
