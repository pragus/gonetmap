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
	RingOffset unsafe.Pointer // NmRing is here
}

type Interface struct {
	File *os.File
	Nif  *NetmapIf
}

func (i *Interface) ring(idx uint32) uintptr {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(i.Nif)) + unsafe.Offsetof(i.Nif.RingOffset))
	h := *(*[]uintptr)(PtrSliceFrom(ptr, int(i.Nif.TxRings+i.Nif.RxRings+2)))
	return uintptr(unsafe.Pointer(i.Nif)) + h[idx]

}

func (i *Interface) GetRing(ringIndex interface{}, direction Direction) *NetmapRing {
	var ringPtr uintptr
	idx := ifaceTOuint32(ringIndex)
	if direction == TX {
		ringPtr = i.ring(idx)
	} else {
		ringPtr = i.ring(idx + i.Nif.TxRings + 1)
	}
	return (*NetmapRing)(unsafe.Pointer(ringPtr))
}

func (i *Interface) RxSync() error {
	return NmIoctl(i.File, NRxSync, unsafe.Pointer(uintptr(0)))
}

func (i *Interface) TxSync() error {
	return NmIoctl(i.File, NTxSync, unsafe.Pointer(uintptr(0)))
}
