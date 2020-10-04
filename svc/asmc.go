package svc

import "github.com/racerxdl/gonx/nx/nxtypes"

// All these functions are packaged inside tinygo src

// This function is inside tinygo
//export svcConnectToNamedPort
func svcConnectToNamedPort(session *nxtypes.Handle, name *byte) uint64

// SvcCreateTransferMemory Creates a block of transfer memory.
// svc 0x15
// This function is inside tinygo
//export svcCreateTransferMemory
func svcCreateTransferMemory(handle *nxtypes.Handle, addr uintptr, size uintptr, perm uint32) uint64

// SvcWaitSynchronization Waits on one or more synchronization objects, optionally with a timeout.
// handleCount must not be greater than 40. This is a Horizon Kernel Limitation
// svc 0x18
//export svcWaitSynchronization
func svcWaitSynchronization(index *uint32, handles *nxtypes.Handle, handleCount int32, timeout uint64) uint64

// GetTLS returns a pointer to thread local storage
//export getTLS
func GetTLS() *TLS
