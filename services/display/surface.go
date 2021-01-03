package display

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/graphics"
	"github.com/racerxdl/gonx/nx/memory"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/am"
	"github.com/racerxdl/gonx/services/gpu"
	"github.com/racerxdl/gonx/services/nv"
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
	LayerId         uint64
	IGBP            vi.IGBP
	IGBPIsConnected bool
	State           SurfaceState
	HasRequested    []bool
	CurrentSlot     uint32

	GpuBuffer            *gpu.Buffer
	GpuBufferMemory      []byte
	GpuBufferMemoryAlloc uintptr
	GraphicBuffers       []GraphicBuffer
	CurrentFence         Fence
	GRBuff               *nv.GraphicBuffer

	status int
}

func (s *Surface) Connect() (int, error) {
	if s.IGBPIsConnected {
		return s.status, nil
	}

	status, _, err := IGBPConnect(s.IGBP, NativeWindowAPICPU, false)
	s.status = status
	if err != nil {
		return status, err
	}
	s.IGBPIsConnected = true

	return status, err
}

func (s *Surface) Disconnect() {
	if s.IGBPIsConnected {
		_, _ = IGBPDisconnect(s.IGBP, 2, DisconnectAllLocal)
	}
	s.IGBPIsConnected = false
}

func (s *Surface) Destroy() {
	if s.State == SURFACE_STATE_INVALID {
		return
	}

	_, _ = s.DequeueBuffer()

	s.Disconnect()
	_ = vi.AdjustRefCount(s.IGBP.IgbpBinder.Handle, -1, 1)
	_ = vi.CloseLayer(s.LayerId)

	aruid, err := am.IwcGetAppletResourceUserId()
	if err != nil {
		return
	}
	if aruid == 0 {
		_ = vi.DestroyManagedLayer(s.LayerId)
	}
	s.State = SURFACE_STATE_INVALID

	svc.SetMemoryAttribute(uintptr(unsafe.Pointer(&s.GpuBufferMemory[0])), uintptr(0x3c0000*len(s.GraphicBuffers)), 0, 0)
	_, _, _ = s.GpuBuffer.Destroy()
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

	imageSlice := s.GpuBufferMemory[(s.CurrentSlot * s.GRBuff.TotalSize):]
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

	defer func() {
		if err != nil && surface != nil {
			surface.Disconnect()
		}
	}()

	format := graphics.PixelFormatRgba8888
	width := uint32(1280)
	height := uint32(720)
	numFbs := 2

	nvColorFmt := nv.ColorFormatTable[format-graphics.PixelFormatRgba8888]
	bytesPerPixel := uint32(nvColorFmt>>3) & 0x1f
	blockHeightLog2 := uint32(4) // According to TRM this is the optimal value (SIXTEEN_GOBS)
	blockHeight := uint32(8 * (1 << blockHeightLog2))

	grBuf := nv.GraphicBuffer{}
	grBuf.Header.NumInts = int32(uint64(unsafe.Sizeof(grBuf)-unsafe.Sizeof(nxtypes.NativeHandle{})) / 4)
	grBuf.Unk0 = -1
	grBuf.Magic = 0xDAFFCAFF
	grBuf.PID = 42
	grBuf.Usage = GRALLOC_USAGE_HW_COMPOSER | GRALLOC_USAGE_HW_RENDER | GRALLOC_USAGE_HW_TEXTURE
	grBuf.Format = uint32(format)
	grBuf.ExtFormat = grBuf.Format
	grBuf.NumPlanes = 1

	grBuf.Planes[0].Width = width
	grBuf.Planes[0].Height = height
	grBuf.Planes[0].ColorFormat = nvColorFmt
	grBuf.Planes[0].Layout = nv.LayoutBlockLinear
	grBuf.Planes[0].Kind = nv.KindPitch
	grBuf.Planes[0].BlockHeightLog2 = blockHeightLog2

	widthAlignedBytes := (width*bytesPerPixel + 63) & ^uint32(63) // GOBs are 64 bytes wide
	widthAligned := widthAlignedBytes / bytesPerPixel
	heightAligned := (height + blockHeight - 1) & ^(blockHeight - 1)
	fbSize := widthAlignedBytes * heightAligned
	bufferSize := int(((uint32(numFbs) * fbSize) + 0xFFF) & ^uint32(0xFFF))

	surface = &Surface{
		LayerId:        layerId,
		IGBP:           igbp,
		State:          SURFACE_STATE_INVALID,
		CurrentSlot:    0,
		GraphicBuffers: []GraphicBuffer{},
		GRBuff:         &grBuf,
	}

	surface.GpuBufferMemory = memory.AllocPages(bufferSize, bufferSize)
	if surface.GpuBufferMemory == nil {
		return nil, status, nxerrors.OutOfMemory
	}

	if debugDisplay {
		fmt.Printf("Allocated %d bytes\n", len(surface.GpuBufferMemory))
	}

	surface.GpuBuffer, err = gpu.CreateBuffer(unsafe.Pointer(&surface.GpuBufferMemory[0]), uintptr(bufferSize), 0, 0x1000, grBuf.Planes[0].Kind)
	if err != nil {
		return nil, status, err
	}

	nvmapId, err := surface.GpuBuffer.GetID()
	if err != nil {
		_, _, _ = surface.GpuBuffer.Destroy()
		return nil, status, err
	}

	grBuf.NVMapID = int32(nvmapId)
	grBuf.Stride = widthAligned
	grBuf.TotalSize = fbSize
	grBuf.Planes[0].Pitch = widthAlignedBytes
	grBuf.Planes[0].Size = uint64(fbSize)

	for i := 0; i < numFbs; i++ {
		grBuf.Planes[0].Offset = uint32(i) * fbSize
		surface.GraphicBuffers = append(surface.GraphicBuffers, GraphicBuffer{
			GRBuff:    &grBuf,
			Width:     grBuf.Planes[0].Width,
			Height:    grBuf.Planes[0].Height,
			Stride:    grBuf.Stride,
			Format:    format,
			Length:    fbSize,
			Usage:     grBuf.Usage,
			GPUBuffer: surface.GpuBuffer,
		})
		surface.HasRequested = append(surface.HasRequested, false)
		if debugDisplay {
			fmt.Printf("Pre-allocating buffer %d\n", i)
		}
		err = IGBPSetPreallocatedBuffer(surface.IGBP, i, &surface.GraphicBuffers[i])
		if err != nil {
			return nil, 0, err
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
