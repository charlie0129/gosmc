# gosmc

A Go library to read and write Apple SMC on Macs.

This uses C-binding to Go, so you will need to use `cgo` for it to work.

## Usage

run `go get github.com/charlie0129/gosmc`

```go
package main

import (
	"fmt"

	"github.com/charlie0129/gosmc"
)

func main() {
	c := gosmc.New()

	// Open connection to SMC.
	_ = c.Open()

	// Close conncetion once we are done.
	defer c.Close()

	// Write 0x2 to CH0B or CH0C (to disable battery charging)
	_ = c.Write("CH0B", []byte{0x2})

	// Read value from CH0B
	v, _ := c.Read("CH0B")
}
```

## CLI

There is a command line program included to RW SMC. Run `make` to build it. The binary will be `bin/gosmc`.

Example:

```shell
# Read from CH0B
bin/gosmc -k CH0B
# Write 0x02 to CH0B
bin/gosmc -k CH0B -v 02
```

## Future work

Currently, it only does read and write using byte array. It provides no data conversion to common types like `int`, `float`, `string`. It can be beneficial to provide such conversions.
