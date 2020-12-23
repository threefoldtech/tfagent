package pkg

type Log interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

type NOOPLogger struct{}

func (n *NOOPLogger) Debug(v ...interface{})                 {}
func (n *NOOPLogger) Debugf(v ...interface{})                {}
func (n *NOOPLogger) Info(v ...interface{})                  {}
func (n *NOOPLogger) Infof(format string, v ...interface{})  {}
func (n *NOOPLogger) Warn(v ...interface{})                  {}
func (n *NOOPLogger) Warnf(format string, v ...interface{})  {}
func (n *NOOPLogger) Error(v ...interface{})                 {}
func (n *NOOPLogger) Errorf(format string, v ...interface{}) {}
