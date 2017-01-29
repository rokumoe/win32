package user32

import (
	"syscall"
	"unsafe"
)

/*
#include <windef.h>

int onEnumWindowsProc(HWND hWnd, LPARAM lParam);

BOOL CALLBACK _enumwindows_proc(HWND hWnd, LPARAM lParam) {
	return onEnumWindowsProc(hWnd, lParam) == 1;
}
*/
import "C"

type enumargs struct {
	enumproc EnumWindowsFunc
	param    uintptr
}

func EnumWindowsGo(fn EnumWindowsFunc, param uintptr) error {
	args := &enumargs{
		enumproc: fn,
		param:    param,
	}
	return EnumWindows(uintptr(C._enumwindows_proc), uintptr(unsafe.Pointer(args)))
}

func EnumChildWindowsGo(parentwnd syscall.Handle, fn EnumWindowsFunc, param uintptr) error {
	args := &enumargs{
		enumproc: fn,
		param:    param,
	}
	return EnumChildWindows(parentwnd, uintptr(C._enumwindows_proc), uintptr(unsafe.Pointer(args)))
}
