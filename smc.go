package gosmc

/*
#cgo LDFLAGS: -framework IOKit
#include <stdlib.h>
#include "smc.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type SMCVal struct {
	Key      string
	DataType string
	Bytes    []byte
}

var (
	ErrKeyLength = errors.New("key must be 4 characters long")
	ErrNoData    = errors.New("key has no data, check if it is valid")
)

type AppleSMCConn uint

func New() Conn {
	c := AppleSMCConn(0)
	return &c
}

func (c *AppleSMCConn) Open() error {
	ret := int(C.SMCOpen((*C.uint)(unsafe.Pointer(c))))
	// TODO: pass errors strings from C
	if ret != 0 {
		return fmt.Errorf("error when opening SMC, ret=%d", ret)
	}
	return nil
}

func (c *AppleSMCConn) Close() error {
	ret := int(C.SMCClose(C.uint(*c)))
	if ret != 0 {
		return fmt.Errorf("error when closing SMC, ret=%d", ret)
	}
	return nil
}

func (c *AppleSMCConn) Write(key string, val []byte) error {
	if len(key) != 4 {
		return ErrKeyLength
	}

	var ckey *C.char = C.CString(key)
	var cval unsafe.Pointer = C.CBytes(val)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(cval)

	ret := int(C.SMCWriteSimple(ckey, (*C.uchar)(cval), C.int(len(val)), C.uint(*c)))
	if ret != 0 {
		return fmt.Errorf("error when writing %s to %s, ret=%d", val, key, ret)
	}

	return nil
}

func (c *AppleSMCConn) Read(key string) (SMCVal, error) {
	if len(key) != 4 {
		return SMCVal{}, ErrKeyLength
	}

	var ckey *C.char = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	v := C.SMCVal_t{}

	ret := int(C.SMCReadKey2(ckey, (*C.SMCVal_t)(unsafe.Pointer(&v)), C.uint(*c)))
	if ret != 0 {
		return SMCVal{}, fmt.Errorf("error when reading %s, ret=%d", key, ret)
	}

	if uint32(v.dataSize) == 0 {
		return SMCVal{}, ErrNoData
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
