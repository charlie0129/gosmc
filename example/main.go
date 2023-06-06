package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

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

	key := ""
	val := ""
	flag.StringVar(&key, "k", "", "SMC key to read/write")
	flag.StringVar(&val, "v", "", "Value to write to SMC key")
	flag.Parse()

	if key == "" {
		fmt.Printf("You must specify a key to read/write with -k\n")
		os.Exit(1)
	}

	// Read value from SMC key.
	if val == "" {
		v, err := c.Read(key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %v\n", v.Key, v.Bytes)
		return
	}

	// Write value to SMC key.
	if len(val)%2 != 0 {
		fmt.Printf("Value must be hex encoded\n")
		os.Exit(1)
	}

	b := make([]byte, len(val)/2)
	for i := 0; i < len(val); i += 2 {
		v, err := strconv.ParseUint(val[i:i+2], 16, 8)
		if err != nil {
			panic(err)
		}
		b[i/2] = byte(v)
	}

	err = c.Write(key, b)
	if err != nil {
		panic(err)
	}
}
