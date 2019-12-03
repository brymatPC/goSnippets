//+build windows

package mem

import "fmt"

func (mem *Mem) Open() (err error) {
	return fmt.Errorf("Not implemented")
}

func (mem *Mem) Close() error {
	return fmt.Errorf("Not implemented")
}
func (mem *Mem) Read(addr uint32) uint32 {
	return 0
}
func (mem *Mem) Write(addr uint32, val uint32) {
}
func (mem *Mem) ReadByte(addr uint32) byte {
	return 0
}
func (mem *Mem) ReadBytes(addr uint32, bytes []byte) int {
	return 0
}
func (mem *Mem) WriteByte(addr uint32, val byte) {

}
func (mem *Mem) WriteBytes(addr uint32, bytes []byte) int {
	return 0
}
