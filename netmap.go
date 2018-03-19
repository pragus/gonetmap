package gonetmap

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"
)

var (
	OpenFailed   = errors.New("open netmap failed")
	BufIsNull    = errors.New("buffer is nil")
	InjectFailed = errors.New("netmap inject failed")
)

type Netmap struct {
	File   *os.File
	MemReg uintptr
	lock   sync.Mutex
}

func (n *Netmap) RegIf(r *Request) (i *Interface, err error) {
	file, _ := n.open()
	if err = nmRegIf(file, r); err == nil {
		fmt.Printf("ioctl: %v, fd: %v, %+v\n", err, file.Fd(), *r)
		if err = n.mmap(file, r); err == nil {
			nif := (*NetmapIf)(unsafe.Pointer(n.MemReg + uintptr(r.Offset)))
			return &Interface{Nif: nif, File: file}, err
		}
	}
	fmt.Printf("ioctl: %v, %+v\n", err, *r)
	return i, err
}

func (n *Netmap) mmap(file *os.File, r *Request) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.MemReg != 0 {
		return nil
	}

	if ptr, err := nmMmap(file, r); err == nil {
		n.MemReg = ptr
		return nil
	} else {
		return err
	}
}

func (n *Netmap) open() (*os.File, error) {
	file, err := os.OpenFile(netmapDev, os.O_RDWR, 0644)
	if err == nil {
		if n.File == nil {
			n.lock.Lock()
			n.File = file
			n.lock.Unlock()
		}
	}
	return file, err
}

func (n *Netmap) Info(r *Request) error {
	return nmInfo(n.File, r)
}

func New() *Netmap {
	return &Netmap{File: nil, MemReg: uintptr(0)}
}
