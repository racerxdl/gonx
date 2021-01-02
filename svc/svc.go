package svc

//
//
//// SvcSetMemoryPermission Sets memory permission
//// svc 0x02
////export svcSetMemoryPermission
//func SvcSetMemoryPermission(addr uintptr, size uintptr, perm uint32) uint64
//
//// SvcSetMemoryAttribute Sets memory attributes
//// svc 0x03
////export svcSetMemoryAttribute
//func SvcSetMemoryAttribute(addr uintptr, size uintptr, state0, state1 uint32) uint64
//
//// SvcMapMemory Maps a memory range into a different range. Mainly used for adding guard pages around stack.
//// svc 0x04
////export svcMapMemory
//func SvcMapMemory(dstAddr, srcAddr uintptr, size uintptr) uint64
//
//// SvcUnmapMemory Unmaps a region that was previously mapped with SvcMapMemory.
//// svc 0x05
////export svcUnmapMemory
//func SvcUnmapMemory(dstAddr, srcAddr uintptr, size uintptr) uint64
//
//// SvcQueryMemory Query information about an address. Will always fetch the lowest page-aligned mapping that contains the provided address.
//// svc 0x06
////export svcQueryMemory
//func SvcQueryMemory(memoryInfo, pageInfo uintptr, size uintptr) uint64
//
//
//// SvcCreateThread Creates a thread.
//// svc 0x08
////export svcCreateThread
//func SvcCreateThread(handle *uint32, entry, arg, stackTop uintptr, prio, cpuid int) uint64
//
//// SvcStartThread Starts a freshly created thread.
//// svc 0x09
////export svcStartThread
//func SvcStartThread(handle uint32) uint64
//
//// SvcExitThread Exits the current thread.
//// svc 0x0A
////export svcExitThread
//func SvcExitThread(handle uint32) uint64
//
//// SvcExitThread Gets a thread's priority.
//// svc 0x0C
////export svcGetThreadPriority
//func svcGetThreadPriority(priority, handle *uint32) uint64
//
//// SvcSetThreadPriority Gets a thread's priority.
//// svc 0x0D
////export svcSetThreadPriority
//func SvcSetThreadPriority(handle uint32, priority *uint32) uint64
//
//// SvcGetThreadCoreMask Gets a thread's core mask.
//// svc 0x0E
////export svcGetThreadCoreMask
//func SvcGetThreadCoreMask(preferedCore *uint32, affinityMask *uint64, handle uint32) uint64
//
//// SvcSetThreadCoreMask Sets a thread's core mask.
//// svc 0x0F
////export svcSetThreadCoreMask
//func SvcSetThreadCoreMask(handle uint32, preferedCore uint32, affinityMask uint64) uint64
//
//// SvcGetCurrentProcessorNumber Gets the current processor's number.
//// svc 0x10
////export svcGetCurrentProcessorNumber
//func SvcGetCurrentProcessorNumber() uint32
//
//// SvcSignalEvent Sets an event's signalled status.
//// svc 0x11
////export svcSignalEvent
//func SvcSignalEvent(handle uint32) uint64
//
//// SvcClearEvent Clears an event's signalled status.
//// svc 0x12
////export svcClearEvent
//func SvcClearEvent(handle uint32) uint64
//
//// SvcMapSharedMemory Maps a block of shared memory.
//// svc 0x13
////export svcMapSharedMemory
//func SvcMapSharedMemory(handle uint32, addr uintptr, size uintptr, perm uint32) uint64
//
//// SvcUnmapSharedMemory Unmaps a block of shared memory.
//// svc 0x14
////export svcUnmapSharedMemory
//func SvcUnmapSharedMemory(handle uint32, addr uintptr, size uintptr) uint64
//
//// SvcCreateTransferMemory Creates a block of transfer memory.
//// svc 0x15
////export svcCreateTransferMemory
//func SvcCreateTransferMemory(handle *uint32, addr uintptr, size uintptr, perm uint32) uint64
//
//// SvcCloseHandle Closes a handle, decrementing the reference count of the corresponding kernel object.
//// This might result in the kernel freeing the object.
//// svc 0x16.
////export svcCloseHandle
//func SvcCloseHandle(session uint32) uint64
//
//// SvcResetSignal Resets a signal.
//// svc 0x17.
////export svcResetSignal
//func SvcResetSignal(handle uint32) uint64
//
//// SvcWaitSynchronization Waits on one or more synchronization objects, optionally with a timeout.
//// handleCount must not be greater than 40. This is a Horizon Kernel Limitation
//// svc 0x18
////export svcWaitSynchronization
//func SvcWaitSynchronization(index *uint32, handles *uint32, handleCount int32, timeout uint64) uint64
//
//// SvcCancelSynchronization Cancels a svcWaitSynchronization operation being done on a synchronization object in another thread.
//// svc 0x19
////export svcCancelSynchronization
//func SvcCancelSynchronization(thread uint32) uint64
//
//// SvcArbitrateLock Arbitrates a mutex lock operation in userspace.
//// svc 0x1A
////export svcArbitrateLock
//func SvcArbitrateLock(waitTag uint32, tagLocation *uint32, selfTag uint32) uint64
//
//// SvcArbitrateUnlock Arbitrates a mutex unlock operation in userspace.
//// svc 0x1B
////export svcArbitrateUnlock
//func SvcArbitrateUnlock(tagLocation *uint32) uint64
//
//// SvcWaitProcessWideKeyAtomic Performs a condition variable wait operation in userspace.
//// svc 0x1C
////export svcWaitProcessWideKeyAtomic
//func SvcWaitProcessWideKeyAtomic(key, tagLocation *uint32, selfTag uint32, timeout uint64) uint64
//
//// SvcSignalProcessWideKey Performs a condition variable wake-up operation in userspace.
//// svc 0x1D
////export svcSignalProcessWideKey
//func SvcSignalProcessWideKey(key *uint32, num uint32)
//
//// SvcGetSystemTick Gets the current system tick.
//// svc 0x1E
////export svcGetSystemTick
//func SvcGetSystemTick() uint64
//
//// SvcConnectToNamedPort Connects to a registered named port.
//// Expects byte to be a null terminated string
//// svc 0x1F
////export svcConnectToNamedPort
//func SvcConnectToNamedPort(session *uint32, name *byte) uint64
//
//// SvcSendSyncRequestLight Sends a light IPC synchronization request to a session.
//// svc 0x20
////export svcSendSyncRequestLight
//func SvcSendSyncRequestLight(session uint64) uint64
//
//// SvcSendSyncRequest Sends an IPC synchronization request to a session.
//// svc 0x21
////export svcSyncRequest
//func SvcSendSyncRequest(session uint64) uint64
//
//// SvcSendSyncRequestWithUserBuffer Sends an IPC synchronization request to a session from an user allocated buffer.
//// Size must be page-aligned (0x1000)
//// svc 0x22
////export svcSendSyncRequestWithUserBuffer
//func SvcSendSyncRequestWithUserBuffer(userBuffer uintptr, size uintptr, session uint64) uint64
//
//// SvcSendAsyncRequestWithUserBuffer Sends an IPC synchronization request to a session from an user allocated buffer (asynchronous version).
//// Size must be page-aligned (0x1000)
//// svc 0x23
////export svcSendAsyncRequestWithUserBuffer
//func SvcSendAsyncRequestWithUserBuffer(handle *uint32, userBuffer uintptr, size uintptr, session uint64) uint64
//
//// SvcGetProcessId Gets the PID associated with a process.
//// svc 0x24
////export svcGetProcessId
//func SvcGetProcessId(processID *uint64, handle uint32) uint64
//
//// SvcGetThreadId Gets the TID associated with a process.
//// svc 0x25
////export svcGetThreadId
//func SvcGetThreadId(threadID *uint64, handle uint32) uint64
//
//// SvcBreak  Breaks execution. Panic.
//// Used for triggering a debugger
//// svc 0x26
////export svcBreak
//func SvcBreak(breakReason uint32, inval1, inval2 uint64) uint64
//
//
//// SvcReturnFromException Returns from an exception.
//// NO RETURN
//// svc 0x28
////export svcReturnFromException
//func SvcReturnFromException(result uint64)

