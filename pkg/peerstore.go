package pkg

// TwinInfo information about a digital twin
type TwinInfo struct {
	// PubKey of this digital twin
	PubKey [32] byte
	// Peer id of the libp2p daemon this twin is connected to
	Peer string
}

// PeerStore allows looking up info for a peer digital twin ID
type PeerStore interface {
	// PeerId currently associated with this digital twin
	PeerId(dtid uint64) (string, error)
	// PublicKey of this digital twin
	PublicKey(dtid uint64) ([32]byte, error)
	// SetPeerId of this digital twin. This should override a cached peer ID of
	// a digital twin. NOTE: because the transport does not validate keys of remote
	// digital twins, a malicous entity could emulate a peer ID and poison the peer
	// cache. As such, this should only be used in development.
	SetPeerId(dtid uint64, pid string)
}
