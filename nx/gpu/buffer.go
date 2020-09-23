package gpu

import (
	"github.com/racerxdl/gonx/nx/nv"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"unsafe"
)

type Buffer struct {
	nvMapHandle uint32
	size        uintptr
	alignment   uint32
	kind        uint8
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

	buff.nvMapHandle = nvIocFromIdArgs.handle

	nvParam := nvmapIocParamArgs{
		handle: buff.nvMapHandle,
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

	buff.size = uintptr(nvParam.value)

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

	buff.alignment = nvParam.value

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

	buff.kind = uint8(nvParam.value)

	return buff, nil
}

func (b *Buffer) GetID() (id uint32, err error) {
	if gpuInitializations <= 0 {
		return 0, nxerrors.GPUNotInitialized
	}

	nvIdArgs := nvmapIocGetIdArgs{
		handle: b.nvMapHandle,
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
		handle: b.nvMapHandle,
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

func CreateBuffer(addr unsafe.Pointer, size uintptr, heapMask uint32, flags uint32, alignment uint32, kind uint8) (*Buffer, error) {
	if gpuInitializations <= 0 {
		return nil, nxerrors.GPUNotInitialized
	}

	uaddr := uintptr(addr)

	if uint64(uaddr)&(uint64(alignment)-1) != 0 {
		// GPU Driver crashes if this is not checked
		return nil, nxerrors.GPUBufferUnaligned
	}

	gpuB := &Buffer{
		size:      size,
		kind:      kind,
		alignment: alignment,
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

	gpuB.nvMapHandle = nvmCreate.handle

	nvmAlloc := nvmapIocAllocArgs{
		handle:   nvmCreate.handle,
		heapmask: heapMask,
		flags:    flags,
		align:    alignment,
		kind:     kind,
		addr:     uint64(uintptr(addr)),
	}

	handle, err = nv.Ioctl(nvmapFd, NVMAP_IOC_ALLOC, unsafe.Pointer(&nvmAlloc), unsafe.Sizeof(nvmAlloc))
	if err != nil {
		return nil, err
	}
	if handle != 0 {
		return nil, nxerrors.IPCError{
			Message: "error calling NVMAP_IOC_ALLOC",
			Result:  uint64(handle),
		}
	}

	return gpuB, nil
}
