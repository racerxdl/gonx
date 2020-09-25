// +build nintendoswitch

// Named wrappers to Runtime SVC
package svc

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"runtime"
	"time"
	"unsafe"
)

func GetInfo(output *uint64, id0 uint32, handle nxtypes.Handle, id1 uint64) uint64 {
	return runtime.SvcGetInfo(output, id0, uint32(handle), id1)
}

// SendSyncRequest Sends an IPC synchronization request to a session.
// svc 0x21
func SendSyncRequest(session uint64) uint64 {
	return runtime.SvcSendSyncRequest(session)
}

// CloseHandle Closes a handle, decrementing the reference count of the corresponding kernel object.
// This might result in the kernel freeing the object.
// svc 0x16.
func CloseHandle(session nxtypes.Handle) uint64 {
	return runtime.SvcCloseHandle(uint32(session))
}

// ConnectToNamedPort Connects to a registered named port.
// Expects byte to be a null terminated string
// svc 0x1F
func ConnectToNamedPort(session *nxtypes.Handle, name *byte) uint64 {
	r := (*uint32)(unsafe.Pointer(session))
	return runtime.SvcConnectToNamedPort(r, name)
}

// CreateTransferMemory Creates a block of transfer memory.
// svc 0x15
func CreateTransferMemory(handle *nxtypes.Handle, addr uintptr, size uintptr, perm uint32) uint64 {
	r := (*uint32)(unsafe.Pointer(handle))
	return runtime.SvcCreateTransferMemory(r, addr, size, perm)
}

// SetMemoryAttribute Sets memory attributes
// svc 0x03
func SetMemoryAttribute(addr uintptr, size uintptr, state0, state1 uint32) uint64 {
	return runtime.SvcSetMemoryAttribute(addr, size, state0, state1)
}

// WaitSynchronization Waits the specified handles to be finished or the specified timeout
// Returns the Handle Index and error if timeout
// svc 0x18
func WaitSynchronization(handles []nxtypes.Handle, timeout time.Duration) (uint32, error) {
	index := uint32(0)
	handlesPtr := (*uint32)(unsafe.Pointer(&handles[0]))
	r := runtime.SvcWaitSynchronization(&index, handlesPtr, int32(len(handles)), uint64(timeout))
	if r != nxtypes.ResultOK {
		return index, nxerrors.Timeout
	}

	return index, nil
}

// WaitSynchronization Waits a single handle to be finished with the specified timeout
// Calls WaitSynchronization
func WaitSynchronizationSingle(handle nxtypes.Handle, timeout time.Duration) error {
	index := uint32(0)
	handlePtr := (*uint32)(unsafe.Pointer(&handle))
	r := runtime.SvcWaitSynchronization(&index, handlePtr, 1, uint64(timeout))
	if r != nxtypes.ResultOK {
		return nxerrors.Timeout
	}

	return nil
}

func GetIPCBuffer() *[64]uint32 {
	return &runtime.GetTLS().IPCBuffer
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
