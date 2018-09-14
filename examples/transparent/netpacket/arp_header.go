package netpacket

type ARPHdr struct {
	HType     [2]byte             // Hardware type, e.g. 1 for Ethernet
	PType     [2]byte             // Protocol type, e.g. 0x0800 for IPv4
	HLen      uint8               // Hardware address length, e.g. 6 for MAC length
	PLen      uint8               // Protocol address length, e.g. 4 for IPv4 address length
	Operation uint16              // Operation type, see ARP constants
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
