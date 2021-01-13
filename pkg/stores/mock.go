package stores

import (
	"encoding/hex"

	"github.com/threefoldtech/tfagent/pkg"
)

// pubkey hex: 74856cfef93872537edaebd19504e6494beabc33f61abac91da7301f0f37f655
// sig hex: 60e2359b98c9f021df38c45c0c498da5faab789c927589a37d195471b14bfb5a81d13018d5fbd1f656a0a6957cbd62019f32450c2631e01f81f9d082bfc61606

// MockStore always returns a default key
type MockStore struct {}

// PeerID implements pkg.PeerStore 
func (m MockStore) PeerID(dtid uint64) (string, error) {
	return "", nil	
}
// PublicKey implements pkg.PeerStore
func (m MockStore) PublicKey(dtid uint64) ([pkg.PublicKeySize]byte, error) {
	key := [32]byte{}
	sb, err := hex.DecodeString("74856cfef93872537edaebd19504e6494beabc33f61abac91da7301f0f37f655")
	copy(key[:], sb)
	return key, err 
}
// SetPeerID implements pkg.PeerStore
func (m MockStore) SetPeerID(dtid uint64, pid string) {
	
}
