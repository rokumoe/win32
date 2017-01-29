// +build rio

package mswsock

import (
	"syscall"
	"unsafe"

	_ "github.com/vizee/win32/ws2_32"
)

var WSAID_MULTIPLE_RIO = syscall.GUID{0x8509e081, 0x96dd, 0x4005, [8]byte{0xb1, 0x65, 0x9e, 0x2e, 0xe8, 0xc7, 0x9e, 0x3f}}

type RIOBufferId uintptr
type RIOCQ uintptr
type RIORQ uintptr

type RIOResult struct {
	Status         int32
	Transferred    uint32
	SocketContext  uint64
	RequestContext uint64
}

type RIOBuf struct {
	BufferId RIOBufferId
	Offset   uint32
	Length   uint32
}

const (
	RIOInvalidBufferId = RIOBufferId(uintptr(0xFFFFFFFF))
	RIOInvalidCQ       = RIOCQ(0)
	RIOInvalidRQ       = RIORQ(0)

	RIOMaxCQSize = 0x800000
	RIOCorruptCQ = 0xFFFFFFFF
)

const (
	RIOEventCompletion = 1
	RIOIOCPCompletion  = 2
)

type RIONotificationCompletionEvent struct {
	Type        uint32
	EventHandle syscall.Handle
	NotifyReset bool
}

type RIONotificationCompletionIocp struct {
	Type          uint32
	Handle        syscall.Handle
	CompletionKey uintptr
	Overlapped    uintptr
}

type RIOExtensionFunctionTable struct {
	Size                     uint32
	RIOReceive               uintptr
	RIOReceiveEx             uintptr
	RIOSend                  uintptr
	RIOSendEx                uintptr
	RIOCloseCompletionQueue  uintptr
	RIOCreateCompletionQueue uintptr
	RIOCreateRequestQueue    uintptr
	RIODequeueCompletion     uintptr
	RIODeregisterBuffer      uintptr
	RIONotify                uintptr
	RIORegisterBuffer        uintptr
	RIOResizeCompletionQueue uintptr
	RIOResizeRequestQueue    uintptr
}

var riofuncs RIOExtensionFunctionTable

