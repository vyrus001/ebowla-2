package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/binary"
	"flag"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

var seedFilePath, seedFileSufix, payloadFilePath string

func checkFatalErr(reason string, err error) {
	if err != nil {
		println(reason + ":")
		panic(err)
	}
}

func init() {
	flag.StringVar(&seedFilePath, "s", "", "path to file to use as the cryptographic seed")
	flag.StringVar(&payloadFilePath, "p", "", "path to payload file")
	flag.Parse()
	bailOut := func(reason string) {
		println(reason)
		flag.Usage()
		os.Exit(0)
	}
	if len(payloadFilePath) < 1 {
		bailOut("no payload given")
	}
	if len(seedFilePath) < 1 {
		bailOut("no seed file given")
	}
}

func main() {
	seedFile, err := ioutil.ReadFile(seedFilePath)
	checkFatalErr("failed to read seed file", err)
	payload, err := ioutil.ReadFile(payloadFilePath)
	checkFatalErr("failed to read payload file", err)

	rand.Seed(time.Now().Unix())
	offset := rand.Intn(len(seedFile) - 32)
	offsetBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(offsetBytes, uint64(offset))
	hasher := sha512.New()
	hasher.Write(seedFile[offset : offset+32]) // key

	cypher, err := aes.NewCipher(hasher.Sum(nil)[:32])
	checkFatalErr("failed to create cipher object", err)
	gcm, err := cipher.NewGCM(cypher)
	checkFatalErr("failed to create GCM", err)
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	checkFatalErr("failed populate nonce", err)

	payload = gcm.Seal(nonce, nonce, payload, nil)
	payload = append(offsetBytes, payload...)

	checkFatalErr(
		"failed to write payload",
		ioutil.WriteFile("package", payload, 0777),
	)
}
