// generated by genapi.go
// GOFILE=user32.go GOPACKAGE=user32
// DO NOT EDIT!
package user32

import (
	"syscall"
	"unsafe"
)

var _ unsafe.Pointer // keep unsafe

var (
	pfnEnumWindows         uintptr
	pfnEnumChildWindows    uintptr
	pfnFindWindowW         uintptr
	pfnFindWindowExW       uintptr
	pfnGetClassLongW       uintptr
	pfnGetClassNameW       uintptr
	pfnGetWindowTextW      uintptr
	pfnGetDC               uintptr
	pfnGetDCEx             uintptr
	pfnGetDesktopWindow    uintptr
	pfnGetForegroundWindow uintptr
)

func mustload(libname string) syscall.Handle {
	hlib, err := syscall.LoadLibrary(libname)
	if err != nil {
		panic(err)
	}
	return hlib
}

var (
	pfngetprocaddress uintptr
)

func mustfind(hmodule syscall.Handle, procname string) uintptr {
	ptr := uintptr(0)
	if procname[0] == '#' {
		for i := 1; i < len(procname); i++ {
			c := procname[i]
			if c < '0' || c > '9' {
				break
			}
			ptr = ptr*10 + uintptr(c-'0')
		}
	} else {
		ptr = *(*uintptr)(unsafe.Pointer(&procname))
	}
	proc, _, err := syscall.Syscall(pfngetprocaddress, 2,
		uintptr(hmodule),
		ptr,
		0)
	if proc == 0 {
		panic(err)
	}
	return proc
}

func boolcast(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func EnumWindows(enumproc uintptr, param uintptr) error {
	_, _, en := syscall.Syscall(pfnEnumWindows, 2,
		enumproc,
		param,
		0)
	var err error
	if en != 0 {
		err = en
	}
	return err
}

func EnumChildWindows(parentwnd syscall.Handle, enumproc uintptr, param uintptr) error {
	_, _, en := syscall.Syscall(pfnEnumChildWindows, 3,
		uintptr(parentwnd),
		enumproc,
		param,
	)
	var err error
	if en != 0 {
		err = en
	}
	return err
}

func FindWindowW(classname *uint16, windowname *uint16) (syscall.Handle, error) {
	r1, _, en := syscall.Syscall(pfnFindWindowW, 2,
		uintptr(unsafe.Pointer(classname)),
		uintptr(unsafe.Pointer(windowname)),
		0)
	var err error
	if en != 0 {
		err = en
	}
	return syscall.Handle(r1), err
}

func FindWindowExW(parentwnd syscall.Handle, childafter syscall.Handle, classname *uint16, windowname *uint16) (syscall.Handle, error) {
	r1, _, en := syscall.Syscall6(pfnFindWindowExW, 4,
		uintptr(parentwnd),
		uintptr(childafter),
		uintptr(unsafe.Pointer(classname)),
		uintptr(unsafe.Pointer(windowname)),
		0, 0)
	var err error
	if en != 0 {
		err = en
	}
	return syscall.Handle(r1), err
}

func GetClassLongW(wnd syscall.Handle, index int) (uint, error) {
	r1, _, en := syscall.Syscall(pfnGetClassLongW, 2,
		uintptr(wnd),
		uintptr(index),
		0)
	var err error
	if en != 0 {
		err = en
	}
	return uint(r1), err
}

func GetClassNameW(wnd syscall.Handle, classname *uint16, maxcount int) int {
	r1, _, _ := syscall.Syscall(pfnGetClassNameW, 3,
		uintptr(wnd),
		uintptr(unsafe.Pointer(classname)),
		uintptr(maxcount),
	)
	return int(r1)
}

func GetWindowTextW(wnd syscall.Handle, windowtext *uint16, maxcount int) int {
	r1, _, _ := syscall.Syscall(pfnGetWindowTextW, 3,
		uintptr(wnd),
		uintptr(unsafe.Pointer(windowtext)),
		uintptr(maxcount),
	)
	return int(r1)
}

func GetDC(wnd syscall.Handle) (syscall.Handle, error) {
	r1, _, en := syscall.Syscall(pfnGetDC, 1,
		uintptr(wnd),
		0, 0)
	var err error
	if en != 0 {
		err = en
	}
	return syscall.Handle(r1), err
}

func GetDCEx(wnd syscall.Handle, hrgnClip syscall.Handle, flags uint32) (syscall.Handle, error) {
	r1, _, en := syscall.Syscall(pfnGetDCEx, 3,
		uintptr(wnd),
		uintptr(hrgnClip),
		uintptr(flags),
	)
	var err error
	if en != 0 {
		err = en
	}
	return syscall.Handle(r1), err
}

func GetDesktopWindow() (syscall.Handle, error) {
	r1, _, en := syscall.Syscall(pfnGetDesktopWindow, 0,
		0, 0, 0)
	var err error
	if en != 0 {
		err = en
	}
	return syscall.Handle(r1), err
}

func GetForegroundWindow() (syscall.Handle, error) {
	r1, _, en := syscall.Syscall(pfnGetForegroundWindow, 0,
		0, 0, 0)
	var err error
	if en != 0 {
		err = en
	}
	return syscall.Handle(r1), err
}

func init() {
	hkernel32 := mustload("kernel32.dll")
	var err error
	pfngetprocaddress, err = syscall.GetProcAddress(hkernel32, "GetProcAddress")
	if err != nil {
		panic(err)
	}
	huser32 := mustload("user32.dll")
	_ = huser32
	pfnEnumWindows = mustfind(huser32, "EnumWindows\000")
	pfnEnumChildWindows = mustfind(huser32, "EnumChildWindows\000")
	pfnFindWindowW = mustfind(huser32, "FindWindowW\000")
	pfnFindWindowExW = mustfind(huser32, "FindWindowExW\000")
	pfnGetClassLongW = mustfind(huser32, "GetClassLongW\000")
	pfnGetClassNameW = mustfind(huser32, "GetClassNameW\000")
	pfnGetWindowTextW = mustfind(huser32, "GetWindowTextW\000")
	pfnGetDC = mustfind(huser32, "GetDC\000")
	pfnGetDCEx = mustfind(huser32, "GetDCEx\000")
	pfnGetDesktopWindow = mustfind(huser32, "GetDesktopWindow\000")
	pfnGetForegroundWindow = mustfind(huser32, "GetForegroundWindow\000")
}
