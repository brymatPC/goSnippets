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
	var help bool

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.IntVar(&port, "port", 9200, "TCP port")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()
	if help {
		fmt.Printf("%s compileDate: %v\n", time.Now().Local(), compileDate)
		fmt.Printf("IP server - console to tcp/ip\n")
		flag.PrintDefaults()
	} else {
		fmt.Printf("compileDate: %v compiler: %v port: %v\r\n", compileDate, runtime.Version(), port)
		toTCPChan := make(chan []byte)
		fromTCPChan := make(chan []byte)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Console read  ERROR: %s", err.Error())
					//break
				} else {
					toTCPChan <- []byte(text)
				}
			}
		}()
		go func() {
			for {
				buf := <-fromTCPChan
				fmt.Printf("%v", string(buf))
			}
		}()
		ipchan.DoListen(port, toTCPChan, fromTCPChan)
	}
}
