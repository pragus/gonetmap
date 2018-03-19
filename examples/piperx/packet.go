package main

import (
	"fmt"
	"github.com/intel-go/nff-go/common"
	"unsafe"
)

// EtherHdr L2 header from DPDK: lib/librte_ether/rte_ehter.h
type EtherHdr struct {
	DAddr     [EtherAddrLen]uint8 // Destination address
	SAddr     [EtherAddrLen]uint8 // Source address
	EtherType uint16              // Frame type
}

func (hdr *EtherHdr) String() string {
	r0 := fmt.Sprintf("L2 protocol: Ethernet\nEtherType: 0x%02x\n", hdr.EtherType)
	s := hdr.SAddr
	r1 := fmt.Sprintf("Ethernet Source: %02x:%02x:%02x:%02x:%02x:%02x\n", s[0], s[1], s[2], s[3], s[4], s[5])
	d := hdr.DAddr
	r2 := fmt.Sprintf("Ethernet Destination: %02x:%02x:%02x:%02x:%02x:%02x\n", d[0], d[1], d[2], d[3], d[4], d[5])
	return r0 + r1 + r2
}

// IPv4Hdr L3 header from DPDK: lib/librte_net/rte_ip.h
type IPv4Hdr struct {
	VersionIhl     uint8  // version and header length
	TypeOfService  uint8  // type of service
	TotalLength    uint16 // length of packet
	PacketID       uint16 // packet ID
	FragmentOffset uint16 // fragmentation offset
	TimeToLive     uint8  // time to live
	NextProtoID    uint8  // protocol ID
	HdrChecksum    uint16 // header checksum
	SrcAddr        uint32 // source address
	DstAddr        uint32 // destination address
}

func (hdr *IPv4Hdr) String() string {
	r0 := "    L3 protocol: IPv4\n"
	s := hdr.SrcAddr
	r1 := fmt.Sprintln("    IPv4 Source:", byte(s), ":", byte(s>>8), ":", byte(s>>16), ":", byte(s>>24))
	d := hdr.DstAddr
	r2 := fmt.Sprintln("    IPv4 Destination:", byte(d), ":", byte(d>>8), ":", byte(d>>16), ":", byte(d>>24))
	return r0 + r1 + r2
}

// IPv6Hdr L3 header from DPDK: lib/librte_net/rte_ip.h
type IPv6Hdr struct {
	VtcFlow    uint32             // IP version, traffic class & flow label
	PayloadLen uint16             // IP packet length - includes sizeof(ip_header)
	Proto      uint8              // Protocol, next header
	HopLimits  uint8              // Hop limits
	SrcAddr    [IPv6AddrLen]uint8 // IP address of source host
	DstAddr    [IPv6AddrLen]uint8 // IP address of destination host(s)
}

func (hdr *IPv6Hdr) String() string {
	r0 := "    L3 protocol: IPv6\n"
	s := hdr.SrcAddr
	r1 := fmt.Sprintf("    IPv6 Source: %02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x\n", s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11], s[12], s[13], s[14], s[15])
	d := hdr.DstAddr
	r2 := fmt.Sprintf("    IPv6 Destination %02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x\n", d[0], d[1], d[2], d[3], d[4], d[5], d[6], d[7], d[8], d[9], d[10], d[11], d[12], d[13], d[14], d[15])
	return r0 + r1 + r2
}

// TCPHdr L4 header from DPDK: lib/librte_net/rte_tcp.h
type TCPHdr struct {
	SrcPort  uint16   // TCP source port
	DstPort  uint16   // TCP destination port
	SentSeq  uint32   // TX data sequence number
	RecvAck  uint32   // RX data acknowledgement sequence number
	DataOff  uint8    // Data offset
	TCPFlags TCPFlags // TCP flags
	RxWin    uint16   // RX flow control window
	Cksum    uint16   // TCP checksum
	TCPUrp   uint16   // TCP urgent pointer, if any
}

func (hdr *TCPHdr) String() string {
	r0 := "        L4 protocol: TCP\n"
	r1 := fmt.Sprintf("        L4 Source: %d\n", SwapBytesUint16(hdr.SrcPort))
	r2 := fmt.Sprintf("        L4 Destination: %d\n", SwapBytesUint16(hdr.DstPort))
	return r0 + r1 + r2
}

