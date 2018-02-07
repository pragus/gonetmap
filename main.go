package main

import (
	"flag"
	"log"
	"os"
	"gonetmap"
)

func ProcessRing(r *gonetmap.NmRing) (uint16) {
	var i uint16
	cur := r.Cur
	for i := 0; !r.RingIsEmpty(); i++ {
		slot_ptr := r.Slot(cur)
		buf_ptr := r.SlotBuffer(slot_ptr)
		_ = buf_ptr
		cur = r.Next(cur)
	}

	return i

}

func main() {
	var (
		iface = flag.String("i", "", "interface")
	)
	flag.Parse()
	if *iface == "" {
		log.Println("usage nm -i netmap:p{0")
		os.Exit(1)
	}
	nm := gonetmap.New()
	if err := nm.COpen(*iface); err != nil {
		log.Println(err)
		return
	}

	ring := nm.OpenRing(uint32(nm.Desc.LastRxRing), false)
	for {
		nm.Poll(-1)
		ProcessRing(ring)
	}
	defer nm.Close()
}
