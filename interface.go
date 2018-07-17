package gonetmap

import (
	"os"
	"unsafe"
)

type NetmapIf struct {
	Name       [16]byte
	Version    uint32
	Flags      uint32
	TxRings    uint32
	RxRings    uint32
	BufsHead   uint32
	Spare1     [5]uint32
	RingOffset unsafe.Pointer //NmRing is here
}

type Interface struct {
	File *os.File
	Nif  *NetmapIf
}

func (i *Interface) ring(idx uint32) uintptr {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(i.Nif)) + unsafe.Offsetof(i.Nif.RingOffset))
	h := *(*[]uintptr)(ptrSliceFrom(ptr, int(i.Nif.TxRings+i.Nif.RxRings+2)))
	return uintptr(unsafe.Pointer(i.Nif)) + h[idx]

}

func (i *Interface) GetRing(RingIndex interface{}, direction Direction) *NetmapRing {
	var ring_ptr uintptr
	idx := ifaceTOuint32(RingIndex)
	if direction == TX {
		ring_ptr = i.ring(idx)
	} else {
		ring_ptr = i.ring(idx + i.Nif.TxRings + 1)
	}
	return (*NetmapRing)(unsafe.Pointer(ring_ptr))
}

func (i *Interface) RxSync() error {
	return NmIoctl(i.File, NRxSync, unsafe.Pointer(uintptr(0)))
}

func (i *Interface) TxSync() error {
	return NmIoctl(i.File, NTxSync, unsafe.Pointer(uintptr(0)))
}
