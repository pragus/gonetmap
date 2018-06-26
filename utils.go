package gonetmap

import (
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

func ifaceTOuint16(i interface{}) uint16 {
	var idx uint16
	switch i.(type) {
	case int:
		idx = uint16(i.(int))
	case int16:
		idx = uint16(i.(int16))
	case int32:
		idx = uint16(i.(int32))
	case int64:
		idx = uint16(i.(int64))

	case uint:
		idx = uint16(i.(uint))
	case uint16:
		idx = uint16(i.(uint16))
	case uint32:
		idx = uint16(i.(uint32))
	case uint64:
		idx = uint16(i.(uint64))
	}
	return idx
}

func ifaceTOuint32(i interface{}) uint32 {
	var idx uint32
	switch i.(type) {
	case int:
		idx = uint32(i.(int))
	case int16:
		idx = uint32(i.(int16))
	case int32:
		idx = uint32(i.(int32))
	case int64:
		idx = uint32(i.(int64))

	case uint:
		idx = uint32(i.(uint))
	case uint16:
		idx = uint32(i.(uint16))
	case uint32:
		idx = uint32(i.(uint32))
	case uint64:
		idx = uint32(i.(uint64))
	}
	return idx
}

func nmMmap(file *os.File, r *Request) (ptr uintptr, err error) {
	fd := int(file.Fd())
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	if data, err := syscall.Mmap(fd, 0, int(r.Memsize), prot, syscall.MAP_SHARED); err == nil {
		ptr = (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data
	}
	return
}
