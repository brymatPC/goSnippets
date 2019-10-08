package ipchan

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func doTCPDialOnce(host string, port int, tcpWriteChan, tcpReadChan chan []byte) *net.TCPConn {
	result := (*net.TCPConn)(nil)
	dest := host + ":" + strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp", dest)
	if err != nil {
		fmt.Printf("Could not resolve %v %v  ERROR: %s\r\n", host, port, err.Error())
	} else {
		fmt.Printf("Connecting to [%v]\r\n", dest)
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Printf("Could not connect to: [%v] ERROR: %s\r\n", dest, err.Error())
		} else {
			fmt.Printf("Connected to: [%v]\r\n", dest)
			result = conn
		}
	}
	return result
}

// DoTCPDial connects to a remote host and port
func DoTCPDial(host string, port int, tcpWriteChan, tcpReadChan chan []byte) {
	connectDelay := int64(0)
	for {
		ok := false
		conn := doTCPDialOnce(host, port, tcpWriteChan, tcpReadChan)
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

func doTCPListenOnce(port int, tcpWriteChan, tcpReadChan chan []byte) *net.TCPConn {
	var result *net.TCPConn
	addr := net.TCPAddr{IP: nil, Port: port, Zone: ""}
	listener, err := net.ListenTCP("tcp", &addr)
	defer listener.Close()
	if err != nil {
		fmt.Printf("Could not listen on port: %v  ERROR: %s\r\n", port, err.Error())
		panic(err)
	} else {
		fmt.Printf("Listening on: [%v]\r\n", addr)
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Could not accept connection  ERROR: %s\r\n", err.Error())
		} else {
			fmt.Printf("New connection %v <--> %v\r\n", conn.LocalAddr(), conn.RemoteAddr())
			result = conn
		}
	}
	return result
}

// DoTCPListen listens on the selected port
func DoTCPListen(port int, tcpWriteChan, tcpReadChan chan []byte) {
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
		conn := doTCPListenOnce(port, tcpWriteChan, tcpReadChan)
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
