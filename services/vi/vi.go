package vi

import (
	"fmt"
	"github.com/racerxdl/gonx/internal"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/services/ipc"
	"github.com/racerxdl/gonx/services/sm"
	"unsafe"
)

const (
	viDebug = false
)

var viDomain *ipc.Domain
var imrsObject ipc.Object   // nn::visrv::sf::IManagerRootService
var iadsObject ipc.Object   // nn::visrv::sf::IApplicationDisplayService
var imdsObject ipc.Object   // nn::visrv::sf::IManagerDisplayService
var isdsObject ipc.Object   // nn::visrv::sf::ISystemDisplayService
var ihosbdObject ipc.Object // nn::visrv::sf::IHOSBinderDriver

var viInitializations = 0

func GetObject(iface ipc.Object, command int) (ipc.Object, error) {
	if viDebug {
		fmt.Printf("VI::GetObject(--, %d)\n", command)
	}
	if viInitializations <= 0 {
		return ipc.Object{}, nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(uint32(command))
	rs := ipc.ResponseFmt{}
	rs.Objects = make([]ipc.Object, 1)

	err := ipc.Send(iface, &rq, &rs)
	if err != nil {
		return ipc.Object{}, err
	}

	return rs.Objects[0], nil
}

func Init() (err error) {
	if viDebug {
		println("VI::Init()")
	}
	viInitializations++
	if viInitializations > 1 {
		return nil
	}

	smInit := false
	imrsInit := false
	domainInit := false
	iadsInit := false
	ihosbdInit := false
	isdsInit := false

	defer func() {
		if err != nil {
			viInitializations--

			if isdsInit {
				_ = ipc.Close(isdsObject)
			}
			if ihosbdInit {
				_ = ipc.Close(ihosbdObject)
			}
			if iadsInit {
				_ = ipc.Close(iadsObject)
			}
			if imrsInit {
				_ = ipc.Close(imrsObject)
			}
			if domainInit {
				_ = ipc.CloseSession(viDomain.Session)
			}
		}

		if smInit {
			sm.Finalize()
		}
	}()

	if viDebug {
		println("VI::Init() - SM Initialize")
	}
	// SM Initialize
	err = sm.Init()
	if err != nil {
		return err
	}

	smInit = true

	if viDebug {
		println("VI::Init() - VI:M Initialize")
	}
	// vi:m initialize
	err = sm.GetService(&imrsObject, "vi:m")
	if err != nil {
		return err
	}

	imrsInit = true

	if viDebug {
		println("VI::Init() - VI:M Convert to Domain")
	}
	// Domain Initialize
	viDomain, err = ipc.ConvertToDomain(&imrsObject)
	if err != nil {
		return err
	}
	domainInit = true

	if viDebug {
		println("VI::Init() - IADS Init")
	}
	// iads initialize
	rq := ipc.MakeDefaultRequest(2)
	rq.SetRawDataFromUint32Slice([]uint32{1})

	rs := ipc.ResponseFmt{}
	rs.Objects = make([]ipc.Object, 1)

	err = ipc.Send(imrsObject, &rq, &rs)
	if err != nil {
		return err
	}

	iadsObject = rs.Objects[0]
	iadsInit = true

	if viDebug {
		println("VI::Init() - IHOSBD Init")
	}
	ihosbdObject, err = GetObject(iadsObject, 100)
	if err != nil {
		return err
	}
	ihosbdInit = true

	if viDebug {
		println("VI::Init() - ISDS Init")
	}
	isdsObject, err = GetObject(iadsObject, 101)
	if err != nil {
		return err
	}
	isdsInit = true

	if viDebug {
		println("VI::Init() - IMDS Init")
	}
	imdsObject, err = GetObject(iadsObject, 102)
	if err != nil {
		return err
	}

	return nil
}

func TransactParcel(handle int32, transaction, flags uint32, rqParcel []byte, rsParcel []byte) error {
	if viDebug {
		fmt.Printf("VI::TransactParcel(%d, %d, %d, [%d]byte, [%d]byte)\n", handle, transaction, flags, len(rqParcel), len(rsParcel))
	}
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rqBuffer := ipc.Buffer{
		Addr: uintptr(unsafe.Pointer(&rqParcel[0])),
		Size: uint64(len(rqParcel)),
		Type: 5,
	}

	rsBuffer := ipc.Buffer{
		Addr: uintptr(unsafe.Pointer(&rsParcel[0])),
		Size: uint64(len(rsParcel)),
		Type: 6,
	}

	raw := struct {
		handle      int32
		transaction uint32
		flags       uint32
	}{
		handle:      handle,
		transaction: transaction,
		flags:       flags,
	}

	rq := ipc.MakeDefaultRequest(0)
	rq.Buffers = []*ipc.Buffer{&rqBuffer, &rsBuffer}
	rq.RawData = make([]byte, unsafe.Sizeof(raw))
	internal.Memcpy(unsafe.Pointer(&rq.RawData[0]), unsafe.Pointer(&raw), uintptr(len(rq.RawData)))

	rs := ipc.ResponseFmt{}

	return ipc.Send(ihosbdObject, &rq, &rs)
}

func AdjustRefCount(handle, addVal, Type int32) error {
	if viDebug {
		fmt.Printf("VI::AdjustRefCount(%d, %d, %d)\n", handle, addVal, Type)
	}
	if viInitializations <= 0 {
		return nxerrors.VINotInitialized
	}

	rq := ipc.MakeDefaultRequest(1)
	rq.SetRawDataFromUint32Slice([]uint32{uint32(handle), uint32(addVal), uint32(Type)})

	rs := ipc.ResponseFmt{}

	return ipc.Send(ihosbdObject, &rq, &rs)
}

func forceFinalize() {
	if viDebug {
		println("VI::ForceFinalize()")
	}
	_ = ipc.Close(isdsObject)
	_ = ipc.Close(ihosbdObject)
	_ = ipc.Close(iadsObject)
	_ = ipc.Close(imrsObject)
	_ = ipc.CloseSession(viDomain.Session)
	viInitializations = 0
}

func Finalize() {
	if viDebug {
		println("VI::Finalize()")
	}
	viInitializations--
	if viInitializations < 0 {
		forceFinalize()
	}
}
