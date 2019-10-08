package main

import (
	"bufio"
	"flag"
	"fmt"
	"ipchan"
	"os"
	"runtime"
	"time"
)

var compileDate string

func main() {
	var port int
	var isUDP bool
	var help bool

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.IntVar(&port, "port", 9200, "IP port")
	flag.BoolVar(&isUDP, "udp", false, "Select UDP instead of TCP")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()
	if help {
		fmt.Printf("%s compileDate: %v\n", time.Now().Local(), compileDate)
		fmt.Printf("IP server - console to ip\n")
		flag.PrintDefaults()
	} else {
		fmt.Printf("compileDate: %v compiler: %v port: %v\r\n", compileDate, runtime.Version(), port)
		toIPChan := make(chan []byte)
		fromIPChan := make(chan []byte)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Console read  ERROR: %s", err.Error())
					//break
				} else {
					toIPChan <- []byte(text)
				}
			}
		}()
		go func() {
			for {
				buf := <-fromIPChan
				fmt.Printf("%v", string(buf))
			}
		}()

		if !isUDP {
			ipchan.DoTCPListen(port, toIPChan, fromIPChan)
		} else {
			ipchan.DoUDPListen(port, toIPChan, fromIPChan)
		}
	}
}
