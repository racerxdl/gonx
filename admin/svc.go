package admin

// SvcDumpInfo Causes the kernel to dump debug information.
//
// This is a privileged syscall.
// Only on firmware 1.0.0 to 3.0.2
// svc 0x3C
//export svcDumpInfo
func SvcDumpInfo(dumpInfoType uint32, arg0 uint64)

// SvcKernelDebug Performs a debugging operation on the kernel.
//
// This is a privileged syscall.
// Only on firmware 4.0.0+
// svc 0x3C
//export svcKernelDebug
func SvcKernelDebug(kernelDebugType uint32, arg0, arg1, arg2 uint64)

// SvcChangeKernelTraceState Performs a debugging operation on the kernel
//
// This is a privileged syscall.
// Only on firmware 4.0.0+
// svc 0x3D
//export svcChangeKernelTraceState
func SvcChangeKernelTraceState(kernelTraceState uint32)

// SvcCreateSession Creates an IPC session.
//
// This is a privileged syscall.
// svc 0x40
//export svcCreateSession
func SvcCreateSession(serverHandle, clientHandle *uint32, unk0 uint32, unk1 uint64) uint64

// SvcAcceptSession Accepts an IPC session.
//
// This is a privileged syscall.
// svc 0x41
//export svcAcceptSession
func SvcAcceptSession(serverHandle *uint32, portHandle uint32) uint64

// SvcReplyAndReceiveLight Performs light IPC input/output.
//
// This is a privileged syscall.
// svc 0x42
//export svcReplyAndReceiveLight
func SvcReplyAndReceiveLight(handle uint32) uint64

// SvcReplyAndReceive Performs light IPC input/output.
//
// This is a privileged syscall.
// svc 0x43
//export svcReplyAndReceive
func SvcReplyAndReceive(index *int32, handles *uint32, handleCount uint32, replyTarget uint32, timeout uint64) uint64

// SvcReplyAndReceiveWithUserBuffer Performs IPC input/output from an user allocated buffer
//
// This is a privileged syscall.
// svc 0x44
//export svcReplyAndReceiveWithUserBuffer
func SvcReplyAndReceiveWithUserBuffer(index *int32, userBuffer uintptr, size uint64, handles *uint32, handleCount uint32, replyTarget uint32, timeout uint64) uint64

// SvcCreateEvent  Creates a system event.
//
// This is a privileged syscall.
// svc 0x45
//export svcCreateEvent
func SvcCreateEvent(serverHandle, clientHandle *uint32) uint64

// SvcMapPhysicalMemoryUnsafe Maps unsafe memory (usable for GPU DMA) for a system module at the desired address.
//
// Only on firmware 5.0.0+
// This is a privileged syscall.
// svc 0x48
//export svcMapPhysicalMemoryUnsafe
func SvcMapPhysicalMemoryUnsafe(addr, size uintptr) uint64

// SvcMapPhysicalMemoryUnsafe Undoes the effects of SvcMapPhysicalMemoryUnsafe
//
// Only on firmware 5.0.0+
// This is a privileged syscall.
// svc 0x49
//export svcUnmapPhysicalMemoryUnsafe
func SvcUnmapPhysicalMemoryUnsafe(addr, size uintptr) uint64

// SvcSetUnsafeLimit Sets the system-wide limit for unsafe memory mappable using svcMapPhysicalMemoryUnsafe.
//
// Only on firmware 5.0.0+
// This is a privileged syscall.
// svc 0x4A
//export svcSetUnsafeLimit
func SvcSetUnsafeLimit(size uintptr) uint64

// SvcCreateCodeMemory  Creates code memory in the caller's address space
//
// Only on firmware 4.0.0+
// This is a privileged syscall.
// svc 0x4B
//export svcCreateCodeMemory
func SvcCreateCodeMemory(handle *uint32, srcAddr, size uintptr) uint64

// SvcControlCodeMemory Maps code memory in the caller's address space
//
// Only on firmware 4.0.0+
// This is a privileged syscall.
// svc 0x4C
//export svcControlCodeMemory
func SvcControlCodeMemory(codeHandle uint32, op uint64, dstAddr, size uintptr, perm uint64) uint64

// SvcSleepSystem Causes the system to enter deep sleep.
//
// This is a privileged syscall.
// svc 0x4D
//export svcSleepSystem
func SvcSleepSystem()

// SvcReadWriteRegister Reads/writes a protected MMIO register.
//
// This is a privileged syscall.
// svc 0x4E
//export svcReadWriteRegister
func SvcReadWriteRegister(outVal *uint32, regAddr uint64, rwMask, inVal uint32) uint64

// SvcSetProcessActivity Configures the pause/unpause status of a process.
//
// This is a privileged syscall.
// svc 0x4F
//export svcSetProcessActivity
func SvcSetProcessActivity(process uint32, paused bool) uint64

// SvcCreateSharedMemory Creates a block of shared memory.
//
// This is a privileged syscall.
// svc 0x50
//export svcCreateSharedMemory
func SvcCreateSharedMemory(handle *uint32, size uintptr, localPerm, otherPerm uint32) uint64

// SvcMapTransferMemory Maps a block of transfer memory.
//
// This is a privileged syscall.
// svc 0x51
//export svcMapTransferMemory
func SvcMapTransferMemory(tmemHandle uint32, addr, size uintptr, perm uint32) uint64

// SvcUnmapTransferMemory Unmaps a block of transfer memory.
//
// This is a privileged syscall.
// svc 0x52
//export svcUnmapTransferMemory
func SvcUnmapTransferMemory(tmemHandle uint32, addr, size uintptr) uint64

// SvcCreateInterruptEvent Creates an event and binds it to a specific hardware interrupt.
//
// This is a privileged syscall.
// svc 0x53
//export svcCreateInterruptEvent
func SvcCreateInterruptEvent(handle *uint32, irqNum uint64, flag uint32) uint64

// SvcQueryPhysicalAddress Queries information about a certain virtual address, including its physical address.
//
// This is a privileged syscall.
// svc 0x54
//export svcQueryPhysicalAddress
func SvcQueryPhysicalAddress(out uintptr, virtaddr uintptr) uint64
