package gosmc

type Connection interface {
	Open() error
	Close() error
	Read(key string) (SMCVal, error)
	Write(key string, value []byte) error
}
