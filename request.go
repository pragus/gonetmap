package gonetmap

type Request struct {
	Name    [16]byte
	Version uint32
	Offset  uint32
	Memsize uint32
	TxSlots uint32
	RxSlots uint32
	TxRings uint16
	RxRings uint16
	RingId  uint16
	Cmd     uint16
	Arg1    uint16
	Arg2    uint16
	Arg3    uint32
	Flags   uint32
	Spare2  [1]uint32
}

func (r *Request) SetName(ifname string) {
	copy(r.Name[:], ifname)
}

func NewRequest() *Request {
	return &Request{Version: 11}
}
