package netpacket

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

type VlanHdr struct {
	VlanTCI   [2]byte
	VlanProto [2]byte
}

type EtherVlanHdr struct {
	EtherHdr
	VlanHdr
}
