package display

import (
	"github.com/racerxdl/gonx/nx/gpu"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/vi"
	"unsafe"
)

type Fence struct {
	IsValid uint32
	Sync    [4]gpu.Fence
}

func UnflattenFence(p *vi.Parcel) (fence Fence, err error) {
	size := p.ReadU32()
	numFds := p.ReadU32()

	if uintptr(size) != unsafe.Sizeof(fence) {
		return fence, nxerrors.DisplayInvalidFence
	}

	if numFds != 0 {
		return fence, nxerrors.DisplayFenceTooManyFds
	}

	fence.IsValid = p.ReadU32()

	for i := range fence.Sync {
		fence.Sync[i].SyncptId = 0xFFFFFFFF // Fill with default values
	}

	for i := 0; i < len(fence.Sync) && p.Remaining() > 8; i++ {
		fence.Sync[i].SyncptId = p.ReadU32()
		fence.Sync[i].SyncptValue = p.ReadU32()
	}

	return fence, nil
}
