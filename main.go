package main

import (
	"flag"
	"log"
	"os"
	"gonetmap"
	"golang.org/x/sys/unix"
	"unsafe"

)

type Eth struct {
	Dst  [6]byte
	Src  [6]byte
	Type uint16
}


type Buffer struct {
	DatPtr unsafe.Pointer
	Slot *gonetmap.NmSlot
}

type Bvec struct {
	Bufs *[]Buffer
	Len  uint16
}

func SwapAddr(eth *Eth) {
	(*eth).Src, (*eth).Dst = (*eth).Dst, (*eth).Src
}

func ProcessVector(vec *Bvec) {
	for i := uint32(0); i < uint32((*vec).Len); i++ {
		eth := (*Eth)((*(*vec).Bufs)[i].DatPtr)
		SwapAddr(eth)
	}
}




func ProcessRing(r *gonetmap.NmRing, vec *Bvec) (uint32) {
	avail := gonetmap.GetAvail(r)
	bufs := (*vec).Bufs
	(*vec).Len = uint16(avail)
	base_ptr, buf_size := gonetmap.NmRingBasePtr(r)
	i := uint32(0)
	cur := r.Cur
	for !gonetmap.RingIsEmpty(r) {
		slot_ptr := gonetmap.PtrSlotRing(r, cur)
		slot := *slot_ptr
		ptr := unsafe.Pointer(base_ptr + uintptr(slot.Idx) * buf_size)

		(*bufs)[i].DatPtr = ptr
		(*bufs)[i].Slot = slot_ptr
		cur = gonetmap.RingNext(r, cur)
		i++

	}
	ProcessVector(vec)

	return avail

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
	nm, err := gonetmap.OpenNetmap(*iface)
	if err != nil {
		log.Println(err)
		return
	}
	i := nm.Desc.LastRxRing
	ring := gonetmap.NetmapRing(nm.Desc.NmIf, uint32(i), false)
	buf := make([]Buffer, ring.NumSlots, ring.NumSlots)
	vec := Bvec{Bufs:&buf, Len:uint16(ring.NumSlots)}
	for {
		unix.Poll(nm.Pollset, -1)
		ProcessRing(ring, &vec)
	}
	defer nm.Close()
}
