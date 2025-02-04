package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/charlie0129/gosmc"
)

func main() {
	key := ""
	val := ""
	flag.StringVar(&key, "k", "", "SMC key to read/write. Value will be printed as hex.")
	flag.StringVar(&val, "v", "", "Value (as hex) to write to SMC key. If not specified, read value from SMC key")
	flag.Parse()

	if key == "" {
		log.Fatal("You must specify a key to read/write with -k")
	}

	c := gosmc.NewAppleSMCConn()

	// Open connection to SMC.
	err := c.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Close connection once we are done.
	defer c.Close()

	// Read value from SMC key.
	if val == "" {
		v, err := c.Read(key)
		if err != nil {
			log.Fatal(err)
		}
		// Print bytes as hex.
		for _, b := range v.Bytes {
			fmt.Printf("%02x ", b)
		}
		fmt.Println()
		return
	}

	// Write value to SMC key.
	if len(val)%2 != 0 {
		log.Fatal("Value must be hex encoded\n")
	}

	b := make([]byte, len(val)/2)
	for i := 0; i < len(val); i += 2 {
		v, err := strconv.ParseUint(val[i:i+2], 16, 8)
		if err != nil {
			log.Fatal(err)
		}
		b[i/2] = byte(v)
	}

	err = c.Write(key, b)
	if err != nil {
		log.Fatal(err)
	}
}
