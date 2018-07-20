package netpacket

import (
	"fmt"
	"unsafe"
)

func (e *EtherHdr) String() string {
	return fmt.Sprintf("%#04x, %02x:%02x:%02x:%02x:%02x:%02x -> %02x:%02x:%02x:%02x:%02x:%02x",
		e.EtherType,
		e.SAddr[0], e.SAddr[1], e.SAddr[2], e.SAddr[3], e.SAddr[4], e.SAddr[5],
		e.DAddr[0], e.DAddr[1], e.DAddr[2], e.DAddr[3], e.DAddr[4], e.DAddr[5],
	)
}

func (e *EtherHdr) getNext() unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(e)) + EtherLen)
}

func (e *EtherHdr) GetIP() *IPv4Hdr {
	return (*IPv4Hdr)(e.getNext())
}

func (e *EtherHdr) GetARP() *ARPHdr {
	return (*ARPHdr)(e.getNext())
}