// UDPHdr L4 header from DPDK: lib/librte_net/rte_udp.h
type UDPHdr struct {
	SrcPort    uint16 // UDP source port
	DstPort    uint16 // UDP destination port
	DgramLen   uint16 // UDP datagram length
	DgramCksum uint16 // UDP datagram checksum
}

func (hdr *UDPHdr) String() string {
	r0 := "        L4 protocol: UDP\n"
	r1 := fmt.Sprintf("        L4 Source: %d\n", SwapBytesUint16(hdr.SrcPort))
	r2 := fmt.Sprintf("        L4 Destination: %d\n", SwapBytesUint16(hdr.DstPort))
	return r0 + r1 + r2
}

// ICMPHdr L4 header.
type ICMPHdr struct {
	Type       uint8  // ICMP message type
	Code       uint8  // ICMP message code
	Cksum      uint16 // ICMP checksum
	Identifier uint16 // ICMP message identifier in some messages
	SeqNum     uint16 // ICMP message sequence number in some messages
}

func (hdr *ICMPHdr) String() string {
	r0 := "        L4 protocol: ICMP\n"
	r1 := fmt.Sprintf("        ICMP Type: %d\n", hdr.Type)
	r2 := fmt.Sprintf("        ICMP Code: %d\n", hdr.Code)
	r3 := fmt.Sprintf("        ICMP Cksum: %d\n", SwapBytesUint16(hdr.Cksum))
	r4 := fmt.Sprintf("        ICMP Identifier: %d\n", SwapBytesUint16(hdr.Identifier))
	r5 := fmt.Sprintf("        ICMP SeqNum: %d\n", SwapBytesUint16(hdr.SeqNum))
	return r0 + r1 + r2 + r3 + r4 + r5
}

// Packet is a set of pointers in NFF-GO library. Each pointer points to one of five headers:
// Mac, IPv4, IPv6, TCP and UDP plus raw pointer.
//
// Empty packet means that only raw pointer is not nil: it points to beginning of packet data
// – raw bits. User should extract packet data somehow.
//
// Parsing means to fill required header pointers with corresponding headers. For example,
// after user fills IPv4 pointer to right place inside packet he can use its fields like
// packet.IPv4.SrcAddr or packet.IPv4.DstAddr.
type Packet struct {
	L3   unsafe.Pointer // Pointer to L3 header in mbuf
	L4   unsafe.Pointer // Pointer to L4 header in mbuf
	Data unsafe.Pointer // Pointer to the packet payload data

	// Last two fields of this structure is filled during InitMbuf macros inside low.c file
	// Need to change low.c for all changes in these fields or adding/removing fields before them.
	Ether *EtherHdr      // Pointer to L2 header in mbuf. It is always parsed and point beginning of packet.
	CMbuf unsafe.Pointer // Private pointer to mbuf. Users shouldn't know anything about mbuf

	Next *Packet // non nil if packet consists of several chained mbufs
}

func (packet *Packet) unparsed() unsafe.Pointer {
	ether := unsafe.Pointer(packet.Ether)
	return unsafe.Pointer(uintptr(ether) + EtherLen)
}

// StartAtOffset function return pointer to first byte of packet
// with given offset.
func (packet *Packet) StartAtOffset(offset uintptr) unsafe.Pointer {
	start := unsafe.Pointer(packet.Ether)
	return unsafe.Pointer(uintptr(start) + offset)
}

// ParseL3 set pointer to start of L3 header
func (packet *Packet) ParseL3() {
	packet.L3 = packet.unparsed()
}

// GetIPv4 ensures if EtherType is IPv4 and casts L3 pointer to IPv4Hdr type.
func (packet *Packet) GetIPv4() *IPv4Hdr {
	if packet.Ether.EtherType == SwapBytesUint16(common.IPV4Number) {
		return (*IPv4Hdr)(packet.L3)
	}
	return nil
}

// GetIPv4NoCheck casts L3 pointer to IPv4Hdr type.
func (packet *Packet) GetIPv4NoCheck() *IPv4Hdr {
	return (*IPv4Hdr)(packet.L3)
}

// GetARP ensures if EtherType is ARP and casts L3 pointer to ARPHdr type.
func (packet *Packet) GetARP() *ARPHdr {
	if packet.Ether.EtherType == SwapBytesUint16(common.ARPNumber) {
		return (*ARPHdr)(packet.L3)
	}
	return nil
}

// GetARPNoCheck casts L3 pointer to ARPHdr type.
func (packet *Packet) GetARPNoCheck() *ARPHdr {
	return (*ARPHdr)(packet.L3)
}

