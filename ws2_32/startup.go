package ws2_32

import (
	"syscall"
)

func init() {
	var data syscall.WSAData
	err := syscall.WSAStartup(0x202, &data)
	if err != nil {
		panic(err)
	}
}
