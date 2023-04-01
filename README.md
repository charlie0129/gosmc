# gosmc

A Go library to read and write Apple SMC on Macs.

This uses C-binding to Go, so you will need to use `cgo` for it to work.

## Usage

run `go get github.com/charlie0129/gosmc`

See `example/main.go` for a complete example.

```go
package main

import (
	"fmt"

	"github.com/charlie0129/gosmc"
)

func main() {
	c := gosmc.NewConnection()

	// Open connection to SMC.
	_ = c.Open()

	// Close conncetion once we are done.
	defer c.Close()

	// Write 0x2 to CH0B or CH0C (to disable battery charging)
	_ = c.Write("CH0B", "02")

	// Read value from CH0B
	v, _ := c.Read("CH0B")
}
```

## Future work

Currently, it only does read and write using byte array. It provides no data conversion to common types like `int`, `float`, `string`. It can be beneficial to provide such conversions.

A CLI using this library is also a good way to go.
