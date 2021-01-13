package pkg

// connection from a digital twin
type connection interface {
	Auth(dtid uint64, rawSig []byte) error
	LPush() error
	LPop() error
	LLen() error
	LRange() error
}
