package vi

import (
	"encoding/binary"
	"fmt"
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
	if viDebug {
		fmt.Printf("VI::OpenDisplay(%s)\n", name)
	}
	if viInitializations <= 0 {
		return nil, nxerrors.VINotInitialized
	}

	raw := make([]byte, 0x40)
	copy(raw, name)
	l := len(name) + 1
	if l > 0x40 {
		l = 0x40
	}
	raw[l-1] = 0x00 // force null-terminated

	if viDebug {
		println("VI::OpenDisplay() - IPC Call")
		fmt.Printf("%+v\n", raw)
	}
	rq := ipc.MakeDefaultRequest(requestOpenDisplay)
	rq.RawData = raw

	disp := &Display{}

	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 8) // uint64

	err := ipc.Send(iadsObject, &rq, &rs)
	if err != nil {
		if viDebug {
			fmt.Printf("Error calling IPC: %s\n", err)
		}
		return nil, err
	}

	disp.ID = binary.LittleEndian.Uint64(rs.RawData)

	return disp, nil
}

func CloseDisplay(d *Display) error {
	if viDebug {
		fmt.Printf("VI::CloseDisplay(%d)\n", d.ID)
	}
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(requestCloseDisplay)
	rq.SetRawDataFromUint64(d.ID)

	rs := ipc.ResponseFmt{}

	return ipc.Send(iadsObject, &rq, &rs)
}

func GetDisplayVsyncEvent(d *Display) error {
	if viDebug {
		fmt.Printf("VI::GetDisplayVsyncEvent(%d)\n", d.ID)
	}
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

func OpenLayer(displayName string, layerId uint64, aruid nxtypes.ARUID) (*IGBP, error) {
	if viDebug {
		fmt.Printf("VI::OpenLayer(%s, %d, %d)\n", displayName, layerId, aruid)
	}
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
		aruid       nxtypes.ARUID
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
		IgbpBinder: *b,
	}

	return igbp, nil
}

func CloseLayer(layerId uint64) error {
	if viDebug {
		fmt.Printf("VI::CloseLayer(%d)\n", layerId)
	}
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(requestCloseLayer)
	rq.SetRawDataFromUint64(layerId)

	rs := ipc.ResponseFmt{}

	return ipc.Send(iadsObject, &rq, &rs)
}

func DestroyManagedLayer(layerId uint64) error {
	if viDebug {
		fmt.Printf("VI::DestroyManagedLayer(%d)\n", layerId)
	}
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(requestDestroyManagedLayer)
	rq.SetRawDataFromUint64(layerId)

	rs := ipc.ResponseFmt{}

	return ipc.Send(imdsObject, &rq, &rs)
}

func IadsSetLayerScalingMode(scalingMode uint32, layerId uint64) error {
	if viDebug {
		fmt.Printf("VI::IadsSetLayerScalingMode(%d, %d)\n", scalingMode, layerId)
	}
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	params := struct {
		scalingMode uint32
		layerId     uint64
	}{
		scalingMode: scalingMode,
		layerId:     layerId,
	}

	buff := make([]byte, unsafe.Sizeof(params))
	internal.Memcpy(unsafe.Pointer(&buff[0]), unsafe.Pointer(&params), uintptr(len(buff)))

	rq := ipc.MakeDefaultRequest(requestSetLayerScalingMode)
	rq.RawData = buff

	rs := ipc.ResponseFmt{}

	return ipc.Send(imdsObject, &rq, &rs)
}
