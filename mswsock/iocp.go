// +build iocp
package mswsock

import (
	"syscall"
	"unsafe"

	_ "github.com/vizee/win32/ws2_32"
)

var (
	WSAID_TRANSMITFILE         = syscall.GUID{0xb5367df0, 0xcbac, 0x11cf, [8]byte{0x95, 0xca, 0x00, 0x80, 0x5f, 0x48, 0xa1, 0x92}}
	WSAID_ACCEPTEX             = syscall.GUID{0xb5367df1, 0xcbac, 0x11cf, [8]byte{0x95, 0xca, 0x00, 0x80, 0x5f, 0x48, 0xa1, 0x92}}
	WSAID_GETACCEPTEXSOCKADDRS = syscall.GUID{0xb5367df2, 0xcbac, 0x11cf, [8]byte{0x95, 0xca, 0x00, 0x80, 0x5f, 0x48, 0xa1, 0x92}}
	WSAID_CONNECTEX            = syscall.GUID{0x25a207b9, 0xddf3, 0x4660, [8]byte{0x8e, 0xe9, 0x76, 0xe5, 0x8c, 0x74, 0x06, 0x3e}}
	WSAID_DISCONNECTEX         = syscall.GUID{0x7fda2e11, 0x8630, 0x436f, [8]byte{0xa0, 0x31, 0xf5, 0x36, 0xa6, 0xee, 0xc1, 0x57}}
)

type TransmitFileBuffer struct {
	Head       uintptr
	HeadLength uint32
	Tail       uintptr
	TailLength uint32
}

var (
	pfnTransmitFile         uintptr
	pfnAcceptEx             uintptr
	pfnGetAcceptExSockaddrs uintptr
	pfnConnectEx            uintptr
	pfnDisconnectEx         uintptr
)

func TransmitFile(s syscall.Handle, file syscall.Handle, nwrite uint32, bps uint32, overlapped *syscall.Overlapped, buf *TransmitFileBuffer) error {
	b, _, en := syscall.Syscall9(pfnTransmitFile, 7,
		uintptr(s),
		uintptr(file),
		uintptr(nwrite),
		uintptr(bps),
		uintptr(unsafe.Pointer(overlapped)),
		uintptr(unsafe.Pointer(buf)),
		0,
		0, 0)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func AcceptEx(sliten syscall.Handle, saccept syscall.Handle, buf *byte, nrecv uint32, nladdr uint32, nraddr uint32, nreturn *uint32, overlapped *syscall.Overlapped) error {
	b, _, en := syscall.Syscall9(pfnAcceptEx, 8,
		uintptr(sliten),
		uintptr(saccept),
		uintptr(unsafe.Pointer(buf)),
		uintptr(nrecv),
		uintptr(nladdr),
		uintptr(nraddr),
		uintptr(unsafe.Pointer(nreturn)),
		uintptr(unsafe.Pointer(overlapped)),
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

func GetAcceptExSockaddrs(buf *byte, nrecv uint32, nladdr uint32, nraddr uint32, lsa **syscall.RawSockaddrAny, nlsa *int32, rsa **syscall.RawSockaddrAny, nrsa *int32) {
	syscall.Syscall9(pfnGetAcceptExSockaddrs, 8,
		uintptr(unsafe.Pointer(buf)),
		uintptr(nrecv),
		uintptr(nladdr),
		uintptr(nraddr),
		uintptr(unsafe.Pointer(lsa)),
		uintptr(unsafe.Pointer(nlsa)),
		uintptr(unsafe.Pointer(rsa)),
		uintptr(unsafe.Pointer(nrsa)),
		0)
}

func ConnectEx(s syscall.Handle, psa uintptr, salen int32, buf *byte, nsend uint32, nsent *uint32, overlapped *syscall.Overlapped) error {
	b, _, en := syscall.Syscall9(pfnConnectEx, 7,
		uintptr(s),
		psa,
		uintptr(salen),
		uintptr(unsafe.Pointer(buf)),
		uintptr(nsend),
		uintptr(unsafe.Pointer(nsent)),
		uintptr(unsafe.Pointer(overlapped)),
		0, 0)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func DisconnectEx(s syscall.Handle, overlapped *syscall.Overlapped, flags uint32) error {
	b, _, en := syscall.Syscall6(pfnConnectEx, 4,
		uintptr(s),
		uintptr(unsafe.Pointer(overlapped)),
		uintptr(flags),
		0,
		0, 0)
	if b == 0 {
		if en == 0 {
			return syscall.EINVAL
		} else {
			return en
		}
	}
	return nil
}

func mustgetfunc(s syscall.Handle, guid *syscall.GUID) uintptr {
	var (
		pfn  uintptr
		nret uint32
	)
	err := syscall.WSAIoctl(
		s,
		syscall.SIO_GET_EXTENSION_FUNCTION_POINTER,
		(*byte)(unsafe.Pointer(guid)),
		uint32(unsafe.Sizeof(syscall.GUID{})),
		(*byte)(unsafe.Pointer(&pfn)),
		uint32(unsafe.Sizeof(pfn)),
		&nret,
		nil,
		0,
	)
	if err != nil {
		panic(err)
	}
	return pfn
}

func init() {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}

	pfnTransmitFile = mustgetfunc(s, &WSAID_TRANSMITFILE)
	pfnAcceptEx = mustgetfunc(s, &WSAID_ACCEPTEX)
	pfnGetAcceptExSockaddrs = mustgetfunc(s, &WSAID_GETACCEPTEXSOCKADDRS)
	pfnConnectEx = mustgetfunc(s, &WSAID_CONNECTEX)
	pfnDisconnectEx = mustgetfunc(s, &WSAID_DISCONNECTEX)

	syscall.Closesocket(s)
}
