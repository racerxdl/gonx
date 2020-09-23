package gpu

import (
	"github.com/racerxdl/gonx/nx/nv"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"unsafe"
)

var (
	nvasFd             int32
	nvmapFd            int32
	nvhostCtrlFd       int32
	gpuInitializations = 0
)

type gpuFence struct {
	syncptId    uint32
	syncptValue uint32
}

func (g *gpuFence) Wait(timeout uint32) error {
	if gpuInitializations <= 0 {
		return nxerrors.GPUNotInitialized
	}

	wait := nvhostIocCtrlSyncPtWaitArgs{
		syncptId:  g.syncptId,
		threshold: g.syncptValue,
		timeout:   timeout,
	}

	handle, err := nv.Ioctl(nvmapFd, NVHOST_IOC_CTRL_SYNCPT_WAIT, unsafe.Pointer(&wait), unsafe.Sizeof(wait))
	if err != nil {
		return err
	}

	if handle != 0 {
		return nxerrors.IPCError{
			Message: "error calling NVHOST_IOC_CTRL_SYNCPT_WAIT",
			Result:  uint64(handle),
		}
	}

	return nil
}

func Init() (err error) {
	gpuInitializations++
	if gpuInitializations > 1 {
		return nil
	}

	nvmapInit := false
	nvasInit := false

	defer func() {
		if err != nil {
			if nvmapInit {
				_ = nv.Close(nvmapFd)
			}
			if nvasInit {
				_ = nv.Close(nvasFd)
			}

			nv.Finalize()
			gpuInitializations--
		}
	}()

	nvasFd, err = nv.Open("/dev/nvhost-as-gpu")
	if err != nil {
		return err
	}
	nvasInit = true

	nvmapFd, err = nv.Open("/dev/nvmap")
	if err != nil {
		return err
	}
	nvmapInit = true

	nvhostCtrlFd, err = nv.Open("/dev/nvhost-ctrl")

	return err
}

func forceFinalize() {
	_ = nv.Close(nvhostCtrlFd)
	_ = nv.Close(nvmapFd)
	_ = nv.Close(nvasFd)
	nv.Finalize()
	gpuInitializations = 0
}

func Finalize() {
	gpuInitializations--
	if gpuInitializations <= 0 {
		forceFinalize()
	}
}
