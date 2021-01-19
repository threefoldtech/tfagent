package pkg

// connection from a digital twin
type connection interface {
	Auth(dtid uint64, rawSig []byte) error
	LPush(receiverDtid uint64, subject string, payload []byte) error
	LPop(dtid uint64, subject string) (Message, error)
	LLen() error
	LRange() error
}