// GetIPv6 ensures if EtherType is IPv6 and cast L3 pointer to IPv6Hdr type.
func (packet *Packet) GetIPv6() *IPv6Hdr {
	if packet.Ether.EtherType == SwapBytesUint16(IPV6Number) {
		return (*IPv6Hdr)(packet.L3)
	}
	return nil
}

// GetIPv6NoCheck ensures if EtherType is IPv6 and cast L3 pointer to
// IPv6Hdr type.
func (packet *Packet) GetIPv6NoCheck() *IPv6Hdr {
	return (*IPv6Hdr)(packet.L3)
}

// ParseL4ForIPv4 set L4 to start of L4 header, if L3 protocol is IPv4.
func (packet *Packet) ParseL4ForIPv4() {
	packet.L4 = unsafe.Pointer(uintptr(packet.unparsed()) + uintptr((packet.GetIPv4NoCheck().VersionIhl&0x0f)<<2))
}

// ParseL4ForIPv6 set L4 to start of L4 header, if L3 protocol is IPv6.
func (packet *Packet) ParseL4ForIPv6() {
	packet.L4 = unsafe.Pointer(uintptr(packet.unparsed()) + uintptr(IPv6Len))
}

// GetTCPForIPv4 ensures if L4 type is TCP and cast L4 pointer to TCPHdr type.
func (packet *Packet) GetTCPForIPv4() *TCPHdr {
	if packet.GetIPv4NoCheck().NextProtoID == TCPNumber {
		return (*TCPHdr)(packet.L4)
	}
	return nil
}

// GetTCPNoCheck casts L4 pointer to TCPHdr type.
func (packet *Packet) GetTCPNoCheck() *TCPHdr {
	return (*TCPHdr)(packet.L4)
}

// GetTCPForIPv6 ensures if L4 type is TCP and cast L4 pointer to *TCPHdr type.
func (packet *Packet) GetTCPForIPv6() *TCPHdr {
	if packet.GetIPv6NoCheck().Proto == TCPNumber {
		return (*TCPHdr)(packet.L4)
	}
	return nil
}

// GetUDPForIPv4 ensures if L4 type is UDP and cast L4 pointer to *UDPHdr type.
func (packet *Packet) GetUDPForIPv4() *UDPHdr {
	if packet.GetIPv4NoCheck().NextProtoID == UDPNumber {
		return (*UDPHdr)(packet.L4)
	}
	return nil
}

// GetUDPNoCheck casts L4 pointer to *UDPHdr type.
func (packet *Packet) GetUDPNoCheck() *UDPHdr {
	return (*UDPHdr)(packet.L4)
}

// GetUDPForIPv6 ensures if L4 type is UDP and cast L4 pointer to *UDPHdr type.
func (packet *Packet) GetUDPForIPv6() *UDPHdr {
	if packet.GetIPv6NoCheck().Proto == UDPNumber {
		return (*UDPHdr)(packet.L4)
	}
	return nil
}

// GetICMPForIPv4 ensures if L4 type is ICMP and cast L4 pointer to *ICMPHdr type.
// L3 supposed to be parsed before and of IPv4 type.
func (packet *Packet) GetICMPForIPv4() *ICMPHdr {
	if packet.GetIPv4NoCheck().NextProtoID == ICMPNumber {
		return (*ICMPHdr)(packet.L4)
	}
	return nil
}

// GetICMPNoCheck casts L4 pointer to *ICMPHdr type.
func (packet *Packet) GetICMPNoCheck() *ICMPHdr {
	return (*ICMPHdr)(packet.L4)
}

// GetICMPForIPv6 ensures if L4 type is ICMP and cast L4 pointer to *ICMPHdr type.
// L3 supposed to be parsed before and of IPv6 type.
func (packet *Packet) GetICMPForIPv6() *ICMPHdr {
	if packet.GetIPv6NoCheck().Proto == ICMPNumber {
		return (*ICMPHdr)(packet.L4)
	}
	return nil
}

// ParseAllKnownL3 parses L3 field and returns pointers to parsed headers.
func (packet *Packet) ParseAllKnownL3() (*IPv4Hdr, *IPv6Hdr, *ARPHdr) {
	packet.ParseL3()
	if packet.GetIPv4() != nil {
		return packet.GetIPv4NoCheck(), nil, nil
	} else if packet.GetIPv6() != nil {
		return nil, packet.GetIPv6NoCheck(), nil
	} else if packet.GetARP() != nil {
		return nil, nil, packet.GetARPNoCheck()
	}
	return nil, nil, nil
}

