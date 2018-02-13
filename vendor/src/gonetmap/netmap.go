package gonetmap

import (
	"errors"
	"time"
	"unsafe"
	"syscall"
	"reflect"
	"golang.org/x/sys/unix"
	"github.com/paypal/gatt/linux/gioctl"
	"os"
	"sync"
//	"fmt"
)

const netmapDev  = "/dev/netmap"

const (
	NetmapOpenNoMmap = 0x040000 /* reuse mmap from parent */
	NetmapOpenIfname = 0x080000 /* nr_name, nr_ringid, nr_flags */
	NetmapOpenArg1 = 0x100000
	NetmapOpenArg2 = 0x200000
	NetmapOpenArg3 = 0x400000
	NetmapOpenRingCfg = 0x800000 /* tx|rx rings|slots */
)
const NR_REG_MASK = 0xf
const IOINT = uintptr('i')
const NmReqSize = unsafe.Sizeof(Request{})

var nInfo = gioctl.IoRW(IOINT, 145, NmReqSize) // _IOWR('i', 145, struct nmreq)
var nRegIf = gioctl.IoRW(IOINT, 146, NmReqSize) // _IOWR('i', 146, struct nmreq)
var nTxSync = gioctl.Io(IOINT, 148) // _IO('i', 148) /* sync tx queues */
var nRxSync = gioctl.Io(IOINT, 149) // _IO('i', 149) /* sync rx queues */


var (
	OpenFailed = errors.New("open netmap failed")
	BufIsNull = errors.New("buffer is nil")
	InjectFailed = errors.New("netmap inject failed")
)

func ifaceTOuint16(i interface{}) uint16 {
	var idx uint16
	switch i.(type) {
	case int:
		idx = uint16(i.(int))
	case int16:
		idx = uint16(i.(int16))
	case int32:
		idx = uint16(i.(int32))
	case int64:
		idx = uint16(i.(int64))

	case uint:
		idx = uint16(i.(int))
	case uint16:
		idx = uint16(i.(uint16))
	case uint32:
		idx = uint16(i.(uint32))
	case uint64:
		idx = uint16(i.(uint64))
	}
	return idx
}

func ifaceTOuint32(i interface{}) uint32 {
	var idx uint32
	switch i.(type) {
	case int:
		idx = uint32(i.(int))
	case int16:
		idx = uint32(i.(int16))
	case int32:
		idx = uint32(i.(int32))
	case int64:
		idx = uint32(i.(int64))

	case uint:
		idx = uint32(i.(int))
	case uint16:
		idx = uint32(i.(uint16))
	case uint32:
		idx = uint32(i.(uint32))
	case uint64:
		idx = uint32(i.(uint64))
	}
	return idx
}

type Netmap struct {
	file    *os.File
	MemReg  uintptr
	Pollset []unix.PollFd
	lock    sync.Mutex
}

type Packet struct {
	Time   time.Time // packet time
	Caplen uint32    // bytes stored in the file (caplen <= len)
	Len    uint32    // bytes sent/received
	Data   []byte    // raw packet data
}

type Stat struct {
	Received  uint32
	Dropped   uint32
	IfDropped uint32
}

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

type NetmapIf struct {
	Name       [16]byte
	Version    uint32
	Flags      uint32
	TxRings    uint32
	RxRings    uint32
	BufsHead   uint32
	Spare1     [5]uint32
	RingOffset unsafe.Pointer //NmRing is here
}

type Slot struct {
	Idx   uint32
	Len   uint16
	Flags uint16
	Ptr   uintptr
}

