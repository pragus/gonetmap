package netpacket

type IPv6Hdr struct {
	VtcFlow    [2]byte           // IP version, traffic class & flow label
	PayloadLen [2]byte           // IP packet length - includes sizeof(ip_header)
	Proto      uint8             // Protocol, next header
	HopLimits  uint8             // Hop limits
	SrcAddr    [IPv6AddrLen]byte // IP address of source host
	DstAddr    [IPv6AddrLen]byte // IP address of destination host(s)
}
