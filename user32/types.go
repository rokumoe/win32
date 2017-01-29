package user32

import "syscall"

type EnumWindowsFunc func(wnd syscall.Handle, param uintptr) int32
