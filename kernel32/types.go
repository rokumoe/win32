package kernel32

import "syscall"

type OVERLAPPED_ENTRY struct {
	CompletionKey uintptr
	Overlapped    *syscall.Overlapped
	Internal      uintptr
	Transferred   uint32
}
