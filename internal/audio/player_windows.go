//go:build windows

package audio

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	winmm          = syscall.NewLazyDLL("winmm.dll")
	mciSendStringW = winmm.NewProc("mciSendStringW")
)

func mciSend(cmd string) error {
	cmdPtr, _ := syscall.UTF16PtrFromString(cmd)
	ret, _, _ := mciSendStringW.Call(
		uintptr(unsafe.Pointer(cmdPtr)),
		0, 0, 0,
	)
	if ret != 0 {
		return fmt.Errorf("MCI error: %d", ret)
	}
	return nil
}

func playWindows(path string) {
	// Close any previous instance
	mciSend("close azan")

	openCmd := fmt.Sprintf(`open "%s" type mpegvideo alias azan`, path)
	if err := mciSend(openCmd); err != nil {
		return
	}
	if err := mciSend("play azan wait"); err != nil {
		mciSend("close azan")
		return
	}
	mciSend("close azan")
}

func stopPlayback() {
	mciSend("stop azan")
	mciSend("close azan")
}
