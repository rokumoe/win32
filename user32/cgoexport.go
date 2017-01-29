package user32

import (
	"syscall"
	"unsafe"
)

import "C"

//export onEnumWindowsProc
func onEnumWindowsProc(wnd uintptr, lparam uintptr) int32 {
	args := *(**enumargs)(unsafe.Pointer(&lparam))
	return args.enumproc(syscall.Handle(wnd), args.param)
}
