//+build linux

package mem

import (
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

const (
	memBase = 0xFFFF0000
	memSpan = 0x0000F000
)

// Open and memory map memory range from /dev/mem .
// Some reflection magic is used to convert it to a unsafe []uint32 pointer
func (mem *Mem) Open() (err error) {
	var file *os.File
	// Open fd for rw mem access; try dev/mem first (need root)
	file, err = os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return
	}
	// FD can be closed after memory mapping
	defer file.Close()

	mem.memlock.Lock()
	defer mem.memlock.Unlock()

	mem.uint32ptr, ocm.uint8ptr, err = memMap(file.Fd(), memBase, memSpan)
	if err != nil {
		return
	}
	return nil
}

func memMap(fd uintptr, base int64, length int) (mem []uint32, mem8 []byte, err error) {

	mem8, err = syscall.Mmap(
		int(fd),
		base,
		length,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)

	if err != nil {
		return
	}

	//Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&mem8))
	header.Len /= (32 / 8) // (32 bit = 4 bytes)
	header.Cap /= (32 / 8)
	mem = *(*[]uint32)(unsafe.Pointer(&header))
	return
}

// Close unmaps memory
func (mem *Mem) Close() error {
	mem.memlock.Lock()
	defer mem.memlock.Unlock()
	if err := syscall.Munmap(mem.uint8ptr); err != nil {
		return err
	}
	return nil
}
func (mem *Mem) Read(addr uint32) uint32 {
	return mem.uint32ptr[addr]
}
func (mem *Mem) Write(addr uint32, val uint32) {
	mem.uint32ptr[addr] = val
}
func (mem *Mem) ReadByte(addr uint32) byte {
	return mem.uint8ptr[addr]
}
func (mem *Mem) ReadBytes(addr uint32, bytes []byte) int {
	numBytes := copy(bytes, mem.uint8ptr[addr:])
	return numBytes
}
func (mem *Mem) WriteByte(addr uint32, val byte) {
	mem.uint8ptr[addr] = val
}
func (mem *Mem) WriteBytes(addr uint32, bytes []byte) int {
	numBytes := copy(mem.uint8ptr[addr:], bytes)
	return numBytes
}
