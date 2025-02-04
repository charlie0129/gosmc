package gosmc

/*
#cgo LDFLAGS: -framework IOKit

#include <stdlib.h>
#include <IOKit/IOKitLib.h>

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

type AppleSMCConn struct {
	conn      C.io_connect_t
	infoCache map[uint32]C.SMCKeyData_keyInfo_t
}

var _ Connection = &AppleSMCConn{}

func NewAppleSMCConn() *AppleSMCConn {
	c := AppleSMCConn{
		conn: C.io_connect_t(0),
	}
	return &c
}

func (c *AppleSMCConn) Open() error {
	var result C.kern_return_t
	var mainPort C.mach_port_t
	var iterator C.io_iterator_t
	var device C.io_object_t

	result = C.IOMainPort(0, &mainPort)
	if result != C.kIOReturnSuccess {
		return fmt.Errorf("error when getting main port, ret=%x", int32(result))
	}

	appleSMCString := C.CString("AppleSMC")
	defer C.free(unsafe.Pointer(appleSMCString))

	var matchingDict C.CFDictionaryRef = (C.CFDictionaryRef)(C.IOServiceMatching(appleSMCString))
	result = C.IOServiceGetMatchingServices(mainPort, matchingDict, &iterator)
	if result != C.kIOReturnSuccess {
		return fmt.Errorf("error when getting matching services, ret=%x", int32(result))
	}

	device = C.IOIteratorNext(iterator)
	C.IOObjectRelease(iterator)
	if device == C.io_object_t(0) {
		return fmt.Errorf("no AppleSMC found")
	}

	result = C.IOServiceOpen(device, C.mach_task_self_, C.uint32_t(0), &c.conn)
	C.IOObjectRelease(device)
	if result != C.kIOReturnSuccess {
		return fmt.Errorf("error when opening AppleSMC, ret=%x", int32(result))
	}

	return nil
}

func (c *AppleSMCConn) Close() error {
	ret := C.IOServiceClose(c.conn)
	if ret != C.kIOReturnSuccess {
		return fmt.Errorf("error when closing AppleSMC, ret=%x", int32(ret))
	}
	return nil
}

func (c *AppleSMCConn) call(selector C.uint32_t, input *C.SMCKeyData_t, output *C.SMCKeyData_t) error {
	ret := C.IOConnectCallStructMethod(
		C.mach_port_t(c.conn),
		C.uint32_t(selector),
		input,
		C.sizeof(C.SMCKeyData_t),
		output,
		C.sizeof(C.SMCKeyData_t),
	)
	if ret != C.kIOReturnSuccess {
		return fmt.Errorf("error when calling method %d, ret=%x", uint32(selector), int32(ret))
	}

	return nil
}

func (c *AppleSMCConn) getKeyInfo(key C.Uint32, keyInfo *C.SMCKeyData_keyInfo_t) error {

}

func (c *AppleSMCConn) Write(key string, val []byte) error {
	if len(key) != 4 {
		return ErrKeyLength
	}

	var ckey *C.char = C.CString(key)
	var cval unsafe.Pointer = C.CBytes(val)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(cval)

	ret := int(C.SMCWriteSimple(ckey, (*C.uchar)(cval), C.int(len(val)), c.conn))
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

	ret := int(C.SMCReadKey2(ckey, (*C.SMCVal_t)(unsafe.Pointer(&v)), c.conn))
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
