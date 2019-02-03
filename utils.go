package gonetmap

import (
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

func PtrSliceFrom(p unsafe.Pointer, s int) unsafe.Pointer {
	return unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(p), Len: s, Cap: s})
}

func ifaceTOuint16(i interface{}) uint16 {
	var idx uint16
	switch i := i.(type) {
	case int:
		idx = uint16(i)
	case int16:
		idx = uint16(i)
	case int32:
		idx = uint16(i)
	case int64:
		idx = uint16(i)

	case uint:
		idx = uint16(i)
	case uint16:
		idx = uint16(i)
	case uint32:
		idx = uint16(i)
	case uint64:
		idx = uint16(i)
	}
	return idx
}

func ifaceTOuint32(i interface{}) uint32 {
	var idx uint32
	switch i := i.(type) {
	case int:
		idx = uint32(i)
	case int16:
		idx = uint32(i)
	case int32:
		idx = uint32(i)
	case int64:
		idx = uint32(i)

	case uint:
		idx = uint32(i)
	case uint16:
		idx = uint32(i)
	case uint32:
		idx = uint32(i)
	case uint64:
		idx = uint32(i)
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
