package netpacket

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

func PtrSliceFrom(p unsafe.Pointer, s int) unsafe.Pointer {
	return unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(p), Len: s, Cap: s})
}

type IPv4Hdr struct {
	VersionIhl     uint8             // version and header length
	TypeOfService  uint8             // type of service
	TotalLength    uint16            // length of packet
	PacketID       uint16            // packet ID
	FragmentOffset uint16            // fragmentation offset
	TimeToLive     uint8             // time to live
	NextProtoID    uint8             // protocol ID
	HdrChecksum    uint16            // header checksum
	SrcAddr        [IPv4AddrLen]byte // source address
	DstAddr        [IPv4AddrLen]byte // destination address
}

func (h *IPv4Hdr) UpdateChecksum() {
	var chk uint32
	b := *(*[]byte)(PtrSliceFrom(unsafe.Pointer(h), IPv4MinLen))
	for i := 0; i < IPv4MinLen; i += 2 {
		// Iterating two bytes at a time; checksum bytes occur at offsets
		// 10 and 11.  Skip them.
		if i == 10 {
			continue
		}

		chk += uint32(binary.BigEndian.Uint16(b[i : i+2]))
	}

	// "The first 4 bits are the carry and will be added to the rest of
	// the value."
	carry := uint16(chk >> 16)
	csum := ^(carry + uint16(chk&0x0ffff))
	h.HdrChecksum = Htons(csum)
}
