package display

import (
	"github.com/racerxdl/gonx/internal"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/services/gpu"
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
	Format            PixelFormat
	Usage             uint32
	GPUBuffer         *gpu.Buffer
	Index             int32
	PixelBufferOffset uint32
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
	gpuBufferId, err := gb.GPUBuffer.GetID()
	if err != nil {
		return err
	}

	gpuBufferCopy, err := gpu.InitializeFromId(gpuBufferId)
	if err != nil {
		return err
	}

	/*
	  RFBG, width, height, stride,
	  format, usage, 0x2a [mId >> 32?], [some kind of native handle thing?] [v38 points here] [mId & 0xFFFFFFFF] index
	  0x0 [numFds?], 0x51 [numInts?], -1 {592}, gpu_buffer_id {593},
	  0x0, 0xdaffcaff {582}, -1 [0x2a?] {583}, v39 [0x0] {584},
	  v5 [0xb00] {585}, v4 [0x1] {586}, v4 [0x1] {587}, 0 [0x500] {588}
	  v31 [0x3c0000] {589}, v22 [0x1] {590}, uninit?,
	  memcpied from &v53, length 88 * v22
	  zeroes? {581 clears from v38+12 to the end of this block}
	  0x0 {594}, 0x0 {594}
	*/

	template := []uint32{
		0x47424652, gb.Width, gb.Height, gb.Stride,
		uint32(gb.Format), gb.Usage, 0x0000002a, uint32(gb.Index),
		0x00000000, 0x00000051, 0xffffffff, gpuBufferId,
		0x00000000, 0xdaffcaff, 0x0000002a, 0x00000000,
		0x00000b00, 0x00000001, 0x00000001, 0x00000500,
		0x003c0000, 0x00000001, 0x00000000, 0x00000500,
		0x000002d0, 0x00532120, 0x00000001, 0x00000003,
		0x00001400, gpuBufferCopy.NvMapHandle, gb.PixelBufferOffset, 0x000000fe,
		0x00000004, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x003c0000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x0, 0x0, // these last two words are some kind of address but I don't think it really matters
	}
	p.WriteInPlaceU32(template)

	return nil
}
