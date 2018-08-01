package gonetmap

import (
	"reflect"
	"syscall"
	"unsafe"
)

type Stat struct {
	Received  uint32
	Dropped   uint32
	IfDropped uint32
}

type Slot struct {
	Idx   uint32
	Len   uint16
	Flags uint16
	Ptr   uintptr
}

type PacketHeader struct {
	Ts     syscall.Timeval
	Caplen uint32
	Len    uint32
	Flags  uint64
	Desc   *Descriptor
	Slot   *Slot
	Buf    *uint8
}

type Descriptor struct {
	Self        *Descriptor
	Fd          int32
	pad0        [4]byte
	Mem         *byte
	Memsize     uint32
	DoneMmap    int32
	NetmapIf    *NetmapIf
	FirstTxRing uint16
	LastTxRing  uint16
	CurTxRing   uint16
	FirstRxRing uint16
	LastRxRing  uint16
	CurRxRing   uint16
	Request     Request
	Header      PacketHeader
	SomeRing    *NetmapRing
	BufStart    *byte
	BufEnd      *byte
	Snaplen     int32
	Promisc     int32
	ToMs        int32
	pad1        [4]byte
	ErrBuf      *int8
	IfaceFlags  uint32
	IfaceReqcap uint32
	IfaceCurcap uint32
	Stat        Stat
	Msg         [512]byte
}

