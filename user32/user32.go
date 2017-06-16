// +build genapi

package user32

/*
func EnumWindows(enumproc uintptr, param uintptr) error
func EnumChildWindows(parentwnd syscall.Handle, enumproc uintptr, param uintptr) error
func FindWindowW(classname *uint16, windowname *uint16) (syscall.Handle, error)
func FindWindowExW(parentwnd syscall.Handle, childafter syscall.Handle, classname *uint16, windowname *uint16) (syscall.Handle, error)
func GetClassLongW(wnd syscall.Handle, index int) (uint, error)
func GetClassNameW(wnd syscall.Handle, classname *uint16, maxcount int) int
func GetWindowTextW(wnd syscall.Handle, windowtext *uint16, maxcount int) int
func GetDC(wnd syscall.Handle) (syscall.Handle, error)
func GetDCEx(wnd syscall.Handle, hrgnClip syscall.Handle, flags uint32) (syscall.Handle, error)
func ReleaseDC(wnd syscall.Handle, hdc syscall.Handle)
func GetDesktopWindow() (syscall.Handle, error)
func GetForegroundWindow() (syscall.Handle, error)
func GetSystemMetrics(index int32) (int32, error)
*/

//go:generate go run ../genapi.go
