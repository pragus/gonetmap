package netpacket

import (
	"encoding/binary"
	"unsafe"
	"github.com/pragus/gonetmap"

)

const (
	Idx0 = 0
	Idx1 = 2
	Idx2 = 4
	Idx3 = 6
	Idx4 = 8
	Idx5 = 12
	Idx6 = 14
	Idx7 = 16
	Idx8 = 18

)


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

	b := *(*[]byte)(gonetmap.PtrSliceFrom(unsafe.Pointer(h), IPv4MinLen))

	p0 := uint32(binary.BigEndian.Uint16(b[Idx0: Idx0+2]))
	p1 := uint32(binary.BigEndian.Uint16(b[Idx1: Idx1+2]))
	p2 := uint32(binary.BigEndian.Uint16(b[Idx2: Idx2+2]))
	p3 := uint32(binary.BigEndian.Uint16(b[Idx3: Idx3+2]))
	p4 := uint32(binary.BigEndian.Uint16(b[Idx4: Idx4+2]))
	p5 := uint32(binary.BigEndian.Uint16(b[Idx5: Idx5+2]))
	p6 := uint32(binary.BigEndian.Uint16(b[Idx6: Idx6+2]))
	p7 := uint32(binary.BigEndian.Uint16(b[Idx7: Idx7+2]))
	p8 := uint32(binary.BigEndian.Uint16(b[Idx8: Idx8+2]))
	chk := p0 + p1 + p2 + p3 + p4 + p5 + p6 + p7 + p8

	// "The first 4 bits are the carry and will be added to the rest of
	// the value."
	carry := uint16(chk >> 16)
	csum := ^(carry + uint16(chk & 0x0ffff))
	h.HdrChecksum = Htons(csum)
}
