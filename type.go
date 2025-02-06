package gosmc

type Conn interface {
	// Open opens the connection to SMC. You must call this before any other operations. After you are done,
	// you should call Close.
	Open() error
	// Close closes the connection to SMC. You should call this after you are done with the connection.
	Close() error
	// Read reads a value from SMC. The key must be 4 characters long.
	Read(key string) (SMCVal, error)
	// Write writes a value to SMC. The key must be 4 characters long.
	Write(key string, value []byte) error
}
