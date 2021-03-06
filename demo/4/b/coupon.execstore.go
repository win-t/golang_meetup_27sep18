package main

import (
	"sync"
)

// you cannot pass golang type to C world
// so we use indirection
// we use uint32_t as the key, because uint32_t is subset of C.double (C.lua_Number)

var execstorage = struct {
	sync.RWMutex
	next uint32
	data map[uint32]*execstorageentry
}{
	sync.RWMutex{},
	0,
	make(map[uint32]*execstorageentry),
}

type execstorageentry struct {
	sync.RWMutex
	uid     int
	cdataid int
	s       storage
	data    map[string]string
}

func execstorageGet(index uint32) *execstorageentry {
	execstorage.RLock()
	defer execstorage.RUnlock()

	return execstorage.data[index]
}

func execstorageNew(data *execstorageentry) uint32 {
	execstorage.Lock()
	defer execstorage.Unlock()

	index := execstorage.next
	for {
		if execstorage.data[index] == nil {
			break
		}
		index++
	}
	execstorage.next = index + 1

	execstorage.data[index] = data
	return index
}

func execstorageDelete(index uint32) {
	execstorage.Lock()
	defer execstorage.Unlock()

	delete(execstorage.data, index)
}