// ParseAllKnownL4ForIPv4 parses L4 field if L3 type is IPv4 and returns pointers to parsed headers.
func (packet *Packet) ParseAllKnownL4ForIPv4() (*TCPHdr, *UDPHdr, *ICMPHdr) {
	packet.ParseL4ForIPv4()
	if packet.GetTCPForIPv4() != nil {
		return packet.GetTCPNoCheck(), nil, nil
	} else if packet.GetUDPForIPv4() != nil {
		return nil, packet.GetUDPNoCheck(), nil
	} else if packet.GetICMPForIPv4() != nil {
		return nil, nil, packet.GetICMPNoCheck()
	}
	return nil, nil, nil
}

// ParseAllKnownL4ForIPv6 parses L4 field if L3 type is IPv6 and returns pointers to parsed headers.
func (packet *Packet) ParseAllKnownL4ForIPv6() (*TCPHdr, *UDPHdr, *ICMPHdr) {
	packet.ParseL4ForIPv6()
	if packet.GetTCPForIPv6() != nil {
		return packet.GetTCPNoCheck(), nil, nil
	} else if packet.GetUDPForIPv6() != nil {
		return nil, packet.GetUDPNoCheck(), nil
	} else if packet.GetICMPForIPv6() != nil {
		return nil, nil, packet.GetICMPNoCheck()
	}
	return nil, nil, nil
}

// ParseL7 fills pointers to all supported headers and data field.
func (packet *Packet) ParseL7(protocol uint) {
	switch protocol {
	case TCPNumber:
		packet.Data = unsafe.Pointer(uintptr(packet.L4) + uintptr(((*TCPHdr)(packet.L4)).DataOff&0xf0)>>2)
	case UDPNumber:
		packet.Data = unsafe.Pointer(uintptr(packet.L4) + uintptr(UDPLen))
	case ICMPNumber:
		packet.Data = unsafe.Pointer(uintptr(packet.L4) + uintptr(ICMPLen))
	}
}

// ParseData parses L3, L4 and fills the field packet.Data.
// returns 0 in case of success and -1 in case of
// failure to parse L3 or L4.
func (packet *Packet) ParseData() int {
	var pktTCP *TCPHdr
	var pktUDP *UDPHdr
	var pktICMP *ICMPHdr

	pktIPv4, pktIPv6, _ := packet.ParseAllKnownL3()
	if pktIPv4 != nil {
		pktTCP, pktUDP, pktICMP = packet.ParseAllKnownL4ForIPv4()
	} else if pktIPv6 != nil {
		pktTCP, pktUDP, pktICMP = packet.ParseAllKnownL4ForIPv6()
	}

	if pktTCP != nil {
		packet.Data = unsafe.Pointer(uintptr(packet.L4) + uintptr(((*TCPHdr)(packet.L4)).DataOff&0xf0)>>2)
	} else if pktUDP != nil {
		packet.Data = unsafe.Pointer(uintptr(packet.L4) + uintptr(UDPLen))
	} else if pktICMP != nil {
		packet.Data = unsafe.Pointer(uintptr(packet.L4) + uintptr(ICMPLen))
	} else {
		return -1
	}
	return 0
}

// ToPacket should be unexported, used in flow package.
func ToPacket(IN uintptr) *Packet {
	return (*Packet)(unsafe.Pointer(IN))
}

// ToUintptr returns start of mbuf for current packet
func (p *Packet) ToUintptr() uintptr {
	return uintptr(unsafe.Pointer(p.CMbuf))
}

// All following functions set Data pointer because it is assumed that user
// need to generate real packets with some information.

// InitEmptyPacket initializes input packet with preallocated plSize of bytes for payload
// and init pointer to Ethernet header.
func InitEmptyPacket(packet *Packet, plSize uint) uint {
	bufSize := plSize + EtherLen
	packet.Data = packet.unparsed()
	return bufSize
}

