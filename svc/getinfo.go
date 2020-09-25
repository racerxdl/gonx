package svc

// ID0 for SvcGetInfo
const (
	InfoType_CoreMask                    = 0          // Bitmask of allowed Core IDs.
	InfoType_PriorityMask                = 1          // Bitmask of allowed Thread Priorities.
	InfoType_AliasRegionAddress          = 2          // Base of the Alias memory region.
	InfoType_AliasRegionSize             = 3          // Size of the Alias memory region.
	InfoType_HeapRegionAddress           = 4          // Base of the Heap memory region.
	InfoType_HeapRegionSize              = 5          // Size of the Heap memory region.
	InfoType_TotalMemorySize             = 6          // Total amount of memory available for process.
	InfoType_UsedMemorySize              = 7          // Amount of memory currently used by process.
	InfoType_DebuggerAttached            = 8          // Whether current process is being debugged.
	InfoType_ResourceLimit               = 9          // Current process's resource limit handle.
	InfoType_IdleTickCount               = 10         // Number of idle ticks on CPU.
	InfoType_RandomEntropy               = 11         // [2.0.0+] Random entropy for current process.
	InfoType_AslrRegionAddress           = 12         // [2.0.0+] Base of the process's address space.
	InfoType_AslrRegionSize              = 13         // [2.0.0+] Size of the process's address space.
	InfoType_StackRegionAddress          = 14         // [2.0.0+] Base of the Stack memory region.
	InfoType_StackRegionSize             = 15         // [2.0.0+] Size of the Stack memory region.
	InfoType_SystemResourceSizeTotal     = 16         // [3.0.0+] Total memory allocated for process memory management.
	InfoType_SystemResourceSizeUsed      = 17         // [3.0.0+] Amount of memory currently used by process memory management.
	InfoType_ProgramId                   = 18         // [3.0.0+] Program ID for the process.
	InfoType_InitialProcessIdRange       = 19         // [4.0.0-4.1.0] Min/max initial process IDs.
	InfoType_UserExceptionContextAddress = 20         // [5.0.0+] Address of the process's exception context (for break).
	InfoType_TotalNonSystemMemorySize    = 21         // [6.0.0+] Total amount of memory available for process, excluding that for process memory management.
	InfoType_UsedNonSystemMemorySize     = 22         // [6.0.0+] Amount of memory used by process, excluding that for process memory management.
	InfoType_IsApplication               = 23         // [9.0.0+] Whether the specified process is an Application.
	InfoType_ThreadTickCount             = 0xF0000002 // Number of ticks spent on thread.
)