func RIOReceive(rq RIORQ, buf *RIOBuf, nbuf uint32, flags uint32, reqctx uintptr) error {
	b, _, en := syscall.Syscall6(riofuncs.RIOReceive, 5,
		uintptr(rq),
		uintptr(unsafe.Pointer(buf)),
		uintptr(nbuf),
		uintptr(flags),
		reqctx,
		0)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func RIOReceiveEx(rq RIORQ, buf *RIOBuf, nbuf uint32, la *RIOBuf, ra *RIOBuf, ctrlctx *RIOBuf, pflags *RIOBuf, flags uint32, reqctx uintptr) error {
	b, _, en := syscall.Syscall9(riofuncs.RIOReceiveEx, 9,
		uintptr(rq),
		uintptr(unsafe.Pointer(buf)),
		uintptr(nbuf),
		uintptr(unsafe.Pointer(la)),
		uintptr(unsafe.Pointer(ra)),
		uintptr(unsafe.Pointer(ctrlctx)),
		uintptr(unsafe.Pointer(pflags)),
		uintptr(flags),
		reqctx,
	)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func RIOSend(rq RIORQ, buf *RIOBuf, nbuf uint32, flags uint32, reqctx uintptr) error {
	b, _, en := syscall.Syscall6(riofuncs.RIOSend, 5,
		uintptr(rq),
		uintptr(unsafe.Pointer(buf)),
		uintptr(nbuf),
		uintptr(flags),
		reqctx,
		0)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func RIOSendEx(rq RIORQ, buf *RIOBuf, nbuf uint32, la *RIOBuf, ra *RIOBuf, ctrlctx *RIOBuf, pflags *RIOBuf, flags uint32, reqctx uintptr) error {
	b, _, en := syscall.Syscall9(riofuncs.RIOSendEx, 9,
		uintptr(rq),
		uintptr(unsafe.Pointer(buf)),
		uintptr(nbuf),
		uintptr(unsafe.Pointer(la)),
		uintptr(unsafe.Pointer(ra)),
		uintptr(unsafe.Pointer(ctrlctx)),
		uintptr(unsafe.Pointer(pflags)),
		uintptr(flags),
		reqctx,
	)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func RIOCloseCompletionQueue(cq RIOCQ) {
	syscall.Syscall(riofuncs.RIOCloseCompletionQueue, 1, uintptr(cq), 0, 0)
}

func RIOCreateCompletionQueue(queuesize uint32, notificationcompletion uintptr) (RIOCQ, error) {
	r0, _, en := syscall.Syscall(riofuncs.RIOCreateCompletionQueue, 3,
		uintptr(queuesize),
		notificationcompletion,
		0)
	var err error
	if RIOCQ(r0) == RIOInvalidCQ {
		if en == 0 {
			err = syscall.EINVAL
		} else {
			err = en
		}
	}
	return RIOCQ(r0), err
}

func RIOCreateRequestQueue(s syscall.Handle, maxpendingrecv uint32, maxrecvbufs uint32, maxpendingsend uint32, maxsendbufs uint32, recvcq RIOCQ, sendcq RIOCQ, sockctx uintptr) (RIORQ, error) {
	r0, _, en := syscall.Syscall9(riofuncs.RIOCreateRequestQueue, 8,
		uintptr(s),
		uintptr(maxpendingrecv),
		uintptr(maxrecvbufs),
		uintptr(maxpendingsend),
		uintptr(maxsendbufs),
		uintptr(recvcq),
		uintptr(sendcq),
		sockctx,
		0)
	var err error
	if RIORQ(r0) == RIOInvalidRQ {
		if en == 0 {
			err = syscall.EINVAL
		} else {
			err = en
		}
	}
	return RIORQ(r0), err
}

func RIODequeueCompletion(cq RIOCQ, arr *RIOResult, narr uint32) uint32 {
	r0, _, _ := syscall.Syscall(riofuncs.RIODequeueCompletion, 3,
		uintptr(cq),
		uintptr(unsafe.Pointer(arr)),
		uintptr(narr),
	)
	return uint32(r0)
}

func RIODeregisterBuffer(bufid RIOBufferId) {
	syscall.Syscall(riofuncs.RIODeregisterBuffer, 1, uintptr(bufid), 0, 0)
}

func RIONotify(cq RIOCQ) error {
	r0, _, _ := syscall.Syscall(riofuncs.RIONotify, 1, uintptr(cq), 0, 0)
	if r0 == 0 {
		return nil
	}
	return syscall.Errno(r0)
}

func RIORegisterBuffer(buf *byte, buflen uint32) (RIOBufferId, error) {
	r0, _, en := syscall.Syscall(riofuncs.RIORegisterBuffer, 2,
		uintptr(unsafe.Pointer(buf)),
		uintptr(buflen),
		0)
	var err error
	if RIOBufferId(r0) == RIOInvalidBufferId {
		if en == 0 {
			err = syscall.EINVAL
		} else {
			err = en
		}
	}
	return RIOBufferId(r0), err
}

func RIOResizeCompletionQueue(cq RIOCQ, queuesize uint32) error {
	b, _, en := syscall.Syscall(riofuncs.RIOResizeCompletionQueue, 2,
		uintptr(cq),
		uintptr(queuesize),
		0)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func RIOResizeRequestQueue(rq RIORQ, maxpendingrecv uint32, maxpendingsend uint32) error {
	b, _, en := syscall.Syscall(riofuncs.RIOResizeRequestQueue, 3,
		uintptr(rq),
		uintptr(maxpendingrecv),
		uintptr(maxpendingsend),
	)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func init() {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}
	var nret uint32
	err = syscall.WSAIoctl(
		s,
		SIO_GET_MULTIPLE_EXTENSION_FUNCTION_POINTER,
		(*byte)(unsafe.Pointer(&WSAID_MULTIPLE_RIO)),
		uint32(unsafe.Sizeof(syscall.GUID{})),
		(*byte)(unsafe.Pointer(&riofuncs)),
		uint32(unsafe.Sizeof(RIOExtensionFunctionTable{})),
		&nret,
		nil,
		0,
	)
	if err != nil {
		panic(err)
	}
	syscall.Closesocket(s)
}
