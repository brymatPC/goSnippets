package ipchan

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func doUDPDialOnce(host string, port int, tcpWriteChan, tcpReadChan chan []byte) *net.UDPConn {
	result := (*net.UDPConn)(nil)
	dest := host + ":" + strconv.Itoa(port)
	addr, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		fmt.Printf("Could not resolve %v %v  ERROR: %s\r\n", host, port, err.Error())
	} else {
		fmt.Printf("Connecting to [%v]\r\n", dest)
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Printf("Could not connect to: [%v] ERROR: %s\r\n", dest, err.Error())
		} else {
			fmt.Printf("Connected to: [%v]\r\n", dest)
			result = conn
		}
	}
	return result
}

// DoUDPDial connects to a remote host and port
func DoUDPDial(host string, port int, tcpWriteChan, tcpReadChan chan []byte) {
	connectDelay := int64(0)
	for {
		ok := false
		conn := doUDPDialOnce(host, port, tcpWriteChan, tcpReadChan)
		if conn != nil {
			bytesRead, bytesWritten, timeUp := doConnection(conn, tcpWriteChan, tcpReadChan)
			ok = bytesRead != 0 || bytesWritten != 0 || timeUp > 3
		}
		if !ok {
			connectDelay = nextDelay(connectDelay)
			fmt.Printf("connectDelay delay: %v mSec\r\n", connectDelay)
			time.Sleep(time.Duration(connectDelay) * time.Millisecond)
		} else {
			connectDelay = 0
		}
	}
}

func doUDPListenOnce(port int, tcpWriteChan, tcpReadChan chan []byte) *net.UDPConn {
	var result *net.UDPConn
	addr := net.UDPAddr{IP: nil, Port: port, Zone: ""}
	listener, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Could not listen on port: %v  ERROR: %s\r\n", port, err.Error())
		panic(err)
	} else {
		fmt.Printf("Listening on: [%v]\r\n", addr)
		result = listener
	}
	return result
}

// DoUDPListen listens on the selected port
func DoUDPListen(port int, tcpWriteChan, tcpReadChan chan []byte) {
	listenDelay := int64(0)
	connectDelay := int64(0)
	connected := make(chan bool)
	for {
		go func() {
			for flag := true; flag; {
				select {
				case <-connected:
					flag = false
				case <-time.After(time.Duration(connectDelay) * time.Millisecond):
					connectDelay = nextDelay(connectDelay)
					fmt.Printf("Listening on port %d %v mSec\r\n", port, connectDelay)
				}
			}
		}()
		ok := false
		conn := doUDPListenOnce(port, tcpWriteChan, tcpReadChan)
		connected <- true
		if conn != nil {
			bytesRead, bytesWritten, timeUp := doConnection(conn, tcpWriteChan, tcpReadChan)
			ok = bytesRead != 0 || bytesWritten != 0 || timeUp > 3
		}
		if !ok {
			listenDelay = nextDelay(listenDelay)
			fmt.Printf("Listening delay: %v sec\r\n", listenDelay)
			time.Sleep(time.Duration(listenDelay) * time.Millisecond)
		} else {
			listenDelay = 0
		}

	}
}
