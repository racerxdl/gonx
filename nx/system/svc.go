// +build nintendoswitch

package system

import (
    _ "github.com/racerxdl/gonx"
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
//export svcBreak
func SvcBreak(reason uint32, a, b uint64)

// Result svcOutputDebugString(const char *str, u64 size)
//export svcOutputDebugString
func svcOutputDebugString(str unsafe.Pointer, size uint64) int64


// SvcOutputDebugString outputs a debug string on Emulator Console
func SvcOutputDebugString(data string) Result {
    sh := (*reflect.StringHeader)(unsafe.Pointer(&data))
    return Result(svcOutputDebugString(unsafe.Pointer(sh.Data), uint64(sh.Len)))
}

func Printf(msg string) {
    msg = msg + "\x00"
    sh := (*reflect.StringHeader)(unsafe.Pointer(&msg))
    printf(unsafe.Pointer(sh.Data))
}

//export printf
func printf(str unsafe.Pointer)

////e    xport nxzinit
//var zinit = func() {
//    SvcBreak(1234,1234,1234)
//}