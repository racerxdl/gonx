package svc

// Default Handles
const (
	currentProcessHandle = 0xFFFF8001 /// Pseudo handle for the current process.
	currentThreadHandle  = 0xFFFF8000 /// Pseudo handle for the current thread.
)

// TLS is the Nintendo Switch Thread Local Storage
// this struct represents what the TLS space are in Horizon OS
// More details: https://switchbrew.org/wiki/Thread_Local_Region
type TLS struct {
	// IPC command buffer.
	IPCBuffer [0x40]uint32
	// If userland sets this to non-zero, kernel will pin the thread and disallow calls to almost all SVCs.
	DisableCounter uint16
	// If a context switch would have occurred when user disable count was non-zero, kernel will set this to
	// 1. This signifies that the user must call SynchronizePreemptionState to unpin itself and regain access other SVCs
	InterruptFlag uint16
	reserved0     uint32
	reserved1     [0x78]uint8
	tls           [0x50]uint8
	LocalePointer uintptr
	ErrnoVal      uintptr
	ThreadData    uintptr
	EhGlobals     uintptr
	ThreadPointer uintptr
	ThreadType    uintptr
}

// ClearIPCBuffer fills the IPC Buffer with zeroes
func (tls *TLS) ClearIPCBuffer() {
	for i := range tls.IPCBuffer {
		tls.IPCBuffer[i] = 0
	}
}
