package main

import (
	"flag"
	"fmt"
	"github.com/pragus/gonetmap"
	"golang.org/x/sys/unix"
)


const MACFmt = "%#04x, %02x:%02x:%02x:%02x:%02x:%02x -> %02x:%02x:%02x:%02x:%02x:%02x, "

func ProcessSlot(r *gonetmap.NetmapRing, s *gonetmap.Slot) {
	buf := r.SlotBuffer(s)
	eth := (*EtherHdr)(buf)
	s.Flags = gonetmap.RingForward

	fmt.Printf(MACFmt, eth.EtherType,
		eth.SAddr[0], eth.SAddr[1], eth.SAddr[2], eth.SAddr[3], eth.SAddr[4], eth.SAddr[5],
		eth.DAddr[0], eth.DAddr[1], eth.DAddr[2], eth.DAddr[3], eth.DAddr[4], eth.DAddr[5],
	)

	switch eth.EtherType {
	case IPV4Number:
		{
		ip := eth.GetIP()
		fmt.Printf("%+v\n", *ip)
		}
	case ARPNumber:
		arp := eth.GetARP()
		fmt.Printf("%+v\n", *arp)
	default:
		fmt.Printf("\n")

	}

}

func ProcessRing(r *gonetmap.NetmapRing) uint16 {
	var i uint16
	r.Flags = gonetmap.RingForward
	cur := r.Cur
	for i := 0; !r.RingIsEmpty(); i++ {
		SlotPtr := r.Slot(cur)
		ProcessSlot(r, SlotPtr)
		cur = r.Next(cur)
	}

	return i

}

func PollingWorker(nif *gonetmap.Interface, timeout int) {
	fd := int32(nif.File.Fd())
	events := make([]unix.PollFd, 1, 1)
	events[0] = unix.PollFd{Fd: fd, Events: unix.POLLIN, Revents: 0}


	for {
		_, err := unix.Poll(events, timeout)
		if err == nil {
			for ringIndex := uint32(0); ringIndex < nif.Nif.RxRings+1; ringIndex++ {
				ring := nif.GetRing(ringIndex, gonetmap.RX)
				ProcessRing(ring)

			}
			nif.RxSync()
		}

	}
}

func main() {
	rawifaceData := flag.String("iface", "eth0", "interface for transparent filtering")
	flag.Parse()

	netmap := gonetmap.New()
	req0 := gonetmap.Request{Version: 11, RingId: 0, Flags: gonetmap.ReqNicSoftware, Arg1: 0}
	req0.SetName(*rawifaceData)

	iface, _ := netmap.RegIf(&req0)
	PollingWorker(iface, 1)

}
