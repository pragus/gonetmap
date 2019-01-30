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
	Slots     Slot // NmSlot is here
}

func (r *NetmapRing) GetSlots() *[]Slot {
	return (*[]Slot)(PtrSliceFrom(unsafe.Pointer(&r.Slots), int(r.NumSlots)))
}

func (r *NetmapRing) Slot(slotIdx uint32) *Slot {
	slotPtrUnsafe := unsafe.Pointer(uintptr(unsafe.Pointer(&r.Slots)) + unsafe.Sizeof(r.Slots)*uintptr(slotIdx))
	return (*Slot)(slotPtrUnsafe)
}

func (r *NetmapRing) Base() (uintptr, uintptr) {
	BasePtr := uintptr(unsafe.Pointer(r)) + r.BufOffset
	BufSize := uintptr(r.BufSize)
	return BasePtr, BufSize

}

func (r *NetmapRing) Next(i uint32) uint32 {
	i++

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

func (r *NetmapRing) SlotBuffer(slotPtr *Slot) unsafe.Pointer {
	idx := uintptr(slotPtr.Idx)
	BasePtr, BufSize := r.Base()
	return unsafe.Pointer(BasePtr + idx*BufSize)
}

func (r *NetmapRing) BufferSlice(slotPtr *Slot) *[]byte {
	return (*[]byte)(PtrSliceFrom(r.SlotBuffer(slotPtr), int(slotPtr.Len)))
}
