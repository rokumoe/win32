// +build genapi

package gdi32

/*
func BitBlt(hdc syscall.Handle, x int, y int, cx int, cy int, hdcSrc syscall.Handle, x1 int, y1 int, rop uint32)
func DeleteDC(hdc syscall.Handle)
func CreateCompatibleDC(hdc syscall.Handle) syscall.Handle
func SelectObject(hdc syscall.Handle, hgdiobj syscall.Handle) syscall.Handle
func CreateDIBSection(hdc syscall.Handle, bmi *BITMAPINFO, usage uint32, ppvBits unsafe.Pointer, section syscall.Handle, offset uint32) syscall.Handle
*/

//go:generate go run ../genapi.go
