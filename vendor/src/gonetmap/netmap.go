package gonetmap


/*
#include <stdio.h>
#define NETMAP_WITH_LIBS
//#cgo CFLAGS: -DNETMAP_WITH_LIBS -g3
#include <net/netmap_user.h>

int nm_getfd(struct nm_desc *d) {
    return NETMAP_FD(d);
}


size_t* netmap_buf(void* ring, uint32_t index) {return (size_t*)NETMAP_BUF((struct netmap_ring*)ring, index);}

size_t* netmap_rxring(void* nifp, uint32_t index) {return (size_t*)NETMAP_RXRING((struct netmap_if*)nifp, index);}

size_t* netmap_txring(void* nifp, uint32_t index) {return (size_t*)NETMAP_TXRING((struct netmap_if*)nifp, index);}

static struct nm_desc *open_netmap(const char *ifname) {
    printf("%s\n", ifname);
    return nm_open(ifname, NULL, 0, NULL);
}

*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
	"syscall"
	"reflect"
	"fmt"
	"golang.org/x/sys/unix"
)

const (
	NM_OPEN_NO_MMAP = 0x040000 /* reuse mmap from parent */
	NM_OPEN_IFNAME = 0x080000 /* nr_name, nr_ringid, nr_flags */
	NM_OPEN_ARG1 = 0x100000
	NM_OPEN_ARG2 = 0x200000
	NM_OPEN_ARG3 = 0x400000
	NM_OPEN_RING_CFG = 0x800000 /* tx|rx rings|slots */
)

var (
	OPEN_FAILED = errors.New("open netmap failed")
	BUFF_IS_NULL = errors.New("buffer is nil")
	INJECT_FAILED = errors.New("netmap inject failed")
)

type Dummy struct {
	size_t *C.size_t
}

type Netmap struct {
	Desc    *NmDesc
	Fd      int
	Pollset []unix.PollFd
}

type Packet struct {
	Time   time.Time // packet time
	Caplen uint32    // bytes stored in the file (caplen <= len)
	Len    uint32    // bytes sent/received
	Data   []byte    // raw packet data
}

type NmStat struct {
	Received  uint32
	Dropped   uint32
	IfDropped uint32
}

type CNmRing struct {
	Ring *C.struct_netmap_ring
}

type NmRing struct {
	BufOffset   uintptr
	NumSlots    uint32
	Nr_buf_size uint16
	RingId      uint16
	Dir         uint16
	Head        uint32
	Cur         uint32
	Tail        uint32
	Flags       uint32
	pad_cgo_0   [4]byte
	Ts          syscall.Timeval
	pad_cgo_1   [72]byte
	Sem         [128]uint8
	Slots       NmSlot //NmSlot is here
}

type NmIf struct {
	Name       [16]byte
	Version    uint32
	Flags      uint32
	TxRings    uint32
	RxRings    uint32
	BufsHead   uint32
	Spare1     [5]uint32
	RingOffset unsafe.Pointer //NmRing is here
}

type NmSlot struct {
	Idx   uint32
	Len   uint16
	Flags uint16
	Ptr   uintptr
}

type NmReq struct {
	Name    [16]byte
	Version uint32
	Offset  uint32
	Memsize uint32
	TxSlots uint32
	RxSlots uint32
	TxRings uint16
	RxRings uint16
	RingId  uint16
	Cmd     uint16
	Arg1    uint16
	Arg2    uint16
	Arg3    uint32
	Flags   uint32
	Spare2  [1]uint32
}

type NmPktHdr struct {
	Ts     syscall.Timeval
	Caplen uint32
	Len    uint32
	Flags  uint64
	Desc   *NmDesc
	Slot   *NmSlot
	Buf    *uint8
}

type NmDesc struct {
	Self        *NmDesc
	Fd          int32
	pad_cgo_0   [4]byte
	Mem         *byte
	Memsize     uint32
	DoneMmap    int32
	NmIf        *NmIf
	FirstTxRing uint16
	LastTxRing  uint16
	CurTxRing   uint16
	FirstRxRing uint16
	LastRxRing  uint16
	CurRxRing   uint16
	Req         NmReq
	Hdr         NmPktHdr
	SomeRing    *NmRing
	BufStart    *byte
	BufEnd      *byte
	Snaplen     int32
	Promisc     int32
	ToMs        int32
	pad_cgo_1   [4]byte
	ErrBuf      *int8
	IfaceFlags  uint32
	IfaceReqcap uint32
	IfaceCurcap uint32
	Stat        NmStat
	Msg         [512]byte
}

