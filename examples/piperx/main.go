package main

import (
	"github.com/pragus/go-netmap"
	"golang.org/x/sys/unix"
	"unsafe"
)

func ProcessSlot(r *gonetmap.NetmapRing, s *gonetmap.Slot) {
	buf := r.SlotBuffer(s)
	eth := (*EtherHdr)(buf)
	ip := (*IPv4Hdr)(unsafe.Pointer(uintptr(buf) + EtherLen))
	_ = eth.EtherType
	ip.TimeToLive = 64
}

func ProcessRing(r *gonetmap.NetmapRing) uint16 {
	var i uint16
	cur := r.Cur
	for i := 0; !r.RingIsEmpty(); i++ {
		SlotPtr := r.Slot(cur)
		ProcessSlot(r, SlotPtr)
		cur = r.Next(cur)
	}

	return i

}

func PollingWorker(nif *gonetmap.Interface, ring *gonetmap.NetmapRing, timeout int) {
	fd := int32(nif.File.Fd())
	events := make([]unix.PollFd, 1, 1)
	events[0] = unix.PollFd{Fd: fd, Events: unix.POLLIN, Revents: 0}
	for {
		unix.Poll(events, timeout)
		ProcessRing(ring)

	}
}

func main() {
	netmap := gonetmap.New()
	flags := uint32(gonetmap.PipeMaster | gonetmap.NetmapNoTxPoll)
	req0 := gonetmap.Request{Version: 11, RingId: 0, Flags: flags, Arg1: 0}
	req0.SetName("p")

	req1 := gonetmap.Request{Version: 11, RingId: 1, Flags: flags, Arg1: 0}
	req1.SetName("p")

	iface0, _ := netmap.RegIf(&req0)
	iface1, _ := netmap.RegIf(&req1)
	rxq0 := iface0.OpenRing(0, gonetmap.RX)
	rxq1 := iface1.OpenRing(0, gonetmap.RX)

	go PollingWorker(iface0, rxq0, 0)
	PollingWorker(iface1, rxq1, 0)

}
