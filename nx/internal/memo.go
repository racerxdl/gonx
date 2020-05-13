package internal

import "unsafe"

//go:export malloc
func Alloc(size uint64) unsafe.Pointer

//export free
func Free(ptr unsafe.Pointer)
