package main

import "unsafe"

const (
	EtherAddrLen = 6
	IPv4AddrLen  = 4
	IPv6AddrLen  = 16
)

const (
	EtherLen   = 14
	VLANLen    = 4
	MPLSLen    = 4
	IPv4MinLen = 20
	IPv6Len    = 40
	ICMPLen    = 8
	TCPMinLen  = 20
	UDPLen     = 8
	ARPLen     = 28
)

// Supported EtherType for L2
var (
	IPV4Number = [2]byte{0x08, 0x00}
	ARPNumber  = [2]byte{0x08, 0x06}
	VLANNumber = [2]byte{0x81, 0x00}
	MPLSNumber = [2]byte{0x88, 0x47}
	IPV6Number = [2]byte{0x86, 0xdd}
)

// Supported L4 types
const (
	ICMPNumber   = 0x01
	IPNumber     = 0x04
	TCPNumber    = 0x06
	UDPNumber    = 0x11
	NoNextHeader = 0x3B
)


// TCPFlags contains set TCP flags.
type TCPFlags uint8

// Constants for values of TCP flags.
const (
	TCPFlagFin = 0x01
	TCPFlagSyn = 0x02
	TCPFlagRst = 0x04
	TCPFlagPsh = 0x08
	TCPFlagAck = 0x10
	TCPFlagUrg = 0x20
	TCPFlagEce = 0x40
	TCPFlagCwr = 0x80
)

// Supported ICMP Types
const (
	ICMPTypeEchoRequest  uint8 = 8
	ICMPTypeEchoResponse uint8 = 0
)


type EtherHdr struct {
	DAddr     [EtherAddrLen]byte
	SAddr     [EtherAddrLen]byte
	EtherType [2]byte
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




type VlanHdr struct {
	VlanTCI   [2]byte
	VlanProto [2]byte
}

type EtherVlanHdr struct {
	EtherHdr
	VlanHdr
}

type ARPHdr struct {
	HType     [2]byte                     // Hardware type, e.g. 1 for Ethernet
	PType     [2]byte                     // Protocol type, e.g. 0x0800 for IPv4
	HLen      uint8                      // Hardware address length, e.g. 6 for MAC length
	PLen      uint8                      // Protocol address length, e.g. 4 for IPv4 address length
	Operation uint16                     // Operation type, see ARP constants
	SHA       [EtherAddrLen]uint8 // Sender hardware address (sender MAC address)
	SPA       [IPv4AddrLen]uint8  // Sender protocol address (sender IPv4 address)
	// array is used to avoid alignment (compiler aligns uint32 on 4 bytes)
	THA [EtherAddrLen]uint8 // Target hardware address (target MAC address)
	TPA [IPv4AddrLen]uint8  // Target protocol address (target IPv4 address)

}

// ARP protocol operations
const (
	ARPRequest = 1
	ARPReply   = 2
)

type IPv4Hdr struct {
	VersionIhl     uint8  // version and header length
	TypeOfService  uint8  // type of service
	TotalLength    uint16 // length of packet
	PacketID       uint16 // packet ID
	FragmentOffset uint16 // fragmentation offset
	TimeToLive     uint8  // time to live
	NextProtoID    uint8  // protocol ID
	HdrChecksum    uint16 // header checksum
	SrcAddr        [IPv4AddrLen]byte // source address
	DstAddr        [IPv4AddrLen]byte // destination address
}

type IPv6Hdr struct {
	VtcFlow    [2]byte             // IP version, traffic class & flow label
	PayloadLen [2]byte             // IP packet length - includes sizeof(ip_header)
	Proto      uint8              // Protocol, next header
	HopLimits  uint8              // Hop limits
	SrcAddr    [IPv6AddrLen]byte // IP address of source host
	DstAddr    [IPv6AddrLen]byte // IP address of destination host(s)
}

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

type UDPHdr struct {
	SrcPort    uint16 // UDP source port
	DstPort    uint16 // UDP destination port
	DgramLen   uint16 // UDP datagram length
	DgramCksum uint16 // UDP datagram checksum
}

type ICMPHdr struct {
	Type       uint8  // ICMP message type
	Code       uint8  // ICMP message code
	Cksum      uint16 // ICMP checksum
	Identifier uint16 // ICMP message identifier in some messages
	SeqNum     uint16 // ICMP message sequence number in some messages
}
