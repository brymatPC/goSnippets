package main

import (
	"flag"
	"fmt"
	"goSnippets/src/uio"
	"runtime"
	"time"
)

func main() {
	var help bool
	var displayUio bool

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.BoolVar(&displayUio, "uio", false, "Display available uio")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()
	if help {
		fmt.Printf("UIO Test\n")
		flag.PrintDefaults()
	} else if displayUio {
		uio.ListDevices()
	} else {
		fmt.Printf("UIO Test\n")
		flag.PrintDefaults()
	}
}

func readUio(uioName string, addr int) {
	dev, err := uio.GetUio(uioName)
	if err != nil {
		fmt.Printf("Failed to get uio, error is %v\r\n", err)
		return
	}
	for {
		val, _ := dev.Read(uint32(addr))
		fmt.Printf("Val is %d\r\n", val)
		time.Sleep(1 * time.Second)
	}
}
