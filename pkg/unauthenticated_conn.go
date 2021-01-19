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

func (conn *unauthenticatedConn) LPush(_ uint64, _ string, _ []byte) error {
	return errNotAuthenticated
}

func (conn *unauthenticatedConn) LPop(_ uint64, _ string) (Message, error) {
	return Message{}, errNotAuthenticated
}

func (conn *unauthenticatedConn) LLen() error {
	return errNotAuthenticated
}

func (conn *unauthenticatedConn) LRange() error {
	return errNotAuthenticated
}
