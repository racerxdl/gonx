package gpu

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/nv"
	"github.com/racerxdl/gonx/svc"
	"unsafe"
)

type Buffer struct {
	NvMapHandle uint32
	Size        uintptr
	Alignment   uint32
	Kind        uint8
}

func InitializeFromId(id uint32) (*Buffer, error) {
	if gpuInitializations <= 0 {
		return nil, nxerrors.GPUNotInitialized
	}

	buff := &Buffer{}

	nvIocFromIdArgs := nvmapIocGetIdArgs{
		id: id,
	}

	handle, err := nv.Ioctl(nvmapFd, NVMAP_IOC_FROM_ID, unsafe.Pointer(&nvIocFromIdArgs), unsafe.Sizeof(nvIocFromIdArgs))
	if err != nil {
		return nil, err
	}

	if handle != 0 {
		return nil, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_FROM_ID",
			Result:  uint64(handle),
		}
	}

	buff.NvMapHandle = nvIocFromIdArgs.handle

	nvParam := nvmapIocParamArgs{
		handle: buff.NvMapHandle,
		param:  1, // SIZE
	}

	handle, err = nv.Ioctl(nvmapFd, NVMAP_IOC_PARAM, unsafe.Pointer(&nvParam), unsafe.Sizeof(nvParam))
	if err != nil {
		return nil, err
	}

	if handle != 0 {
		return nil, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_PARAM",
			Result:  uint64(handle),
		}
	}

	buff.Size = uintptr(nvParam.value)

	nvParam.param = 2 // ALIGNMENT
	handle, err = nv.Ioctl(nvmapFd, NVMAP_IOC_PARAM, unsafe.Pointer(&nvParam), unsafe.Sizeof(nvParam))
	if err != nil {
		return nil, err
	}

	if handle != 0 {
		return nil, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_PARAM",
			Result:  uint64(handle),
		}
	}

	buff.Alignment = nvParam.value

	nvParam.param = 5 // KIND
	handle, err = nv.Ioctl(nvmapFd, NVMAP_IOC_PARAM, unsafe.Pointer(&nvParam), unsafe.Sizeof(nvParam))
	if err != nil {
		return nil, err
	}

	if handle != 0 {
		return nil, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_PARAM",
			Result:  uint64(handle),
		}
	}

	buff.Kind = uint8(nvParam.value)

	return buff, nil
}

func (b *Buffer) GetID() (id uint32, err error) {
	if gpuInitializations <= 0 {
		return 0, nxerrors.GPUNotInitialized
	}

	nvIdArgs := nvmapIocGetIdArgs{
		handle: b.NvMapHandle,
	}

	handle, err := nv.Ioctl(nvmapFd, NVMAP_IOC_GET_ID, unsafe.Pointer(&nvIdArgs), unsafe.Sizeof(nvIdArgs))

	if err != nil {
		return 0, err
	}

	if handle != 0 {
		return 0, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_GET_ID",
			Result:  uint64(handle),
		}
	}
	return nvIdArgs.id, nil
}

func (b *Buffer) Destroy() (refCount uint32, flags uint32, err error) {
	if gpuInitializations <= 0 {
		return 0, 0, nxerrors.GPUNotInitialized
	}

	nvmFree := nvmapIocFreeArgs{
		handle: b.NvMapHandle,
	}

	handle, err := nv.Ioctl(nvmapFd, NVMAP_IOC_FREE, unsafe.Pointer(&nvmFree), unsafe.Sizeof(nvmFree))

	if err != nil {
		return 0, 0, err
	}

	if handle != 0 {
		return 0, 0, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_FREE",
			Result:  uint64(handle),
		}
	}

	return uint32(nvmFree.refcount), nvmFree.flags, nil
}

// CreateBuffer creates a buffer in GPU
// equivalent to nvMapCreate on libnx
func CreateBuffer(addr unsafe.Pointer, size uintptr, heapMask, alignment uint32, kind nv.Kind) (*Buffer, error) {
	if gpuInitializations <= 0 {
		return nil, nxerrors.GPUNotInitialized
	}

	if alignment < 0x1000 {
		alignment = 0x1000
	}

	uaddr := uintptr(addr)

	if uint64(uaddr)&(uint64(alignment)-1) != 0 {
		// GPU Driver crashes if this is not checked
		fmt.Println("A")
		return nil, nxerrors.GPUBufferUnaligned
	}

	if size == 0 || (size&0xFFF > 0) {
		// GPU Driver crashes if this is not checked
		fmt.Println("B")
		return nil, nxerrors.GPUBufferUnaligned
	}

	if addr == nil || (uintptr(addr)&0xFFF > 0) {
		// GPU Driver crashes if this is not checked
		fmt.Println("C")
		return nil, nxerrors.GPUBufferUnaligned
	}

	gpuB := &Buffer{
		Size:      size,
		Kind:      uint8(kind),
		Alignment: alignment,
	}

	nvmCreate := nvmapIocCreateArgs{
		size: uint32(size),
	}

	handle, err := nv.Ioctl(nvmapFd, NVMAP_IOC_CREATE, unsafe.Pointer(&nvmCreate), unsafe.Sizeof(nvmCreate))
	if err != nil {
		return nil, err
	}
	if handle != 0 {
		return nil, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_CREATE",
			Result:  uint64(handle),
		}
	}

	gpuB.NvMapHandle = nvmCreate.handle

	err = nvmapAlloc(nvmCreate.handle, heapMask, 0, alignment, uint32(kind), uintptr(addr))
	if err != nil {
		return nil, err
	}

	r := svc.SetMemoryAttribute(uintptr(addr), uintptr(size), 0x8, 0x8)
	if r != nxtypes.ResultOK {
		return nil, nxerrors.CannotSetMemoryAttributes
	}

	return gpuB, nil
}

func nvmapAlloc(nvmapHandle, heapMask, flags, align, kind uint32, addr uintptr) error {
	nvmAlloc := nvmapIocAllocArgs{
		handle:   nvmapHandle,
		heapmask: heapMask,
		flags:    flags,
		align:    align,
		kind:     uint8(kind),
		addr:     uint64(addr),
	}

	handle, err := nv.Ioctl(nvmapFd, NVMAP_IOC_ALLOC, unsafe.Pointer(&nvmAlloc), unsafe.Sizeof(nvmAlloc))
	if err != nil {
		return err
	}
	if handle != 0 {
		return nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_ALLOC",
			Result:  uint64(handle),
		}
	}

	return nil
}
