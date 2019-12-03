package mem

import "sync"

type Mem struct {
	memlock   sync.Mutex
	uint32ptr []uint32
	uint8ptr  []uint8
}
