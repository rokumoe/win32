// +build genapi

package kernel32

/*
func GetProcAddress(module syscall.Handle, procname uintptr) (uintptr, error)
func CreateIoCompletionPort(filehandle syscall.Handle, completionport syscall.Handle, completionkey uintptr, threadnums uint32) (syscall.Handle, error)
func GetQueuedCompletionStatusEx(completionport syscall.Handle, entries []OVERLAPPED_ENTRY, removed *uint32, timeout uint32, alertable bool) error
func PostQueuedCompletionStatus(completionport syscall.Handle, transferred uint32, completionkey uintptr, overlapped *syscall.Overlapped) error
func SetProcessWorkingSetSize(process syscall.Handle, minimumWorkingSetSize uint, maximumWorkingSetSize uint) error
*/

//go:generate go run ../genapi.go
