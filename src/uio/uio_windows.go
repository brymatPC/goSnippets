// +build windows

package uio

import (
	"fmt"
)

type Uio struct {
	dev       string
	uint8ptr  []byte
	uint32ptr []uint32
}

func GetUio(name string) (*Uio, error) {
	return nil, fmt.Errorf("Function not implemented")
}

func ListDevices() {
	fmt.Printf("Function not implemented\r\n")
}
func (uio *Uio) Read(addr uint32) (uint32, error) {
	return 0, fmt.Errorf("Function not implemented")
}

func (uio *Uio) Write(addr uint32, val uint32) error {
	return fmt.Errorf("Function not implemented")
}
