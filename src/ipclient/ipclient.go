package main

import (
	"bufio"
	"flag"
	"fmt"
	"ipchan"
	"os"
	"runtime"
)

var compileDate string

func main() {
	var port int
	var host string
	var isUDP bool
	var help bool

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.IntVar(&port, "port", 9200, "TCP port")
	flag.StringVar(&host, "host", "localhost", "Host")
	flag.BoolVar(&isUDP, "udp", false, "Select UDP instead of TCP")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()
	fmt.Printf("IP Client - console to ip\n")

	if help {
		fmt.Printf("compileDate: %v compiler: %v\n\n", compileDate, runtime.Version())
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
					fmt.Printf("Console read  ERROR: %s\r\n", err.Error())
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
			ipchan.DoTCPDial(host, port, toIPChan, fromIPChan)
		} else {
			ipchan.DoUDPDial(host, port, toIPChan, fromIPChan)
		}

	}
}
