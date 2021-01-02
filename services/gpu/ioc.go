package gpu

const (
	NVHOST_IOC_CTRL_SYNCPT_READ      = 0xC0080014
	NVHOST_IOC_CTRL_SYNCPT_INCR      = 0x40040015
	NVHOST_IOC_CTRL_SYNCPT_WAIT      = 0xC00C0016
	NVHOST_IOC_CTRL_MODULE_MUTEX     = 0x40080017
	NVHOST_IOC_CTRL_MODULE_REGRDWR   = 0xC0180018
	NVHOST_IOC_CTRL_SYNCPT_WAITEX    = 0xC0100019
	NVHOST_IOC_CTRL_SYNCPT_READ_MAX  = 0xC008001A
	NVHOST_IOC_CTRL_GET_CONFIG       = 0xC183001B
	NVHOST_IOC_CTRL_EVENT_SIGNAL     = 0xC004001C
	NVHOST_IOC_CTRL_EVENT_WAIT       = 0xC010001D
	NVHOST_IOC_CTRL_EVENT_WAIT_ASYNC = 0xC010001E
	NVHOST_IOC_CTRL_EVENT_REGISTER   = 0xC004001F
	NVHOST_IOC_CTRL_EVENT_UNREGISTER = 0xC0040020
	NVHOST_IOC_CTRL_EVENT_KILL       = 0x40080021

	NVMAP_IOC_CREATE  = 0xC0080101
	NVMAP_IOC_FROM_ID = 0xC0080103
	NVMAP_IOC_ALLOC   = 0xC0200104
	NVMAP_IOC_FREE    = 0xC0180105
	NVMAP_IOC_PARAM   = 0xC00C0109
	NVMAP_IOC_GET_ID  = 0xC008010E
)

// nvhostIocCtrlSyncPtWaitArgs Arguments to wait on a syncpt
type nvhostIocCtrlSyncPtWaitArgs struct {
	syncptId  uint32 // In
	threshold uint32 // In
	timeout   uint32 // In
}

// nvhostIocCtrlEventWaitArgs Arguments to wait on a syncpt event
type nvhostIocCtrlEventWaitArgs struct {
	syncptId  uint32 // In
	threshold uint32 // In
	timeout   int32  // In
	value     uint32 // Inout
}

// nvmapIocCreateArgs Args to create an nvmap object
// Identical to Linux Driver
type nvmapIocCreateArgs struct {
	size   uint32 // In
	handle uint32 // Out
}

// nvmapIocFromIdArgs Args to get the handle to an existing nvmap object
// Identical to Linux Driver
type nvmapIocFromIdArgs struct {
	id     uint32 // In
	handle uint32 // Out
}

// nvmapIocAllocArgs Memory allocation args structure for the nvmap object.
// Nintendo extended this one with 16 bytes, and changed it from in to inout.
type nvmapIocAllocArgs struct {
	handle   uint32
	heapmask uint32
	flags    uint32 // 0 = readonly, 1 = readwrite
	align    uint32
	kind     uint8
	pad      [7]uint8
	addr     uint64
}

// nvmapIocFreeArgs Memory freeing args structure for the nvmap object.
type nvmapIocFreeArgs struct {
	handle   uint32
	pad      uint32
	refcount uint64 // out
	size     uint32 // out
	flags    uint32 // out ( 1 = not freed yet )
}

// nvmapIocParamArgs Info query args structure for an nvmap object.
// Identical to Linux driver, but extended with further params.
type nvmapIocParamArgs struct {
	handle uint32
	param  uint32 // 1=SIZE, 2=ALIGNMENT, 3=BASE (returns error), 4=HEAP (always 0x40000000), 5=KIND, 6=COMPR (unused)
	value  uint32
}

// nvmapIocGetIdArgs ID query args structure for an nvmap object.
// Identical to Linux driver.
type nvmapIocGetIdArgs struct {
	id     uint32 // Out ~0 indicates error
	handle uint32 // In
}
