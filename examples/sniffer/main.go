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
	DAddr     [6]byte
	SAddr     [6]byte
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

const MACFmt = "%#04x, %02x:%02x:%02x:%02x:%02x:%02x -> %02x:%02x:%02x:%02x:%02x:%02x "

func ProcessSlot(r *gonetmap.NetmapRing, s *gonetmap.Slot) {
	buf := r.SlotBuffer(s)
	eth := (*EtherHdr)(buf)

	fmt.Printf(MACFmt, eth.EtherType,
		eth.SAddr[0], eth.SAddr[1], eth.SAddr[2], eth.SAddr[3], eth.SAddr[4], eth.SAddr[5],
		eth.DAddr[0], eth.DAddr[1], eth.DAddr[2], eth.DAddr[3], eth.DAddr[4], eth.DAddr[5],
	)

	switch eth.EtherType {
	case [2]byte{0x81, 00}:
		{
			eth := (*EtherVlanHdr)(buf)
			fmt.Printf("tci: %08b proto: %#04x\n", eth.VlanTCI, eth.VlanProto)
		}

	default:
		fmt.Printf("\n")

	}

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
		_, err := unix.Poll(events, timeout)
		if err == nil {
			ProcessRing(ring)
		}

	}
}

func main() {
	rawifaceData := flag.String("iface", "eth0:0", "interface:netmap_ring to listen")
	flag.Parse()

	ifaceData := strings.Split(*rawifaceData, ":")
	ringIndex, _ := strconv.Atoi(ifaceData[1])

	netmap := gonetmap.New()
	req0 := gonetmap.Request{Version: 11, RingId: 0, Flags: gonetmap.ReqAllNic, Arg1: 0}
	req0.SetName(ifaceData[0])

	iface0, _ := netmap.RegIf(&req0)
	rxq0 := iface0.OpenRing(ringIndex, gonetmap.RX)

	PollingWorker(iface0, rxq0, 5)

}
