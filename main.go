package main

import (
	"flag"
	"log"
	"os"
	"gonetmap"

)



func ProcessRing(r *gonetmap.NetmapRing) (uint16) {
	var i uint16
	cur := r.Cur
	for i := 0; !r.RingIsEmpty(); i++ {
		SlotPtr := r.Slot(cur)
		buf_ptr := r.SlotBuffer(SlotPtr)
		_ = buf_ptr
		cur = r.Next(cur)
	}

	return i

}


func main() {
	iface := flag.String("i", "", "interface")
	flag.Parse()

	if *iface == "" {
		log.Println("usage nm -i netmap:p{0")
		os.Exit(1)
	}

	nm := gonetmap.New()



	r0 := gonetmap.Request{Version:11, RingId:0, Flags:gonetmap.PipeMaster}
	r0.Arg1 = 1
	r0.SetName("p")

	n0, _ := nm.RegIf(&r0)
	rx0 := n0.OpenRing(0, gonetmap.RX)
	//fmt.Printf("%+v\n\n", *rx0)


	//os.Exit(1)
	//fmt.Printf("name:%+v %+v\n\n", s, *ring)
	//os.Exit(0)
	//ring := nm.OpenRing(nm.Descriptor.LastRxRing, gonetmap.RX)
	//fmt.Printf("%p=%+v\n\n", &nm.Descriptor, nm.Descriptor)
	//ring := nm.OpenRing(nm.Descriptor.LastRxRing, gonetmap.RX)

	for {
		nm.Poll(-1)
		ProcessRing(rx0)

	}

}
