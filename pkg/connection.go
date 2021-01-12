package pkg

// Connection from a digital twin
type Connection interface {
	Auth() error
	LPush() error
	LPop() error
	LLen() error
	Lrange() error
}