type Request struct {
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

type PacketHeader struct {
	Ts     syscall.Timeval
	Caplen uint32
	Len    uint32
	Flags  uint64
	Desc   *Descriptor
	Slot   *Slot
	Buf    *uint8
}

type Descriptor struct {
	Self        *Descriptor
	Fd          int32
	pad0        [4]byte
	Mem         *byte
	Memsize     uint32
	DoneMmap    int32
	NetmapIf    *NetmapIf
	FirstTxRing uint16
	LastTxRing  uint16
	CurTxRing   uint16
	FirstRxRing uint16
	LastRxRing  uint16
	CurRxRing   uint16
	Request     Request
	Header      PacketHeader
	SomeRing    *NetmapRing
	BufStart    *byte
	BufEnd      *byte
	Snaplen     int32
	Promisc     int32
	ToMs        int32
	pad1        [4]byte
	ErrBuf      *int8
	IfaceFlags  uint32
	IfaceReqcap uint32
	IfaceCurcap uint32
	Stat        Stat
	Msg         [512]byte
}

type Register int

const (
	Default Register = iota        /* backward compat, should not be used. */
	AllNic = iota
	Software = iota
	NicSoftware = iota
	OneNic = iota
	PipeMaster = iota
	PipeSlave = iota
)

type Direction int

const (
	RX Direction = iota
	TX
)

type NetmapBuffer unsafe.Pointer

func (r *Request) SetName(ifname string) {
	copy(r.Name[:], ifname)
}

func ptrSliceFrom(p unsafe.Pointer, s int) (unsafe.Pointer) {
	return unsafe.Pointer(&reflect.SliceHeader{Data:uintptr(p), Len:s, Cap:s})
}

func (r *NetmapRing) GetSlots() (*[]Slot) {
	return (*[]Slot)(ptrSliceFrom(unsafe.Pointer(&r.Slots), int(r.NumSlots)))
}

func (r *NetmapRing) Slot(slot_idx uint32) (*Slot) {
	nm_size := unsafe.Sizeof(r.Slots)
	return (*Slot)(unsafe.Pointer(uintptr(unsafe.Pointer(&r.Slots)) + nm_size * uintptr(slot_idx)))
}

func (r *NetmapRing) Base() (uintptr, uintptr) {
	base_ptr := uintptr(unsafe.Pointer(r)) + r.BufOffset
	buf_size := uintptr(r.BufSize)
	return base_ptr, buf_size

}

func (r *NetmapRing) Next(i uint32) (uint32) {
	i = i + 1

	if i == r.NumSlots {
		i = 0
	}
	r.Cur = i
	r.Head = i
	return i
}

func (r *NetmapRing) GetAvail() (uint32) {
	if r.Tail < r.Cur {
		return r.Tail - r.Cur + r.NumSlots
	} else {
		return r.Tail - r.Cur
	}
}

func (r *NetmapRing) RingIsEmpty() (bool) {
	return (r.Cur == r.Tail)

}

func (r *NetmapRing) SlotBuffer(slot_ptr *Slot) (unsafe.Pointer) {
	idx := uintptr((*slot_ptr).Idx)
	base_ptr := uintptr(unsafe.Pointer(r)) + r.BufOffset
	buf_size := uintptr(r.BufSize)
	ptr := unsafe.Pointer(base_ptr + idx * buf_size)
	return ptr
}

func (r *NetmapRing) BufferSlice(slot_ptr *Slot) (*[]byte) {
	return (*[]byte)(ptrSliceFrom(r.SlotBuffer(slot_ptr), int((*slot_ptr).Len)))
}

func (nif *NetmapIf) ring(idx uint32) (uintptr) {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(nif)) + unsafe.Offsetof(nif.RingOffset))
	h := *(*[]uintptr)(ptrSliceFrom(ptr, int(nif.TxRings + nif.RxRings + 2)))
	return uintptr(unsafe.Pointer(nif)) + h[idx]

}

func (nif *NetmapIf) OpenRing(RingIndex interface{}, direction Direction) (*NetmapRing) {
	var ring_ptr uintptr
	idx := ifaceTOuint32(RingIndex)
	if direction == TX {
		ring_ptr = nif.ring(idx)
	} else {
		ring_ptr = nif.ring(idx + nif.TxRings + 1)
	}

	return (*NetmapRing)(unsafe.Pointer(ring_ptr))
}

func (n *Netmap) RegIf(r *Request) (nif *NetmapIf, err error) {
	if err = gioctl.Ioctl(n.file.Fd(), nRegIf, uintptr(unsafe.Pointer(r))); err == nil {
		//fmt.Printf("ioctl: %v, %+v\n", err, *r)
		if err = n.mmap(r); err == nil {
			nif := (*NetmapIf)(unsafe.Pointer(n.MemReg + uintptr(r.Offset)))
			return nif, err
		}
	}
	//fmt.Printf("ioctl: %v, %+v\n", err, *r)
	return nif, err
}

func (n *Netmap) Info(r *Request) (error) {
	return gioctl.Ioctl(n.file.Fd(), nInfo, uintptr(unsafe.Pointer(r)))
}

func (n *Netmap) mmap(r *Request) (error) {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.MemReg != 0 {
		return nil
	}
	fd := int(n.file.Fd())
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	if data, err := syscall.Mmap(fd, 0, int(r.Memsize), prot, syscall.MAP_SHARED); err == nil {
		n.MemReg = (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data
		return err
	} else {
		return err
	}
}

func New() (*Netmap) {
	file, err := os.OpenFile(netmapDev, os.O_RDWR, 0644)
	if err != nil {
		os.Exit(1)
	}
	return &Netmap{
		MemReg: uintptr(0),
		file:file,
		Pollset:[]unix.PollFd{
			{Fd: int32(file.Fd()), Events:unix.POLLIN, Revents: 0},
		},
	}

}

func (n *Netmap) Poll(timeout int) {
	unix.Poll(n.Pollset, timeout)

}
