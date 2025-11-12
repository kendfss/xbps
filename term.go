package main

import (
	"syscall"
	"unsafe"
)

// winsize holds the dimension and cursor location in of the current
// terminal
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// termSize checks the size of the current environment's terminal
func termSize() (width, height int, err error) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 {
		return 0, 0, errno
	}
	return int(ws.Col), int(ws.Row), nil
}
