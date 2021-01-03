package display

import (
	"github.com/racerxdl/gonx/internal"
	"github.com/racerxdl/gonx/nx/graphics"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/gpu"
	"github.com/racerxdl/gonx/services/nv"
	"github.com/racerxdl/gonx/services/vi"
	"unsafe"
)

// CompositorTiming
// https://source.android.com/reference/hidl/android/hardware/graphics/bufferqueue/1.0/IGraphicBufferProducer#compositortiming
type CompositorTiming struct {
	DeadlineNanoseconds       int64
	InternalNanoseconds       int64
	PresentLatencyNanoseconds int64
}

// FrameEventHistoryDelta
//
// Not fully implemented, just enough that it works.
// https://source.android.com/reference/hidl/android/hardware/graphics/bufferqueue/1.0/IGraphicBufferProducer#frameeventhistorydelta
type FrameEventHistoryDelta struct {
	CompositorTiming CompositorTiming
}

// QueueBufferOutput Values received back from queueBuffer
// https://source.android.com/reference/hidl/android/hardware/graphics/bufferqueue/1.0/IGraphicBufferProducer#queuebufferoutput
type QueueBufferOutput struct {
	Width             uint32
	Height            uint32
	TransformHint     uint32
	NumPendingBuffers uint32
	NextFrameNumber   uint32
	BufferReplaced    bool
	FrameTimestamps   FrameEventHistoryDelta
}

// GraphicBuffer Graphics Buffer
type GraphicBuffer struct {
	Width             uint32
	Height            uint32
	Stride            uint32
	Format            graphics.PixelFormat
	Length            uint32
	Usage             uint32
	GPUBuffer         *gpu.Buffer
	Index             int32
	PixelBufferOffset uint32

	NativeHandle *nxtypes.NativeHandle
	GRBuff       *nv.GraphicBuffer
}

// QueueBufferInput Parameters passed to queueBuffer
// https://source.android.com/reference/hidl/android/hardware/graphics/bufferqueue/1.0/IGraphicBufferProducer#queuebufferinput
type QueueBufferInput struct {
	Size            uint32
	NumFds          uint32
	Timestamp       int64
	IsAutoTimestamp int32
	Crop            Rect
	ScalingMode     int32
	Transform       uint32
	StickyTransform uint32
	Unknown         [2]uint32
	Fence           Fence
}

func UnflattenQueueBufferOutput(p *vi.Parcel) (qbo *QueueBufferOutput, err error) {
	if p.Remaining() < 4*4 { // 4 uint32
		return nil, nxerrors.ParcelDataUnderrun
	}

	qbo = &QueueBufferOutput{
		Width:             p.ReadU32(),
		Height:            p.ReadU32(),
		TransformHint:     p.ReadU32(),
		NumPendingBuffers: p.ReadU32(),
	}

	return qbo, nil
}

func (qbi *QueueBufferInput) Flatten(p *vi.Parcel) {
	buff := make([]byte, unsafe.Sizeof(*qbi))
	internal.Memcpy(unsafe.Pointer(&buff[0]), unsafe.Pointer(qbi), uintptr(len(buff)))
	p.WriteInPlace(buff)
}

func (gb *GraphicBuffer) Flatten(p *vi.Parcel) error {
	buffer := make([]uint32, 0)

	buffer = append(buffer, 0x47424652) // GBFR (Graphic Buffer)
	buffer = append(buffer, gb.Width)
	buffer = append(buffer, gb.Height)
	buffer = append(buffer, gb.Stride)

	buffer = append(buffer, uint32(gb.Format))
	buffer = append(buffer, gb.Usage)
	buffer = append(buffer, 42)
	buffer = append(buffer, 0)

	buffer = append(buffer, 0)
	buffer = append(buffer, uint32(gb.GRBuff.Header.NumInts))

	buffer = append(buffer, uint32(gb.GRBuff.Unk0))
	buffer = append(buffer, uint32(gb.GRBuff.NVMapID))
	buffer = append(buffer, gb.GRBuff.Unk2)
	buffer = append(buffer, gb.GRBuff.Magic)
	buffer = append(buffer, gb.GRBuff.PID)
	buffer = append(buffer, gb.GRBuff.Type)
	buffer = append(buffer, gb.GRBuff.Usage)
	buffer = append(buffer, gb.GRBuff.Format)
	buffer = append(buffer, gb.GRBuff.ExtFormat)
	buffer = append(buffer, gb.GRBuff.Stride)
	buffer = append(buffer, gb.GRBuff.TotalSize)
	buffer = append(buffer, gb.GRBuff.NumPlanes)
	buffer = append(buffer, gb.GRBuff.Unk12)

	planeSize := (uint32(unsafe.Sizeof(gb.GRBuff.Planes[0])) / 4) * 3
	planeBuff := make([]uint32, planeSize)
	internal.Memcpy(unsafe.Pointer(&planeBuff[0]), unsafe.Pointer(&gb.GRBuff.Planes), unsafe.Sizeof(gb.GRBuff.Planes))

	for i := uint32(0); i < planeSize; i++ {
		buffer = append(buffer, planeBuff[i])
	}
	buffer = append(buffer, 0)
	buffer = append(buffer, 0)

	p.WriteU32(uint32(len(buffer) * 4))
	p.WriteU32(0)
	p.WriteInPlaceU32(buffer)

	return nil
}
