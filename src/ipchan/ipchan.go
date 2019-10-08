package ipchan

import (
	"fmt"
	"net"
	"time"
)

type ipchan interface {
}

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

func doConnection(conn net.Conn, tcpWriteChan, tcpReadChan chan []byte) (float64, float64, float64) {
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
						fmt.Printf("IP write ERROR: %s\r\n", err.Error())
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
			fmt.Printf("IP read ERROR: %s\r\n", err.Error())
			conn.Close()
			break
		} else if bytesRead == 0 {
			fmt.Printf("IP read got 0 bytes\r\n")
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
	fmt.Printf("IP close sync read: %v written: %v seconds up: %v\r\n", totalBytesRead, totalBytesWritten, elapsed)
	return totalBytesRead, totalBytesWritten, elapsed
}
