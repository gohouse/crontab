package adapter

type Logger interface {
	Info(arg string) error
	Infof(arg string, args ...interface{}) error
	Error(err error) error
	Read(readRows ...int64) string
}
