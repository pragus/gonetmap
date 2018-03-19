package main

import (
	"fmt"
)

// ARPHdr is protocol structure used in Address Resolution Protocol
// for IPv4 to MAC mapping
type ARPHdr struct {
	HType     uint16              // Hardware type, e.g. 1 for Ethernet
	PType     uint16              // Protocol type, e.g. 0x0800 for IPv4
	HLen      uint8               // Hardware address length, e.g. 6 for MAC length
	PLen      uint8               // Protocol address length, e.g. 4 for IPv4 address length
	Operation uint16              // Operation type, see ARP constants
	SHA       [EtherAddrLen]uint8 // Sender hardware address (sender MAC address)
	SPA       [IPv4AddrLen]uint8  // Sender protocol address (sender IPv4 address)
	// array is used to avoid alignment (compiler alignes uint32 on 4 bytes)
	THA [EtherAddrLen]uint8 // Target hardware address (target MAC address)
	TPA [IPv4AddrLen]uint8  // Target protocol address (target IPv4 address)
	// array is used to avoid alignment (compiler alignes uint32 on 4 bytes)
}

// ARP protocol operations
const (
	ARPRequest = 1
	ARPReply   = 2
)

func (hdr *ARPHdr) String() string {
	return fmt.Sprintf(`    L3 protocol: ARP\n
    HType: %d\n
    PType: %d\n
    HLen:  %d\n
    PLen:  %d\n
    Operation: %d\n
    Sender MAC address: %02x:%02x:%02x:%02x:%02x:%02x\n
    Sender IPv4 address: %d.%d.%d.%d\n
    Target MAC address: %02x:%02x:%02x:%02x:%02x:%02x\n
    Target IPv4 address: %d.%d.%d.%d\n`,
		hdr.HType,
		hdr.PType,
		hdr.HLen,
		hdr.PLen,
		hdr.Operation,
		hdr.SHA[0], hdr.SHA[1], hdr.SHA[2], hdr.SHA[3], hdr.SHA[4], hdr.SHA[5],
		hdr.SPA[0], hdr.SPA[1], hdr.SPA[2], hdr.SPA[3],
		hdr.THA[0], hdr.THA[1], hdr.THA[2], hdr.THA[3], hdr.THA[4], hdr.THA[5],
		hdr.TPA[0], hdr.TPA[1], hdr.TPA[2], hdr.TPA[3])
}
