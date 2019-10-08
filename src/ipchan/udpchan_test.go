package ipchan_test

import (
	"ipchan"
	"testing"
	"time"
)

const (
	udpAddr = "localhost"
	udpPort = 9090
)

func TestUDPClientConnection(t *testing.T) {
	serverToIPChan := make(chan []byte)
	serverFromIPChan := make(chan []byte)
	clientToIPChan := make(chan []byte)
	clientFromIPChan := make(chan []byte)

	go ipchan.DoUDPListen(udpPort, serverToIPChan, serverFromIPChan)
	go ipchan.DoUDPDial(udpAddr, udpPort, clientToIPChan, clientFromIPChan)

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
