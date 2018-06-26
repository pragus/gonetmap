package main

import (
	"github.com/pragus/gonetmap"
	"golang.org/x/sys/unix"
)

type EtherHdr struct {
	DAddr     [6]uint8
	SAddr     [6]uint8
	EtherType [2]uint8
}

func ProcessSlot(r *gonetmap.NetmapRing, s *gonetmap.Slot) {
	buf := r.SlotBuffer(s)
	eth := (*EtherHdr)(buf)
	eth.EtherType = [2]byte{0x81, 00}
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
	events := make([]unix.PollFd, 1)
	events[0] = unix.PollFd{Fd: fd, Events: unix.POLLIN, Revents: 0}
	for {
		unix.Poll(events, timeout)
		ProcessRing(ring)

	}
}

func main() {
	netmap := gonetmap.New()
	req0 := gonetmap.Request{Version: 11, RingId: 0, Flags: gonetmap.ReqPipeMaster, Arg1: 0}
	req0.SetName("p")

	iface0, _ := netmap.RegIf(&req0)
	rxq0 := iface0.OpenRing(0, gonetmap.RX)
	PollingWorker(iface0, rxq0, 0)

}
