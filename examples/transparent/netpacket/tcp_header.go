package netpacket

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
