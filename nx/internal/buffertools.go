package internal

import "unsafe"

const maxLen = 0xFFFFFFFF

func PointerToByteSlice(ptr unsafe.Pointer, size uint32) []byte {
	return (*[maxLen]byte)(ptr)[:size:size]
}