//// SvcFlushEntireDataCache Flushes the entire data cache (by set/way).
//// This is a privileged syscall and is dangerous, and should not be used if not needed
//// svc 0x2A
////export svcFlushEntireDataCache
//func SvcFlushEntireDataCache() uint64
//
//// SvcFlushDataCache Flushes data cache for a virtual address range.
//// svc 0x2B
////export svcFlushDataCache
//func SvcFlushDataCache(addr, size uintptr) uint64
//
//// SvcMapPhysicalMemory Maps new heap memory at the desired address.
//// Only on firmware 3.0.0+
//// svc 0x2C
////export svcMapPhysicalMemory
//func SvcMapPhysicalMemory(addr, size uintptr) uint64
//
//// SvcUnmapPhysicalMemory Undoes the effects of SvcMapPhysicalMemory.
//// Only on firmware 3.0.0+
//// svc 0x2D
////export svcUnmapPhysicalMemory
//func SvcUnmapPhysicalMemory(addr, size uintptr) uint64
//
//// SvcGetDebugFutureThreadInfo Gets information about a thread that will be scheduled in the future.
//// outContext  -> Output LastThreadContext for the thread that will be scheduled.
//// outThreadId -> Output thread id for the thread that will be scheduled.
//// debug       -> handle.
//// nanoseconds -> Nanoseconds in the future to get scheduled thread at.
////
//// This is a privileged syscall
//// Only on firmware 5.0.0+
//// svc 0x2E
////export svcGetDebugFutureThreadInfo
//func SvcGetDebugFutureThreadInfo(outContext uintptr, outThreadId *uint64, debug uint32, nanoseconds uint64) uint64
//
//// SvcGetLastThreadInfo Gets information about the previously-scheduled thread.
//// outContext    -> Output LastThreadContext for the thread that will be scheduled.
//// outTLSAddress -> Output tls address for the previously scheduled thread.
//// outFlags      -> Output flags for the previously scheduled thread.
////
//// svc 0x2F
////export svcGetLastThreadInfo
//func SvcGetLastThreadInfo(outContext uintptr, outTLSAddress *uint64, outFlags *uint32) uint64
//
//// SvcGetResourceLimitLimitValue Gets the maximum value a LimitableResource can have, for a Resource Limit handle.
////
//// This is a privileged syscall.
//// svc 0x30
////export svcGetResourceLimitLimitValue
//func SvcGetResourceLimitLimitValue(out *int64, resLimit uint32, which uint32) uint64
//
//// SvcGetResourceLimitCurrentValue Gets the maximum value a LimitableResource can have, for a Resource Limit handle.
////
//// This is a privileged syscall.
//// svc 0x31
////export svcGetResourceLimitCurrentValue
//func SvcGetResourceLimitCurrentValue(out *int64, resLimit uint32, which uint32) uint64
//
//// SvcSetThreadActivity Configures the pause/unpause status of a thread.
////
//// svc 0x32
////export svcSetThreadActivity
//func SvcSetThreadActivity(thread uint32, paused bool) uint64
//
//// SvcGetThreadContext3 Dumps the registers of a thread paused by svcSetThreadActivity (register groups: all).
////
//// svc 0x33
////export svcGetThreadContext3
//func SvcGetThreadContext3(ctx uintptr, thread uint32) uint64
//
//// SvcWaitForAddress Arbitrates an address depending on type and value.
////
//// Only on firmware 4.0.0+
//// svc 0x34
////export svcWaitForAddress
//func SvcWaitForAddress(addr uintptr, arbType uint32, value uint32, timeout uint64) uint64
//
//// SvcSignalToAddress Signals (and updates) an address depending on type and value.
////
//// Only on firmware 4.0.0+
//// svc 0x35
////export svcSignalToAddress
//func SvcSignalToAddress(addr uintptr, signalType, value, count uint32) uint64
//
//// SvcSynchronizePreemptionState Sets thread preemption state (used during abort/panic).
////
//// Only on firmware 8.0.0+
//// svc 0x36
////export svcSynchronizePreemptionState
//func SvcSynchronizePreemptionState()
