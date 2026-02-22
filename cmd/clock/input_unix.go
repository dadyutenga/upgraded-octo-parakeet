//go:build !windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

func readKeys(ch chan<- byte) {
	// Set terminal to raw mode
	fd := int(os.Stdin.Fd())
	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(getTermiosGet()), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return
	}

	newState := oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(getTermiosSet()), uintptr(unsafe.Pointer(&newState)), 0, 0, 0)

	defer syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(getTermiosSet()), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}
		if n == 3 && buf[0] == 0x1b && buf[1] == '[' {
			ch <- buf[2] // 'D' for left, 'C' for right
		} else {
			ch <- buf[0]
		}
	}
}

func getTermiosGet() uintptr {
	return 0x5401 // TCGETS on Linux
}

func getTermiosSet() uintptr {
	return 0x5402 // TCSETS on Linux
}
