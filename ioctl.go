package gonetmap

import (
	"github.com/paypal/gatt/linux/gioctl"
	"os"
	"unsafe"
)

const IOINT = uintptr('i')
const NmReqSize = unsafe.Sizeof(Request{})

var nInfo = gioctl.IoRW(IOINT, 145, NmReqSize)  // _IOWR('i', 145, struct nmreq)
var nRegIf = gioctl.IoRW(IOINT, 146, NmReqSize) // _IOWR('i', 146, struct nmreq)
var NTxSync = gioctl.Io(IOINT, 148)             // _IO('i', 148) /* sync tx queues */
var NRxSync = gioctl.Io(IOINT, 149)             // _IO('i', 149) /* sync rx queues */

func NmIoctl(file *os.File, op uintptr, arg unsafe.Pointer) error {
	return gioctl.Ioctl(file.Fd(), op, uintptr(unsafe.Pointer(arg)))

}

func nmInfo(file *os.File, r *Request) error {
	return NmIoctl(file, nInfo, unsafe.Pointer(r))
}

func nmRegIf(file *os.File, r *Request) error {
	return NmIoctl(file, nRegIf, unsafe.Pointer(r))
}
