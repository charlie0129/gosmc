package gosmc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMC(t *testing.T) {
	a := assert.New(t)

	c := New()
	err := c.Open()
	a.Nil(err)

	defer c.Close()

	// **requires root access** disable charging
	err = c.Write("CH0B", "02")
	a.Nil(err)

	v, err := c.Read("CH0B")
	a.Nil(err)
	a.Equal("CH0B", v.Key)
	a.Equal("hex_", v.DataType)
	a.Equal([]uint8{0x3}, v.Bytes)
}
