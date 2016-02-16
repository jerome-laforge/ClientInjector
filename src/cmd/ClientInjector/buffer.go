package main

import "sync"

var leakyBuffer = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1500)
	},
}

func GetBuffer() []byte {
	return leakyBuffer.Get().([]byte)
}

func ReleaseBuffer(b []byte) {
	leakyBuffer.Put(b)
}
