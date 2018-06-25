package main

import (
	"flag"
	"fmt"
	"github.com/pragus/gonetmap"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

type EtherHdr struct {
	DAddr     [6]uint
	SAddr     [6]uint8
	EtherType uint16
}

func ProcessSlot(r *gonetmap.NetmapRing, s *gonetmap.Slot) {
	fmt.Printf("%+v\n", r.BufferSlice(s))
	/*	buf := r.SlotBuffer(s)
		eth := (*EtherHdr)(buf)
		eth.EtherType = 0x800
	*/
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
	rawifaceData := flag.String("iface", "eth0:0", "interface:netmap_ring to listen")
	flag.Parse()

	ifaceData := strings.Split(*rawifaceData, ":")
	ringIndex, _ := strconv.Atoi(ifaceData[1])

	netmap := gonetmap.New()
	req0 := gonetmap.Request{Version: 11, RingId: 0, Flags: gonetmap.ReqNicSoftware, Arg1: 0}
	req0.SetName(ifaceData[0])

	iface0, _ := netmap.RegIf(&req0)
	rxq0 := iface0.OpenRing(ringIndex, gonetmap.RX)

	PollingWorker(iface0, rxq0, 1)

}
