package internal

import "unsafe"

// From tinygo
//export memcpy
func Memcpy(dst unsafe.Pointer, src unsafe.Pointer, size uintptr)
