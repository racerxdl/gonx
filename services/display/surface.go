package display

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/memory"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/gpu"
	"github.com/racerxdl/gonx/services/vi"
	"github.com/racerxdl/gonx/svc"
	"image"
	"unsafe"
)

// SurfaceState Keeps track of the internal state of a \ref surface_t
type SurfaceState int

const (
	SURFACE_STATE_INVALID  SurfaceState = iota
	SURFACE_STATE_DEQUEUED SurfaceState = iota
	SURFACE_STATE_QUEUED   SurfaceState = iota
)

type Surface struct {
	LayerId      uint64
	IGBP         vi.IGBP
	State        SurfaceState
	HasRequested [3]bool
	CurrentSlot  uint32

	GpuBuffer            *gpu.Buffer
	GpuBufferMemory      []byte
	GpuBufferMemoryAlloc uintptr
	GraphicBuffers       [3]GraphicBuffer
	CurrentFence         Fence
}

func (s *Surface) Destroy() {
	if s.State == SURFACE_STATE_INVALID {
		return
	}

	_, _ = IGBPDisconnect(s.IGBP, 2, DisconnectAllLocal)
	_ = vi.AdjustRefCount(s.IGBP.IgbpBinder.Handle, -1, 1)
	_ = vi.CloseLayer(s.LayerId)
	_ = vi.DestroyManagedLayer(s.LayerId)
	s.State = SURFACE_STATE_INVALID

	_, _, _ = s.GpuBuffer.Destroy()
	svc.SetMemoryAttribute(uintptr(unsafe.Pointer(&s.GpuBufferMemory[0])), uintptr(0x3c0000*len(s.GraphicBuffers)), 0, 0)
	s.GpuBufferMemory = nil
}

func (s *Surface) CloseLayer() error {
	if displayInitializations < 1 {
		return nxerrors.DisplayNotInitialized
	}

	s.Destroy()
	_ = vi.AdjustRefCount(s.IGBP.IgbpBinder.Handle, -1, 1)
	_ = vi.CloseLayer(s.LayerId)
	_ = vi.DestroyManagedLayer(s.LayerId)

	return nil
}

func (s *Surface) DequeueBuffer() ([]byte, error) {
	if s.State != SURFACE_STATE_QUEUED {
		return nil, nxerrors.SurfaceInvalidState
	}

	bf := s.GraphicBuffers[s.CurrentSlot]

	status, slot, fence, _, err := IGBPDequeueBuffer(s.IGBP, bf.Width, bf.Height, bf.Format, bf.Usage, false)
	if err != nil {
		return nil, err
	}

	s.CurrentSlot = slot
	s.CurrentFence = fence

	if status != 0 {
		return nil, nxerrors.SurfaceBufferDequeueFailed
	}

	if !s.HasRequested[s.CurrentSlot] {
		_, _, err = IGBPRequestBuffer(s.IGBP, s.CurrentSlot)
		if err != nil {
			return nil, err
		}
		s.HasRequested[s.CurrentSlot] = true
	}

	imageSlice := s.GpuBufferMemory[(s.CurrentSlot * 0x3c0000):]
	s.State = SURFACE_STATE_DEQUEUED

	return imageSlice, nil
}

func (s *Surface) QueueBuffer() error {
	if s.State != SURFACE_STATE_DEQUEUED {
		return nxerrors.SurfaceInvalidState
	}

	qbi := &QueueBufferInput{
		Size:    uint32(unsafe.Sizeof(QueueBufferInput{}) - 8),
		Unknown: [2]uint32{0, 1},
		Fence: Fence{
			IsValid: 1,
			Sync: [4]gpu.Fence{
				{SyncptId: 0xffffffff},
				{SyncptId: 0xffffffff},
				{SyncptId: 0xffffffff},
				{SyncptId: 0xffffffff},
			},
		},
	}

	_, status, err := IGBPQueueBuffer(s.IGBP, int(s.CurrentSlot), qbi)
	if err != nil {
		return err
	}

	if status != 0 {
		return nxerrors.SurfaceBufferQueueFailed
	}

	s.State = SURFACE_STATE_QUEUED

	return nil
}

