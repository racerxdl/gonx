package gpu

import (
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/services/nv"
	"unsafe"
)

const gpuDebug = false

var (
	nvasFd             int32
	nvmapFd            int32
	nvhostCtrlFd       int32
	gpuInitializations = 0
)

type Fence struct {
	SyncptId    uint32
	SyncptValue uint32
}

func (g *Fence) Wait(timeout uint32) error {
	if gpuInitializations <= 0 {
		return nxerrors.GPUNotInitialized
	}

	wait := nvhostIocCtrlSyncPtWaitArgs{
		syncptId:  g.SyncptId,
		threshold: g.SyncptValue,
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
	if gpuDebug {
		println("GPU::Init()")
	}
	gpuInitializations++
	if gpuInitializations > 1 {
		return nil
	}

	nvInit := false
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
			if nvInit {
				nv.Finalize()
			}
			gpuInitializations--
		}
	}()

	if gpuDebug {
		println("GPU::Init() - Init NV")
	}

	err = nv.Init()
	if err != nil {
		return err
	}
	nvInit = true

	if gpuDebug {
		println("GPU::Init() - open nvhost-as-gpu")
	}
	nvasFd, err = nv.Open("/dev/nvhost-as-gpu")
	if err != nil {
		return err
	}
	nvasInit = true

	if gpuDebug {
		println("GPU::Init() - open nvmap")
	}
	nvmapFd, err = nv.Open("/dev/nvmap")
	if err != nil {
		return err
	}
	nvmapInit = true

	if gpuDebug {
		println("GPU::Init() - open nvhost-ctrl")
	}
	nvhostCtrlFd, err = nv.Open("/dev/nvhost-ctrl")

	return err
}

func forceFinalize() {
	if gpuDebug {
		println("GPU::ForceFinalize()")
	}
	_ = nv.Close(nvhostCtrlFd)
	_ = nv.Close(nvmapFd)
	_ = nv.Close(nvasFd)
	nv.Finalize()
	gpuInitializations = 0
}

func Finalize() {
	if gpuDebug {
		println("GPU::Finalize()")
	}
	gpuInitializations--
	if gpuInitializations <= 0 {
		forceFinalize()
	}
}
