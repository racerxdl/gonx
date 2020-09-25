package display

import (
	"fmt"
	"github.com/racerxdl/gonx/nx/internal"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/vi"
	"unsafe"
)

const interfaceToken = "android.gui.IGraphicBufferProducer"
const (
	REQUEST_BUFFER          = 0x1
	SET_BUFFER_COUNT        = 0x2
	DEQUEUE_BUFFER          = 0x3
	DETACH_BUFFER           = 0x4
	DETACH_NEXT_BUFFER      = 0x5
	ATTACH_BUFFER           = 0x6
	QUEUE_BUFFER            = 0x7
	CANCEL_BUFFER           = 0x8
	QUERY                   = 0x9
	CONNECT                 = 0xA
	DISCONNECT              = 0xB
	ALLOCATE_BUFFERS        = 0xD
	SET_PREALLOCATED_BUFFER = 0xE
)

type DisconnectMode int

const (
	DisconnectAPI      DisconnectMode = 0
	DisconnectAllLocal DisconnectMode = 1
)

func IGBPDisconnect(igbp vi.IGBP, api int, mode DisconnectMode) (status int, err error) {
	if debugDisplay {
		fmt.Printf("IGBPDisconnect(%d, %d, %d)\n", igbp.IgbpBinder.Handle, api, mode)
	}
	p := &vi.Parcel{}

	p.WriteInterfaceToken(interfaceToken)
	p.WriteU32(uint32(api))
	p.WriteU32(uint32(mode))

	response, err := vi.BinderTransactParcel(igbp.IgbpBinder, DISCONNECT, 0, p)
	if err != nil {
		return 0, err
	}

	if response.Remaining() < 4 {
		return 0, nxerrors.ParcelDataUnderrun
	}

	status = int(response.ReadU32())

	return status, nil
}

func IGBPConnect(igbp vi.IGBP, api int, producerControlledByApp bool) (status int, qbo *QueueBufferOutput, err error) {
	if debugDisplay {
		fmt.Printf("IGBPConnect(%d, %d, %t)\n", igbp.IgbpBinder.Handle, api, producerControlledByApp)
	}
	p := &vi.Parcel{}

	p.WriteInterfaceToken(interfaceToken)
	p.WriteU32(0) // IProducerListener is null
	p.WriteU32(uint32(api))

	producerControlledByAppInt := uint32(0)
	if producerControlledByApp {
		producerControlledByAppInt = 1
	}

	p.WriteU32(producerControlledByAppInt)

	response, err := vi.BinderTransactParcel(igbp.IgbpBinder, CONNECT, 0, p)
	if err != nil {
		return 0, nil, err
	}

	qbo, err = UnflattenQueueBufferOutput(response)
	if err != nil {
		return 0, nil, err
	}

	if response.Remaining() < 4 {
		return 0, nil, nxerrors.ParcelDataUnderrun
	}

	status = int(response.ReadU32())

	return status, qbo, nil
}

func IGBPRequestBuffer(igbp vi.IGBP, slot uint32) (status uint32, gb *GraphicBuffer, err error) {
	if debugDisplay {
		fmt.Printf("IGBPRequestBuffer(%d, %d)\n", igbp.IgbpBinder.Handle, slot)
	}
	p := &vi.Parcel{}
	p.WriteInterfaceToken(interfaceToken)
	p.WriteU32(slot)

	response, err := vi.BinderTransactParcel(igbp.IgbpBinder, REQUEST_BUFFER, 0, p)
	if err != nil {
		return 0, nil, err
	}

	nonNull := response.ReadU32() != 0
	if nonNull {
		length := response.ReadU32()
		if length != 0x16c {
			return 0, nil, nxerrors.DisplayGraphicBufferLengthMismatch
		}
		_ = response.ReadU32() // PixelBufferOffset
		gbBuff := response.ReadInPlace(0x16c)
		gb = &GraphicBuffer{}
		internal.Memcpy(unsafe.Pointer(gb), unsafe.Pointer(&gbBuff[0]), 0x16c)
	}

	status = response.ReadU32()

	return status, gb, err
}

func IGBPSetPreallocatedBuffer(igbp vi.IGBP, slot int, gb *GraphicBuffer) error {
	if debugDisplay {
		fmt.Printf("IGBPSetPreallocatedBuffer(%d, %d, %p)\n", igbp.IgbpBinder.Handle, slot, gb)
	}
	p := &vi.Parcel{}
	p.WriteInterfaceToken(interfaceToken)
	p.WriteU32(uint32(slot))
	p.WriteU32(1) // Unknown

	p.WriteU32(0x16C)
	p.WriteU32(0)

	err := gb.Flatten(p)
	if err != nil {
		return err
	}

	_, err = vi.BinderTransactParcel(igbp.IgbpBinder, SET_PREALLOCATED_BUFFER, 0, p)
	return err
}

func IGBPDequeueBuffer(igbp vi.IGBP, width, height uint32, pixelFormat PixelFormat, usage uint32, getFrameTimestamps bool) (status, slot uint32, fence Fence, outTimestamps *FrameEventHistoryDelta, err error) {
	if debugDisplay {
		fmt.Printf("IGBPDequeueBuffer(%d, %d, %d, %d, %d, %t)\n", igbp.IgbpBinder.Handle, width, height, pixelFormat, usage, getFrameTimestamps)
	}
	if getFrameTimestamps {
		return status, slot, fence, outTimestamps, nxerrors.NotImplemented
	}
	var response *vi.Parcel

	p := &vi.Parcel{}
	p.WriteInterfaceToken(interfaceToken)
	p.WriteU32(uint32(pixelFormat))
	p.WriteU32(width)
	p.WriteU32(height)
	p.WriteU32(0) // getFrameTimestamps
	p.WriteU32(usage)

	response, err = vi.BinderTransactParcel(igbp.IgbpBinder, DEQUEUE_BUFFER, 0, p)
	if err != nil {
		return status, slot, fence, outTimestamps, nxerrors.NotImplemented
	}

	slot = response.ReadU32()
	hasFence := response.ReadU32() > 0

	if hasFence {
		fence, err = UnflattenFence(response)
		if err != nil {
			return status, slot, fence, outTimestamps, nxerrors.NotImplemented
		}
	}

	status = response.ReadU32()

	return status, slot, fence, outTimestamps, nil
}

func IGBPQueueBuffer(igbp vi.IGBP, slot int, qbi *QueueBufferInput) (qbo *QueueBufferOutput, status int, err error) {
	if debugDisplay {
		fmt.Printf("IGBPQueueBuffer(%d, %d, %p)\n", igbp.IgbpBinder.Handle, slot, qbi)
	}
	p := &vi.Parcel{}
	p.WriteInterfaceToken(interfaceToken)
	p.WriteU32(uint32(slot))
	qbi.Flatten(p)

	resp, err := vi.BinderTransactParcel(igbp.IgbpBinder, QUEUE_BUFFER, 0, p)
	if err != nil {
		return nil, 0, err
	}

	qbo, err = UnflattenQueueBufferOutput(resp)
	if err != nil {
		return nil, 0, err
	}

	status = int(resp.ReadU32())

	return qbo, status, nil
}