// InitEmptyIPv4Packet initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet and IPv4 headers.
func InitEmptyIPv4Packet(packet *Packet, plSize uint) uint {
	// TODO After mandatory fields, IPv4 header optionally may have options of variable length
	// Now pre-allocate space only for mandatory fields
	bufSize := plSize + EtherLen + IPv4MinLen

	// After packet is parsed, we can write to packet struct known protocol types
	packet.Ether.EtherType = SwapBytesUint16(IPV4Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv4MinLen)

	// Next fields not required by pktgen to accept packet. But set anyway
	packet.ParseL3()
	packet.GetIPv4NoCheck().VersionIhl = 0x45 // Ipv4, IHL = 5 (min header len)
	packet.GetIPv4NoCheck().TotalLength = SwapBytesUint16(uint16(IPv4MinLen + plSize))
	packet.GetIPv4NoCheck().NextProtoID = NoNextHeader
	return bufSize
}

// InitEmptyIPv6Packet initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet and IPv6 headers.
func InitEmptyIPv6Packet(packet *Packet, plSize uint) uint {
	bufSize := plSize + EtherLen + IPv6Len
	packet.Ether.EtherType = SwapBytesUint16(IPV6Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv6Len)

	packet.ParseL3()
	packet.GetIPv6NoCheck().PayloadLen = SwapBytesUint16(uint16(plSize))
	packet.GetIPv6NoCheck().VtcFlow = SwapBytesUint32(0x60 << 24) // IP version
	packet.GetIPv6NoCheck().Proto = NoNextHeader

	return bufSize
}

// InitEmptyARPPacket initializes empty ARP packet
func InitEmptyARPPacket(packet *Packet) uint {
	var bufSize uint = EtherLen + ARPLen
	packet.Ether.EtherType = SwapBytesUint16(ARPNumber)
	return bufSize
}

// InitEmptyIPv4TCPPacket initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet, IPv4 and TCP headers. This function supposes that IPv4 and TCP
// headers have minimum length. In fact length can be higher due to optional fields.
// Now setting optional fields explicitly is not supported.

func InitEmptyIPv4TCPPacket(packet *Packet, plSize uint) bool {
	// Now user cannot set explicitly optional fields, so len of header is supposed to be equal to TCPMinLen
	// TODO support variable header length (ask header length from user)
	packet.Ether.EtherType = SwapBytesUint16(IPV4Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv4MinLen + TCPMinLen)

	// Next fields not required by pktgen to accept packet. But set anyway
	packet.ParseL3()
	packet.GetIPv4NoCheck().NextProtoID = TCPNumber
	packet.GetIPv4NoCheck().VersionIhl = 0x45 // Ipv4, IHL = 5 (min header len)
	packet.GetIPv4NoCheck().TotalLength = SwapBytesUint16(uint16(IPv4MinLen + TCPMinLen + plSize))

	packet.ParseL4ForIPv4()
	packet.GetTCPNoCheck().DataOff = packet.GetTCPNoCheck().DataOff | 0x50 // TODO check
	return true
}

// InitEmptyIPv4UDPPacket initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet, IPv4 and UDP headers. This function supposes that IPv4
// header has minimum length. In fact length can be higher due to optional fields.
// Now setting optional fields explicitly is not supported.
func InitEmptyIPv4UDPPacket(packet *Packet, plSize uint) uint {
	bufSize := plSize + EtherLen + IPv4MinLen + UDPLen
	packet.Ether.EtherType = SwapBytesUint16(IPV4Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv4MinLen + UDPLen)

	// Next fields not required by pktgen to accept packet. But set anyway
	packet.ParseL3()
	packet.GetIPv4NoCheck().NextProtoID = UDPNumber
	packet.GetIPv4NoCheck().VersionIhl = 0x45 // Ipv4, IHL = 5 (min header len)
	packet.GetIPv4NoCheck().TotalLength = SwapBytesUint16(uint16(IPv4MinLen + UDPLen + plSize))

	packet.ParseL4ForIPv4()
	packet.GetUDPNoCheck().DgramLen = SwapBytesUint16(uint16(UDPLen + plSize))
	return bufSize
}

// InitEmptyIPv4ICMPPacket initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet, IPv4 and ICMP headers. This function supposes that IPv4
// header has minimum length. In fact length can be higher due to optional fields.
// Now setting optional fields explicitly is not supported.
func InitEmptyIPv4ICMPPacket(packet *Packet, plSize uint) uint {
	bufSize := plSize + EtherLen + IPv4MinLen + ICMPLen
	packet.Ether.EtherType = SwapBytesUint16(IPV4Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv4MinLen + ICMPLen)

	// Next fields not required by pktgen to accept packet. But set anyway
	packet.ParseL3()
	packet.GetIPv4NoCheck().NextProtoID = ICMPNumber
	packet.GetIPv4NoCheck().VersionIhl = 0x45 // Ipv4, IHL = 5 (min header len)
	packet.GetIPv4NoCheck().TotalLength = SwapBytesUint16(uint16(IPv4MinLen + ICMPLen + plSize))
	packet.ParseL4ForIPv4()
	return bufSize
}

