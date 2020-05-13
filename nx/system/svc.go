// +build nintendoswitch

package system

import (
    "reflect"
    "unsafe"
)

// InfoType Types of information to use on SvcGetInfo System Call
type InfoType uint32

// Result represents a result state from a System Call
type Result uint32


// Predefined Results
const (
    ResultOK Result = 0
)

// Result svcBreak(u32 breakReason, u64 inval1, u64 inval2);
//go:export svcBreak
func SvcBreak(reason uint32, a, b uint64)

// Result svcOutputDebugString(const char *str, u64 size)
//go:export svcOutputDebugString
func svcOutputDebugString(str unsafe.Pointer, size uint64) Result


// SvcOutputDebugString outputs a debug string on Emulator Console
func SvcOutputDebugString(data string) Result {
    sh := (*reflect.StringHeader)(unsafe.Pointer(&data))
    return svcOutputDebugString(unsafe.Pointer(sh.Data), uint64(sh.Len))
}