func SurfaceCreate(layerId uint64, igbp vi.IGBP) (surface *Surface, status int, err error) {
	if debugDisplay {
		fmt.Printf("Display::SurfaceCreate(%d, %d)\n", layerId, igbp.IgbpBinder.Handle)
	}
	surface = &Surface{
		LayerId:        layerId,
		IGBP:           igbp,
		State:          SURFACE_STATE_INVALID,
		CurrentSlot:    0,
		GraphicBuffers: [3]GraphicBuffer{},
	}

	memoryAttributesSet := false
	igbpConnected := false
	numBuffers := len(surface.GraphicBuffers)
	bufferSize := numBuffers * 0x3c0000

	defer func() {
		if err != nil {
			if memoryAttributesSet {
				svc.SetMemoryAttribute(uintptr(unsafe.Pointer(&surface.GpuBufferMemory[0])), uintptr(bufferSize), 0, 0)
			}
			if igbpConnected {
				_, _ = IGBPDisconnect(surface.IGBP, 2, DisconnectAllLocal)
			}
		}
	}()

	var qbo *QueueBufferOutput

	status, qbo, err = IGBPConnect(igbp, 2, false)
	if err != nil {
		return nil, status, err
	}

	surface.GpuBufferMemory = memory.AllocPages(bufferSize, bufferSize)
	if surface.GpuBufferMemory == nil {
		return nil, status, nxerrors.OutOfMemory
	}

	r := svc.SetMemoryAttribute(uintptr(unsafe.Pointer(&surface.GpuBufferMemory[0])), uintptr(bufferSize), 0x8, 0x8)
	if r != nxtypes.ResultOK {
		return nil, status, nxerrors.CannotSetMemoryAttributes
	}

	surface.GpuBuffer, err = gpu.CreateBuffer(unsafe.Pointer(&surface.GpuBufferMemory[0]), uintptr(bufferSize), 0, 0, 0x1000, 0)
	if err != nil {
		return nil, status, err
	}

	for i := range surface.GraphicBuffers {
		surface.GraphicBuffers[i] = GraphicBuffer{
			Width:     qbo.Width,
			Height:    qbo.Height,
			Stride:    qbo.Width,
			Format:    RGBA_8888,
			Usage:     GRALLOC_USAGE_HW_COMPOSER | GRALLOC_USAGE_HW_RENDER | GRALLOC_USAGE_HW_TEXTURE,
			GPUBuffer: surface.GpuBuffer,
		}
		surface.GraphicBuffers[i].PixelBufferOffset = uint32(0x3c0000 * i)

		err = IGBPSetPreallocatedBuffer(surface.IGBP, i, &surface.GraphicBuffers[i])
		if err != nil {
			return nil, status, err
		}
	}

	surface.State = SURFACE_STATE_QUEUED
	return surface, status, nil
}

func (s *Surface) refreshFrame(f *Frame) error {
	data, err := s.DequeueBuffer()
	if err != nil {
		return err
	}
	slot := s.CurrentSlot

	f.surfaceBuff = data
	f.bounds = image.Rect(0, 0, int(s.GraphicBuffers[slot].Width), int(s.GraphicBuffers[slot].Height))

	b := f.bounds.Size()
	l := b.X * b.Y * 4 // uint32 per pixel

	if len(f.buff) != l {
		if debugDisplay {
			fmt.Printf("Allocating buffer of %d bytes\n", l)
		}
		f.buff = make([]byte, l)
	}

	return nil
}

// GetFrame returns a frame to be draw on screen
func (s *Surface) GetFrame() (*Frame, error) {
	f := &Frame{
		surface:     s,
		buff:        nil,
		surfaceBuff: nil,
	}

	err := s.refreshFrame(f)
	if err != nil {
		return nil, err
	}

	return f, nil
}
