package pkg

// connection from a digital twin
type connection interface {
	Auth(dtid uint64, rawSig []byte) error
	LPush(receiverDtid uint64, subject string, payload []byte) error
	LPop(dtid uint64, subject string) (Message, error)
	LLen(dtid uint64, subject string) (uint64, error)
	LRange(dtid uint64, subject string, start int, end int) ([]Message, error)
}
