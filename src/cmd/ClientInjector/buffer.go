package main

import "sync"

var leakyBuffer = sync.Pool{
	New: func() interface{} {
		return make([]byte, mtu)
	},
}

func GetBuffer() []byte {
	return leakyBuffer.Get().([]byte)
}

func ReleaseBuffer(b []byte) {
	leakyBuffer.Put(b)
}
