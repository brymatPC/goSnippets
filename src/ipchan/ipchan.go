package ipchan

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var maxDataPerRead int

const logInterval = 1000

func init() {
	maxDataPerRead = 2048
}

func nextDelay(delay int64) int64 {
	if delay == 0 {
		delay = 200
	} else {
		if delay < 1000 {
			delay += 200
		} else if delay < 4000 {
			delay += 500
		} else if delay < 10000 {
			delay += 2000
		} else if delay < 60000 {
			delay += 5000
		}
	}
	return delay
}

func doDialOnce(host string, port int, tcpWriteChan, tcpReadChan chan []byte) *net.TCPConn {
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

// DoDial connects to a remote host and port
func DoDial(host string, port int, tcpWriteChan, tcpReadChan chan []byte) {
	connectDelay := int64(0)
	for {
		ok := false
		conn := doDialOnce(host, port, tcpWriteChan, tcpReadChan)
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

func doListenOnce(port int, tcpWriteChan, tcpReadChan chan []byte) *net.TCPConn {
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

// DoListen listens on the selected port
func DoListen(port int, tcpWriteChan, tcpReadChan chan []byte) {
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
		conn := doListenOnce(port, tcpWriteChan, tcpReadChan)
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
func doConnection(conn *net.TCPConn, tcpWriteChan, tcpReadChan chan []byte) (float64, float64, float64) {
	sync := make(chan bool)
	closeChan := make(chan bool, 1)
	startTime := time.Now()
	totalBytesRead := float64(0)
	totalBytesWritten := float64(0)
	go func() {
		currentBytesWritten := float64(0)
		var inbuf []byte
		for notClose := true; notClose; {
			select {
			case inbuf = <-tcpWriteChan:
				break
			case notClose = <-closeChan:
				notClose = false
				break
			}
			if notClose {
				lim := len(inbuf)
				for i := 0; i < lim; {
					writeSize, err := conn.Write(inbuf[i:lim])
					if err != nil {
						fmt.Printf("TCP/IP write ERROR: %s\r\n", err.Error())
						notClose = false
						break
					}
					i += writeSize
				}
				totalBytesWritten += float64(lim)
				currentBytesWritten += float64(lim)
			}
		}
		conn.Close()
		sync <- true
	}()
	buf := make([]byte, maxDataPerRead)
	currentBytesRead := float64(0)
	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("TCP/IP read ERROR: %s\r\n", err.Error())
			conn.Close()
			break
		} else if bytesRead == 0 {
			fmt.Printf("TCP/IP read got 0 bytes\r\n")
		} else {
			totalBytesRead += float64(bytesRead)
			currentBytesRead += float64(bytesRead)
			qbuf := make([]byte, bytesRead)
			copy(qbuf, buf[0:bytesRead])
			tcpReadChan <- qbuf
		}
	}
	closeChan <- true
	<-sync
	elapsed := time.Since(startTime).Seconds()
	fmt.Printf("TCP/IP close sync read: %v written: %v seconds up: %v\r\n", totalBytesRead, totalBytesWritten, elapsed)
	return totalBytesRead, totalBytesWritten, elapsed
}
