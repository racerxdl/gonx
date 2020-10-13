// +build nintendoswitch

// Named wrappers to Runtime SVC
package svc

import (
	"device/arm64"
	"fmt"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"time"
	"unsafe"
)

// SvcGetInfo Retrieves information about the system, or a certain kernel object.
// svc 0x29
//go:inline
func GetInfo(output *uint64, id0 uint32, handle nxtypes.Handle, id1 uint64) uint64 {
	return uint64(arm64.SVCall4(0x29, output, id0, handle, id1))
}

// SendSyncRequest Sends an IPC synchronization request to a session.
// svc 0x21
//go:inline
func SendSyncRequest(session uint64) uint64 {
	return uint64(arm64.SVCall1(0x21, session))
}

// CloseHandle Closes a handle, decrementing the reference count of the corresponding kernel object.
// This might result in the kernel freeing the object.
// svc 0x16.
//go:inline
func CloseHandle(session nxtypes.Handle) uint64 {
	return uint64(arm64.SVCall1(0x16, uint64(session)))
}

// ConnectToNamedPort Connects to a registered named port.
// Expects byte to be a null terminated string
// svc 0x1F
//go:inline
func ConnectToNamedPort(session *nxtypes.Handle, name *byte) uint64 {
	res := arm64.AsmFull(`
		str {session}, [sp, #-16]!
		mov x1, {name}
		svc 0x1F
		mov {}, x0
		ldr x2, [sp], #16
		str w1, [x2]
	`, map[string]interface{}{
		"name":    uintptr(unsafe.Pointer(name)),
		"session": uintptr(unsafe.Pointer(session)),
	})

	return uint64(res)
}

// CreateTransferMemory Creates a block of transfer memory.
// svc 0x15
//go:inline
func CreateTransferMemory(handle *nxtypes.Handle, addr uintptr, size uintptr, perm uint32) uint64 {
	// X1 => Addr
	// X2 => Size
	// W3 => Memory Perms
	// Output Result W0
	// Output TransferMemoryhandle W1
	res := arm64.AsmFull(`
		mov x1, {addr}
		mov x2, {size}
		mov x3, {perms}
		str {handle}, [sp, #-16]!
		svc 0x15
		ldr x2, [sp], #16
		str w1, [x2]
	`, map[string]interface{}{
		"addr":   addr,
		"size":   size,
		"perms":  perm,
		"handle": uintptr(unsafe.Pointer(handle)),
	})
	return uint64(res)
}

// SetMemoryAttribute Sets memory attributes
// svc 0x03
//go:inline
func SetMemoryAttribute(addr uintptr, size uintptr, mask, value uint32) uint64 {
	return uint64(arm64.SVCall4(0x03, addr, size, mask, value))
}

// WaitSynchronization Waits the specified handles to be finished or the specified timeout
// Returns the Handle Index and error if timeout
// svc 0x18
//go:inline
func WaitSynchronization(handles []nxtypes.Handle, timeout time.Duration) (uint32, error) {
	if len(handles) > 0x40 {
		// HOS Kernel Limit
		return 0, nxerrors.TooManyHandles
	}
	index := ^uint32(0)
	//
	r := arm64.AsmFull(`
		str {index}, [sp, #-16]!
		mov x1, {handleptr}
		mov x2, {handlesnum}
		mov x3, {timeout}
		svc 0x18
		mov {}, x0
		ldr x2, [sp], #16
		str w1, [x2]
	`, map[string]interface{}{
		"handleptr":  uintptr(unsafe.Pointer(&handles[0])),
		"handlesnum": len(handles),
		"timeout":    uint64(timeout),
		"index":      uintptr(unsafe.Pointer(&index)),
	})

	if r != nxtypes.ResultOK {
		return uint32(index), nxerrors.Timeout
	}

	return uint32(index), nil
}

// WaitSynchronization Waits a single handle to be finished with the specified timeout
// Calls WaitSynchronization
//go:inline
func WaitSynchronizationSingle(handle nxtypes.Handle, timeout time.Duration) error {
	_, err := WaitSynchronization([]nxtypes.Handle{handle}, timeout)
	return err
}

// GetTLS returns a pointer to thread local storage
func GetTLS() *TLS {
	tlsPtr := arm64.AsmFull(`mrs {}, tpidrro_el0`, nil)
	return (*TLS)(unsafe.Pointer(tlsPtr))
}

func GetIPCBuffer() *[64]uint32 {
	return &GetTLS().IPCBuffer
}

func ClearIPCBuffer() {
	buff := GetIPCBuffer()
	for i := 0; i < 64; i++ {
		buff[i] = 0
	}
}

func DumpIPCBuffer() {
	buff := GetIPCBuffer()
	println("TLS IPC Buffer Dump:")
	for i := 0; i < 64; i++ {
		if i%4 == 0 {
			fmt.Printf("\n%04x: ", i*4)
		}
		fmt.Printf("%08x ", buff[i])
	}
	println("")
}
