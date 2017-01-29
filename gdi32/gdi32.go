// +build genapi

package gdi32

/*
func BitBlt(hdc syscall.Handle, x int, y int, cx int, cy int, hdcSrc syscall.Handle, x1 int, y1 int, rop uint32)
*/

//go:generate go run ../genapi.go
