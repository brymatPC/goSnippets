package ipchan_test

import (
	"ipchan"
	"testing"
	"time"
)

const (
	tcpAddr = "localhost"
	tcpPort = 9090
)

func TestTCPClientConnection(t *testing.T) {
	serverToIPChan := make(chan []byte)
	serverFromIPChan := make(chan []byte)
	clientToIPChan := make(chan []byte)
	clientFromIPChan := make(chan []byte)

	go ipchan.DoTCPListen(tcpPort, serverToIPChan, serverFromIPChan)
	go ipchan.DoTCPDial(tcpAddr, tcpPort, clientToIPChan, clientFromIPChan)

	clientToIPChan <- []byte("test")

	timeoutTimer := time.NewTimer(50 * time.Millisecond)

	select {
	case input := <-serverFromIPChan:
		if string(input) != "test" {
			t.Errorf("Didn't receive correct bytes")
		}
	case <-timeoutTimer.C:
		t.Errorf("Didn't receive response within 50 milliseconds")
	}
}

func TestTCPServerConnection(t *testing.T) {
	serverToIPChan := make(chan []byte)
	serverFromIPChan := make(chan []byte)
	clientToIPChan := make(chan []byte)
	clientFromIPChan := make(chan []byte)

	go ipchan.DoTCPListen(tcpPort, serverToIPChan, serverFromIPChan)
	go ipchan.DoTCPDial(tcpAddr, tcpPort, clientToIPChan, clientFromIPChan)

	serverToIPChan <- []byte("test")

	timeoutTimer := time.NewTimer(50 * time.Millisecond)

	select {
	case input := <-clientFromIPChan:
		if string(input) != "test" {
			t.Errorf("Didn't receive correct bytes")
		}
	case <-timeoutTimer.C:
		t.Errorf("Didn't receive response within 50 milliseconds")
	}
}
