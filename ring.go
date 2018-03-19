package gonetmap

import (
	"syscall"
	"unsafe"
)

type NetmapRing struct {
	BufOffset uintptr
	NumSlots  uint32
	BufSize   uint16
	RingId    uint16
	Direction uint16
	Head      uint32
	Cur       uint32
	Tail      uint32
	Flags     uint32
	pad0      [4]byte
	Ts        syscall.Timeval
	pad1      [72]byte
	Sem       [128]uint8
	Slots     Slot //NmSlot is here
}

func (r *NetmapRing) GetSlots() *[]Slot {
	return (*[]Slot)(ptrSliceFrom(unsafe.Pointer(&r.Slots), int(r.NumSlots)))
}

func (r *NetmapRing) Slot(slotIdx uint32) *Slot {
	slotUptr := unsafe.Pointer(uintptr(unsafe.Pointer(&r.Slots)) + unsafe.Sizeof(r.Slots)*uintptr(slotIdx))
	return (*Slot)(slotUptr)
}

func (r *NetmapRing) Base() (uintptr, uintptr) {
	base_ptr := uintptr(unsafe.Pointer(r)) + r.BufOffset
	buf_size := uintptr(r.BufSize)
	return base_ptr, buf_size

}

func (r *NetmapRing) Next(i uint32) uint32 {
	i = i + 1

	if i == r.NumSlots {
		i = 0
	}
	r.Cur = i
	r.Head = i
	return i
}

func (r *NetmapRing) GetAvail() uint32 {
	if r.Tail < r.Cur {
		return r.Tail - r.Cur + r.NumSlots
	} else {
		return r.Tail - r.Cur
	}
}

func (r *NetmapRing) RingIsEmpty() bool {
	return r.Cur == r.Tail

}

func (r *NetmapRing) SlotBuffer(slot_ptr *Slot) unsafe.Pointer {
	idx := uintptr((*slot_ptr).Idx)
	base_ptr := uintptr(unsafe.Pointer(r)) + r.BufOffset
	buf_size := uintptr(r.BufSize)
	ptr := unsafe.Pointer(base_ptr + idx*buf_size)
	return ptr
}

func (r *NetmapRing) BufferSlice(slot_ptr *Slot) *[]byte {
	return (*[]byte)(ptrSliceFrom(r.SlotBuffer(slot_ptr), int((*slot_ptr).Len)))
}