// InitEmptyIPv6TCPPacket initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet, IPv6 and TCP headers. This function supposes that IPv6 and TCP
// headers have minimum length. In fact length can be higher due to optional fields.
// Now setting optional fields explicitly is not supported.
func InitEmptyIPv6TCPPacket(packet *Packet, plSize uint) uint {
	// TODO support variable header length (ask header length from user)
	bufSize := plSize + EtherLen + IPv6Len + TCPMinLen
	packet.Ether.EtherType = SwapBytesUint16(IPV6Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv6Len + TCPMinLen)

	packet.ParseL3()
	packet.GetIPv6NoCheck().Proto = TCPNumber
	packet.GetIPv6NoCheck().PayloadLen = SwapBytesUint16(uint16(TCPMinLen + plSize))
	packet.GetIPv6NoCheck().VtcFlow = SwapBytesUint32(0x60 << 24) // IP version

	packet.ParseL4ForIPv6()
	packet.GetTCPNoCheck().DataOff = packet.GetTCPNoCheck().DataOff | 0x50
	return bufSize
}

// InitEmptyIPv6UDPPacket initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet, IPv6 and UDP headers. This function supposes that IPv6
// header has minimum length. In fact length can be higher due to optional fields.
// Now setting optional fields explicitly is not supported.
func InitEmptyIPv6UDPPacket(packet *Packet, plSize uint) uint {
	bufSize := plSize + EtherLen + IPv6Len + UDPLen
	packet.Ether.EtherType = SwapBytesUint16(IPV6Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv6Len + UDPLen)

	packet.ParseL3()
	packet.GetIPv6NoCheck().Proto = UDPNumber
	packet.GetIPv6NoCheck().PayloadLen = SwapBytesUint16(uint16(UDPLen + plSize))
	packet.GetIPv6NoCheck().VtcFlow = SwapBytesUint32(0x60 << 24) // IP version

	packet.ParseL4ForIPv6()
	packet.GetUDPNoCheck().DgramLen = SwapBytesUint16(uint16(UDPLen + plSize))
	return bufSize
}

// InitEmptyIPv6ICMPPacket initializes input packet with preallocated plSize of bytes for payload
// and init pointers to Ethernet, IPv6 and ICMP headers.
func InitEmptyIPv6ICMPPacket(packet *Packet, plSize uint) uint {
	bufSize := plSize + EtherLen + IPv6Len + ICMPLen
	packet.Ether.EtherType = SwapBytesUint16(IPV6Number)
	packet.Data = unsafe.Pointer(uintptr(packet.unparsed()) + IPv6Len + ICMPLen)

	// Next fields not required by pktgen to accept packet. But set anyway
	packet.ParseL3()
	packet.GetIPv6NoCheck().Proto = ICMPNumber
	packet.GetIPv6NoCheck().PayloadLen = SwapBytesUint16(uint16(UDPLen + plSize))
	packet.GetIPv6NoCheck().VtcFlow = SwapBytesUint32(0x60 << 24) // IP version
	packet.ParseL4ForIPv6()
	return bufSize
}

// SwapBytesUint16 swaps uint16 in Little Endian and Big Endian
func SwapBytesUint16(x uint16) uint16 {
	return x<<8 | x>>8
}

// SwapBytesUint32 swaps uint32 in Little Endian and Big Endian
func SwapBytesUint32(x uint32) uint32 {
	return ((x & 0x000000ff) << 24) | ((x & 0x0000ff00) << 8) | ((x & 0x00ff0000) >> 8) | ((x & 0xff000000) >> 24)
}

// BytesToIPv4 converts four element address to uint32 representation
func BytesToIPv4(a byte, b byte, c byte, d byte) uint32 {
	return uint32(d)<<24 | uint32(c)<<16 | uint32(b)<<8 | uint32(a)
}

// IPv4ToBytes converts four element address to uint32 representation
func IPv4ToBytes(v uint32) [IPv4AddrLen]byte {
	return [IPv4AddrLen]uint8{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}
}
