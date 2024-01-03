package console

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

var Shutdown = make(chan os.Signal, 1)

func Console() {
	runtime.LockOSThread()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var cmd []string = strings.SplitN(scanner.Text(), " ", 3)
		switch cmd[0] {
		case "Shutdown":
			Shutdown <- os.Interrupt
		case "stop":
			Shutdown <- os.Interrupt
		case "exit":
			Shutdown <- os.Interrupt
		case "GC":
			runtime.GC()
			log.Println("GC invoked")
		case "panic":
			panic("panicked, you told me to :)")
		case "version":
			log.Println("Version: 1.0")
		default:
			log.Println("Unknown command")
		}
	}
}

func Hash() string {
	var SHA string
	file, err := os.Open(os.Args[0])
	if err != nil {
		SHA = "00000000000000000000000000000000"
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		SHA = "00000000000000000000000000000000"
	}
	hBytes := hash.Sum(nil)
	file.Close()
	SHA = hex.EncodeToString(hBytes) //Convert bytes to string
	return SHA
}
