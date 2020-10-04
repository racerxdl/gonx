package nv

import (
	"encoding/binary"
	"fmt"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/ipc"
	"github.com/racerxdl/gonx/services/sm"
	"github.com/racerxdl/gonx/svc"
	"unsafe"
)

const nvDebug = false
const transferMemSize = 3 * 1024 * 1024

//go:align 4096
var gpuSharedBuffer [transferMemSize]byte

var nvObject ipc.Object
var transferMem = nxtypes.Handle(0)
var nvInitializations = 0

func Init() (err error) {
	if nvDebug {
		println("NV::Init()")
	}
	nvInitializations++

	if nvInitializations > 1 {
		return nil
	}

	smInitialize := false
	nvsInitialize := false
	memInitialize := false

	defer func() {
		if err != nil {
			println("got error: %s", err)
			nvInitializations--
			// Only clean-up this if errored
			if memInitialize {
				svc.CloseHandle(transferMem)
			}

			if nvsInitialize {
				_ = ipc.Close(&nvObject)
			}
		}

		// Always de-init sm
		if smInitialize {
			sm.Finalize()
		}
	}()

	err = sm.Init()
	if err != nil {
		return fmt.Errorf("error initializing sm: %s", err)
	}
	smInitialize = true

	//err = sm.GetService(&nvObject, "nvdrv:a")
	err = sm.GetService(&nvObject, "nvdrv")
	if err != nil {
		return fmt.Errorf("error getting \"nvdrv:a\": %s", err)
	}
	nvsInitialize = true

	r := svc.CreateTransferMemory(&transferMem, uintptr(unsafe.Pointer(&gpuSharedBuffer[0])), transferMemSize, 0)

	if r != nxtypes.ResultOK {
		return fmt.Errorf("fail to create transfer memory: result code %d", r)
	}
	memInitialize = true

	handles := []nxtypes.Handle{0xFFFF8001, transferMem}

	rq := ipc.MakeDefaultRequest(3)
	rq.SetRawDataFromUint32Slice([]uint32{transferMemSize})
	rq.CopyHandles = handles

	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 4) // one uint32

	err = ipc.Send(nvObject, &rq, &rs)

	if err != nil {
		return fmt.Errorf("error sending ipc: %s", err)
	}

	response := []uint32{
		binary.LittleEndian.Uint32(rs.RawData[:4]),
	}

	if response[0] != 0 {
		return fmt.Errorf("nvidia error: %d", response[0])
	}

	return nil
}

func nvForceFinalize() {
	if nvDebug {
		println("NV::ForceFinalize()")
	}
	if transferMem != 0 {
		svc.CloseHandle(transferMem)
	}
	_ = ipc.Close(&nvObject)
	nvObject = ipc.Object{}
	nvInitializations = 0
}

func Finalize() {
	if nvDebug {
		println("NV::Finalize()")
	}
	nvInitializations--
	if nvInitializations <= 0 {
		nvForceFinalize()
	}
}

func Open(path string) (int32, error) {
	if nvDebug {
		fmt.Printf("NV::Open(%s)\n", path)
	}
	if nvInitializations <= 0 {
		return -1, nxerrors.NVNotInitialized
	}

	bytePath := []byte(path)

	buff := ipc.Buffer{}
	buff.Type = 0x5
	buff.Size = uint64(len(bytePath))
	buff.Addr = uintptr(unsafe.Pointer(&bytePath[0]))

	rq := ipc.MakeDefaultRequest(0)
	rq.Buffers = append(rq.Buffers, &buff)

	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 2*4) // Two uint32

	err := ipc.Send(nvObject, &rq, &rs)
	if err != nil {
		return -1, err
	}

	response := []uint32{
		binary.LittleEndian.Uint32(rs.RawData[:4]),
		binary.LittleEndian.Uint32(rs.RawData[4:]),
	}

	if response[1] != 0 {
		return int32(response[1]), fmt.Errorf("nvopen failed: %d", response[1])
	}

	return int32(response[0]), nil
}

func Close(fd int32) error {
	if nvDebug {
		fmt.Printf("NV::Close(%d)\n", fd)
	}
	if nvInitializations <= 0 {
		return nxerrors.NVNotInitialized
	}
	rq := ipc.MakeDefaultRequest(2)
	rq.SetRawDataFromUint32Slice([]uint32{uint32(fd)})

	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 4) // one uint32

	err := ipc.Send(nvObject, &rq, &rs)
	if err != nil {
		return fmt.Errorf("error on nvClose: %s", err)
	}

	res := binary.LittleEndian.Uint32(rs.RawData)

	if res != 0 {
		return fmt.Errorf("error on nvClose: result code: %d", res)
	}

	return nil
}

func Ioctl(fd int32, rqid uint32, arg unsafe.Pointer, size uintptr) (uint32, error) {
	if nvDebug {
		fmt.Printf("NV::Ioctl(%d, %d, %p, %d)\n", fd, rqid, arg, size)
	}
	if nvInitializations <= 0 {
		return 0xFFFFFFFF, nxerrors.NVNotInitialized
	}
	InB := ipc.Buffer{
		Addr: uintptr(arg),
		Size: uint64(size),
		Type: 0x21,
	}

	OutB := ipc.Buffer{
		Addr: uintptr(arg),
		Size: uint64(size),
		Type: 0x22,
	}

	rq := ipc.MakeDefaultRequest(1)
	rq.Buffers = []*ipc.Buffer{&InB, &OutB}
	rq.SetRawDataFromUint32Slice([]uint32{uint32(fd), rqid, 0, 0})

	rs := ipc.ResponseFmt{}
	rs.RawData = []byte{0, 0, 0, 0} // One uint32

	err := ipc.Send(nvObject, &rq, &rs)
	if err != nil {
		return 0, err
	}

	res := binary.LittleEndian.Uint32(rs.RawData)
	return res, nil
}
