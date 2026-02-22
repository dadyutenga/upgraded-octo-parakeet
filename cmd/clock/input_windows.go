//go:build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procReadConsole    = kernel32.NewProc("ReadConsoleInputW")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetStdHandle   = kernel32.NewProc("GetStdHandle")
)

const (
	stdInputHandle       = ^uintptr(0) - 10 + 1 // STD_INPUT_HANDLE = -10
	enableProcessedInput = 0x0001
	enableLineInput      = 0x0002
	enableEchoInput      = 0x0004
	enableWindowInput    = 0x0008
	keyEventType         = 0x0001
)

type inputRecord struct {
	eventType uint16
	_         uint16
	keyDown   int32
	repeat    uint16
	vkCode    uint16
	scanCode  uint16
	char      uint16
	state     uint32
}

func readKeys(ch chan<- byte) {
	handle, _, _ := procGetStdHandle.Call(uintptr(stdInputHandle))
	if handle == 0 || handle == ^uintptr(0) {
		// Fallback: read from stdin byte by byte
		readKeysStdin(ch)
		return
	}

	var mode uint32
	procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	// Disable line input and echo
	newMode := mode &^ (enableLineInput | enableEchoInput)
	procSetConsoleMode.Call(handle, uintptr(newMode))

	defer procSetConsoleMode.Call(handle, uintptr(mode))

	var rec inputRecord
	var read uint32
	for {
		ret, _, _ := procReadConsole.Call(handle, uintptr(unsafe.Pointer(&rec)), 1, uintptr(unsafe.Pointer(&read)), 0)
		if ret == 0 || read == 0 {
			continue
		}
		if rec.eventType != keyEventType || rec.keyDown == 0 {
			continue
		}

		switch rec.vkCode {
		case 0x25: // VK_LEFT
			ch <- 'D'
		case 0x27: // VK_RIGHT
			ch <- 'C'
		case 0x20: // VK_SPACE
			ch <- ' '
		default:
			if rec.char > 0 && rec.char < 128 {
				ch <- byte(rec.char)
			}
		}
	}
}

func readKeysStdin(ch chan<- byte) {
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
