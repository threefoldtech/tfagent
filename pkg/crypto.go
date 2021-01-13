package pkg

import "crypto/ed25519"

const (
	PublicKeySize = ed25519.PublicKeySize
	SignatureSize = ed25519.SignatureSize
)

func signatureValid(pk [PublicKeySize]byte, sig [SignatureSize]byte) bool {
	return ed25519.Verify(ed25519.PublicKey(pk[:]), []byte("A"), sig[:])
}
