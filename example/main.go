package main

import (
	"syscall"
	"unsafe"
)

func main() {
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("hello world!"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("message box"))),
		0,
	)
}
