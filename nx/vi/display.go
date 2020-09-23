package vi

import (
	"encoding/binary"
	"github.com/racerxdl/gonx/nx/internal"
	"github.com/racerxdl/gonx/nx/ipc"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"unsafe"
)

type Display struct {
	ID    uint64
	VSync nxtypes.ReventHandle
}

func OpenDisplay(name string) (*Display, error) {
	if viInitializations <= 0 {
		return nil, nxerrors.VINotInitialized
	}

	raw := make([]byte, 0x40)
	copy(raw, name)
	l := len(name)
	if l > 0x40 {
		l = 0x40
	}
	raw[l-1] = 0x00 // force null-terminated

	rq := ipc.MakeDefaultRequest(requestOpenDisplay)
	rq.RawData = raw

	disp := &Display{}

	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 8) // uint64

	err := ipc.Send(iadsObject, &rq, &rs)
	if err != nil {
		return nil, err
	}

	disp.ID = binary.LittleEndian.Uint64(rs.RawData)

	return disp, nil
}

func CloseDisplay(d *Display) error {
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(requestCloseDisplay)
	rq.SetRawDataFromUint64(d.ID)

	rs := ipc.ResponseFmt{}

	return ipc.Send(iadsObject, &rq, &rs)
}

func GetDisplayVsyncEvent(d *Display) error {
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(requestDisplayVsyncEvent)
	rq.SetRawDataFromUint64(d.ID)

	rs := ipc.ResponseFmt{}
	rs.CopyHandles = make([]nxtypes.Handle, 1)

	err := ipc.Send(iadsObject, &rq, &rs)
	if err != nil {
		return err
	}

	d.VSync = nxtypes.ReventHandle(rs.CopyHandles[0])

	return nil
}

func OpenLayer(displayName string, layerId, aruid uint64) (*IGBP, error) {
	if viInitializations <= 0 {
		return nil, nxerrors.VINotInitialized
	}

	parcelBuff := make([]byte, 0x210)

	ipcBuff := &ipc.Buffer{
		Addr: uintptr(unsafe.Pointer(&parcelBuff[0])),
		Size: 0x210,
		Type: 6,
	}

	rqArgs := struct {
		displayName [0x40]byte
		layerId     uint64
		aruid       uint64
	}{}

	rqArgs.layerId = layerId
	rqArgs.aruid = aruid
	copy(rqArgs.displayName[:], displayName)
	rqArgs.displayName[0x3F] = 0 // ensure at least last char is null

	buff := make([]byte, unsafe.Sizeof(rqArgs))
	internal.Memcpy(unsafe.Pointer(&buff[0]), unsafe.Pointer(&rqArgs), uintptr(len(buff)))

	rq := ipc.MakeDefaultRequest(requestOpenLayer)
	rq.Buffers = []*ipc.Buffer{ipcBuff}
	rq.RawData = buff
	rq.SendPID = true

	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 8) // one uint64

	err := ipc.Send(iadsObject, &rq, &rs)
	if err != nil {
		return nil, err
	}

	p, err := ParcelLoad(parcelBuff)
	if err != nil {
		return nil, err
	}

	b, err := p.ReadBinder()
	if err != nil {
		return nil, err
	}

	igbp := &IGBP{
		igbpBinder: Binder{
			handle: b.handle,
		},
	}

	return igbp, nil
}

func CloseLayer(layerId uint64) error {
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(requestCloseLayer)
	rq.SetRawDataFromUint64(layerId)

	rs := ipc.ResponseFmt{}

	return ipc.Send(iadsObject, &rq, &rs)
}
