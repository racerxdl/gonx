// +build nintendoswitch

// Named wrappers to Runtime SVC
package svc

import (
	"device/arm64"
	"fmt"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"time"
)

// SvcGetInfo Retrieves information about the system, or a certain kernel object.
// svc 0x29
//go:inline
func GetInfo(output *uint64, id0 uint32, handle nxtypes.Handle, id1 uint64) uint64 {
	return uint64(arm64.SVCall4(0x29, output, id0, handle, id1))
	//return runtime.SvcGetInfo(output, id0, uint32(handle), id1)
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
	return svcConnectToNamedPort(session, name)
}

// CreateTransferMemory Creates a block of transfer memory.
// svc 0x15
//go:inline
func CreateTransferMemory(handle *nxtypes.Handle, addr uintptr, size uintptr, perm uint32) uint64 {
	return svcCreateTransferMemory(handle, addr, size, perm)
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
	index := uint32(0)
	r := svcWaitSynchronization(&index, &handles[0], int32(len(handles)), uint64(timeout))
	if r != nxtypes.ResultOK {
		return index, nxerrors.Timeout
	}

	return index, nil
}

// WaitSynchronization Waits a single handle to be finished with the specified timeout
// Calls WaitSynchronization
//go:inline
func WaitSynchronizationSingle(handle nxtypes.Handle, timeout time.Duration) error {
	_, err := WaitSynchronization([]nxtypes.Handle{handle}, timeout)
	return err
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
