package pkg

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/pkg/errors"
)

var errNotAuthenticated = errors.New("command requires authentication")

type unauthenticatedConn struct {
	s *Server
}

func newUnauthenticatedConn(s *Server) *unauthenticatedConn {
	return &unauthenticatedConn{
		s: s,
	}
}


func (conn *unauthenticatedConn) Auth(dtid uint64, rawSig []byte) error {
	var sig [64]byte
	switch len(rawSig) {
		case 64:
			copy(sig[:], rawSig)
		case 128:
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

func (conn *unauthenticatedConn) LPush() error {
	return errNotAuthenticated
}

func (conn *unauthenticatedConn) LPop() error {
	return errNotAuthenticated
}

func (conn *unauthenticatedConn) LLen() error {
	return errNotAuthenticated
}

func (conn *unauthenticatedConn) LRange() error {
	return errNotAuthenticated
}

func signatureValid(pk [32]byte, sig [64]byte) bool {
	return ed25519.Verify(ed25519.PublicKey(pk[:]), []byte("A"), sig[:])
}
