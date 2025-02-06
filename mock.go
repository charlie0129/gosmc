package gosmc

import (
	"sync"
)

type MockConn struct {
	data map[string][]byte
	mu   *sync.Mutex
}

func NewMockConn() Conn {
	return &MockConn{
		data: make(map[string][]byte),
		mu:   &sync.Mutex{},
	}
}

func (c *MockConn) Open() error {
	return nil
}

func (c *MockConn) Close() error {
	return nil
}

func (c *MockConn) Read(key string) (SMCVal, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.data[key]
	if !ok {
		return SMCVal{}, ErrNoData
	}

	return SMCVal{
		Key:      key,
		DataType: "hex_",
		Bytes:    v,
	}, nil
}

func (c *MockConn) Write(key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value

	return nil
}