func PtrSliceFrom(p unsafe.Pointer, s int) (unsafe.Pointer) {
	return unsafe.Pointer(&reflect.SliceHeader{Data:uintptr(p), Len:s, Cap:s})
}

func (r *NmRing) GetSlots() (*[]NmSlot) {
	return (*[]NmSlot)(PtrSliceFrom(unsafe.Pointer(&r.Slots), int(r.NumSlots)))
}

func (r *NmRing) Slot(slot_idx uint32) (*NmSlot) {
	nm_size := unsafe.Sizeof(r.Slots)
	return (*NmSlot)(unsafe.Pointer(uintptr(unsafe.Pointer(&r.Slots)) + nm_size * uintptr(slot_idx)))
}

func (r *NmRing) Base() (uintptr, uintptr) {
	base_ptr := uintptr(unsafe.Pointer(r)) + r.BufOffset
	buf_size := uintptr(r.Nr_buf_size)
	return base_ptr, buf_size

}

func (r *NmRing) Next(i uint32) (uint32) {
	i = i + 1

	if i == r.NumSlots {
		i = 0
	}
	r.Cur = i
	r.Head = i
	return i
}

func (r *NmRing) GetAvail() (uint32) {
	if r.Tail < r.Cur {
		return r.Tail - r.Cur + r.NumSlots
	} else {
		return r.Tail - r.Cur
	}
}

func (r *NmRing) RingIsEmpty() (bool) {
	return (r.Cur == r.Tail)

}

func (r *NmRing) SlotBuffer(slot_ptr *NmSlot) (unsafe.Pointer) {
	idx := uintptr((*slot_ptr).Idx)
	base_ptr := uintptr(unsafe.Pointer(r)) + r.BufOffset
	buf_size := uintptr(r.Nr_buf_size)
	ptr := unsafe.Pointer(base_ptr + idx * buf_size)
	return ptr
}

func (r *NmRing) BufferSlice(slot_ptr *NmSlot) (*[]byte) {
	return (*[]byte)(PtrSliceFrom(r.SlotBuffer(slot_ptr), int((*slot_ptr).Len)))
}

func BaseBuf(buf_base_ptr uintptr, buf_size uintptr, idx uint32) (unsafe.Pointer) {
	return unsafe.Pointer(buf_base_ptr + uintptr(idx) * buf_size)
}

func (nif *NmIf) ring(idx uint32) (uintptr) {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(nif)) + unsafe.Offsetof(nif.RingOffset))
	h := *(*[]uintptr)(PtrSliceFrom(ptr, int(nif.TxRings + nif.RxRings + 2)))
	return uintptr(unsafe.Pointer(nif)) + h[idx]

}

func OpenRingbyNif(nif *NmIf, idx uint32, tx bool) (*NmRing) {
	dbg := false
	var ring_ptr uintptr
	var ring_cptr unsafe.Pointer
	if tx {
		ring_ptr = nif.ring(idx)
		ring_cptr = unsafe.Pointer(C.netmap_txring(unsafe.Pointer(nif), C.uint32_t(idx)))
	} else {
		ring_ptr = nif.ring(idx + nif.TxRings + 1)
		ring_cptr = unsafe.Pointer(C.netmap_rxring(unsafe.Pointer(nif), C.uint32_t(idx)))

	}

	if dbg {
		fmt.Printf("idx: %x, tx: %v, base:%p ring: %x, cring: %x\n", idx, tx, nif, ring_ptr, ring_cptr)
	}
	return (*NmRing)(unsafe.Pointer(ring_ptr))
}

func New() (*Netmap) {
	return new(Netmap)

}

func (n *Netmap) OpenRing(idx uint32, tx bool) (*NmRing) {
	return OpenRingbyNif(n.Desc.NmIf, idx, tx)

}

func (n *Netmap) COpen(device string) (err error) {
	dev := C.CString(device)
	defer C.free(unsafe.Pointer(dev))

	n.Desc = (*NmDesc)(unsafe.Pointer(C.nm_open(dev, nil, 0, nil)))
	if n.Desc == nil {
		return OPEN_FAILED
	}
	n.Fd = int(n.Desc.Fd)
	n.Pollset = []unix.PollFd{{Fd: n.Desc.Fd, Events:unix.POLLIN, Revents: 0}}
	return

}

func (n *Netmap) Close() {
	C.nm_close(((*C.struct_nm_desc))(unsafe.Pointer(&n.Desc)))
}

func (n *Netmap) Poll(timeout int) {
	unix.Poll(n.Pollset, timeout)

}
