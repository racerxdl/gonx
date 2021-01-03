package nxtypes

type Handle uint32               // Resource handle
type SharedMemoryHandle Handle   // Shared Memory handle
type TransferMemoryHandle Handle // Transfer Memory handle
type SessionHandle Handle        // Session handle
type ReventHandle Handle         // Revent handle
type ARUID uint64                // Applet resource user id

// NativeHandle from Android native_handle.h
type NativeHandle struct {
	Version int32
	NumFds  int32
	NumInts int32
}

// Event is a Kernel-mode event structure
type Event struct {
	Revent    Handle
	Wevent    Handle
	AutoClear bool
}
