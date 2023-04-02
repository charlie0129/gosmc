package gosmc

/*
#cgo LDFLAGS: -framework IOKit
#include <stdlib.h>
#include "smc.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type SMCVal struct {
	Key      string
	DataType string
	Bytes    []byte
}

type cSMCVal struct {
	val C.SMCVal_t
}

type Connection uint

func New() *Connection {
	c := Connection(0)
	return &c
}

func (c *Connection) Open() error {
	// var conn C.uint
	ret := int(C.SMCOpen((*C.uint)(unsafe.Pointer(c))))
	// TODO: pass errors strings from C
	if ret != 0 {
		return fmt.Errorf("error when opening SMC, ret=%d", ret)
	}
	return nil
}

func (c *Connection) Close() error {
	ret := int(C.SMCClose(C.uint(*c)))
	if ret != 0 {
		return fmt.Errorf("error when closing SMC, ret=%d", ret)
	}
	return nil
}

func (c *Connection) Write(key string, val string) error {
	if len(key) > 4 {
		panic(fmt.Sprintf("key %s too long", key))
	}

	var ckey *C.char = C.CString(key)
	var cval *C.char = C.CString(val)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cval))

	ret := int(C.SMCWriteSimple(ckey, cval, C.uint(*c)))
	if ret != 0 {
		return fmt.Errorf("error when writing %s to %s, ret=%d", val, key, ret)
	}

	return nil
}

func (c *Connection) Read(key string) (SMCVal, error) {
	if len(key) > 4 {
		panic(fmt.Sprintf("key %s too long", key))
	}

	var ckey *C.char = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	v := C.SMCVal_t{}

	ret := int(C.SMCReadKey2(ckey, (*C.SMCVal_t)(unsafe.Pointer(&v)), C.uint(*c)))
	if ret != 0 {
		return SMCVal{}, fmt.Errorf("error when reading %s, ret=%d", key, ret)
	}

	bytes := C.GoBytes(unsafe.Pointer(&v.bytes), 32)
	bytes = bytes[:uint32(v.dataSize)]

	val := SMCVal{
		Key:      C.GoString((*C.char)(unsafe.Pointer(&v.key))),
		DataType: C.GoString((*C.char)(unsafe.Pointer(&v.dataType))),
		Bytes:    bytes,
	}

	return val, nil
}
