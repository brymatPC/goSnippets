//+build linux
package uio

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	uioDir     = "/dev/"
	uioInfoDir = "/sys/class/uio/"
)

type Uio struct {
	dev       string
	uint8ptr  []byte
	uint32ptr []uint32
}

func GetUio(name string) (*Uio, error) {
	uio := Uio{}
	dev, err := findDevice(name)
	if err != nil {
		return nil, err
	}
	uio.dev = dev
	return &uio, nil
}

func memMap(fd uintptr, base int64, length int) (mem []uint32, mem8 []byte, err error) {

	mem8, err = unix.Mmap(
		int(fd),
		base,
		length,
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_SHARED,
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

func (uio *Uio) open() error {
	filename := fmt.Sprintf("%s%s", uioDir, uio.dev)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	uio.uint32ptr, uio.uint8ptr, err = memMap(file.Fd(), 0, 0x2000)
	if err != nil {
		return err
	}
	return nil
}
func (uio *Uio) close() {
	syscall.Munmap(uio.uint8ptr)
}
func (uio *Uio) Read(addr uint32) (uint32, error) {
	err := uio.open()
	if err != nil {
		return 0, err
	}
	defer uio.close()
	return uio.uint32ptr[addr], nil
}
func (uio *Uio) Write(addr uint32, val uint32) error {
	err := uio.open()
	if err != nil {
		return err
	}
	defer uio.close()
	uio.uint32ptr[addr] = val
	return nil
}

func getName(dev string) (string, error) {
	filename := fmt.Sprintf("%s%s/name", uioInfoDir, dev)
	name, err := ioutil.ReadFile(filename)
	return strings.TrimSpace(string(name)), err
}

func getSize(dev string) (int, error) {
	filename := fmt.Sprintf("%s%s/maps/map0/size", uioInfoDir, dev)
	sizeTemp, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("size error %v\r\n", err)
		return 0, err
	}
	sizeStr := strings.TrimSpace(string(sizeTemp))
	size, err := strconv.ParseInt(sizeStr, 16, 32)
	if err != nil {
		fmt.Printf("size error %v\r\n", err)
		return 0, err
	}
	return int(size), nil
}

func findDevice(nameToFind string) (string, error) {
	files, err := ioutil.ReadDir(uioInfoDir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		name, err := getName(file.Name())
		if err == nil {
			if nameToFind == name {
				return file.Name(), nil
			}
		}
	}
	return "", fmt.Errorf("%s not found", nameToFind)
}

func ListDevices() {
	files, err := ioutil.ReadDir(uioInfoDir)
	if err != nil {
		fmt.Printf("Failed to read dir, error is %v\r\n", err)
	}
	for _, file := range files {
		name, err := getName(file.Name())
		if err == nil {
			size, _ := getSize(file.Name())
			fmt.Printf("Uio dev %s, name %s, size %v\r\n", file.Name(), name, size)
		} else {
			fmt.Printf("Failed to get uio name, error is %v\r\n", err)
		}
	}
}
