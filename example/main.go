package main

import (
	"fmt"

	"github.com/charlie0129/gosmc"
)

func main() {
	c := gosmc.New()

	// Open connection to SMC.
	err := c.Open()
	if err != nil {
		panic(err)
	}

	// Close connection once we are done.
	defer c.Close()

	// Write 0x2 to CH0B/CH0C (to disable battery charging)
	err = c.Write("CH0B", "02")
	if err != nil {
		panic(err)
	}
	err = c.Write("CH0C", "02")
	if err != nil {
		panic(err)
	}

	// Read from CH0B/CH0C to check if it is written.
	v, err := c.Read("CH0B")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Read from CH0B: %#v\n", v)

	v, err = c.Read("CH0C")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Read from CH0C: %#v\n", v)
}
