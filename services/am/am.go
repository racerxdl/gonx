package am

import (
	"github.com/racerxdl/gonx/nx/env"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/ipc"
	"github.com/racerxdl/gonx/services/sm"
	"github.com/racerxdl/gonx/svc"
	"time"
)

var debug = false
var amInitializations = 0
var proxyServiceObject ipc.Object
var proxyObject ipc.Object
var iSelfControllerObject ipc.Object
var iWindowController ipc.Object
var amDomain *ipc.Domain

func GetObject(iface ipc.Object, command int) (ipc.Object, error) {
	if amInitializations <= 0 {
		return ipc.Object{}, nxerrors.AMNotInitialized
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
	if debug {
		println("am::Init()")
	}
	amInitializations++
	if amInitializations > 1 {
		return nil
	}

	proxyServiceInit := false
	domainInit := false
	iscInit := false
	proxyInit := false

	defer func() {
		if err != nil {
			amInitializations--
			if iscInit {
				_ = ipc.Close(&iSelfControllerObject)
			}
			if proxyInit {
				_ = ipc.Close(&proxyObject)
			}
			if proxyServiceInit {
				_ = ipc.Close(&proxyServiceObject)
			}
			if domainInit {
				_ = ipc.CloseSession(amDomain.Session)
				amDomain = nil
			}
		}
		sm.Finalize()
	}()

	err = sm.Init()
	if err != nil {
		return err
	}

	appletType := env.GetAppletType()

	switch appletType {
	case nxtypes.AppletTypeDefault:
		appletType = nxtypes.AppletTypeApplication
		fallthrough
	case nxtypes.AppletTypeApplication:
		err = sm.GetService(&proxyServiceObject, "appletOE")
	default:
		err = sm.GetService(&proxyServiceObject, "appletAE")
	}

	if err != nil {
		return err
	}
	proxyServiceInit = true

	amDomain, err = ipc.ConvertToDomain(&proxyServiceObject)
	if err != nil {
		return err
	}
	domainInit = true

	cmdId := uint32(0)
	switch appletType {
	case nxtypes.AppletTypeApplication:
		cmdId = 0
	case nxtypes.AppletTypeSystemApplet:
		cmdId = 100
	case nxtypes.AppletTypeLibraryApplet:
		cmdId = 200
	case nxtypes.AppletTypeOverlayApplet:
		cmdId = 300
	case nxtypes.AppletTypeSystemApplication:
		cmdId = 350
	default:
		return nxerrors.UnknownAppletType
	}

	// Open Application Proxy
	rq := ipc.Request{}
	rs := ipc.ResponseFmt{}
	resCode := uint64(nxerrors.AMBusy)
	for resCode == nxerrors.AMBusy {
		rq = ipc.MakeDefaultRequest(cmdId)
		rq.SetRawDataFromUint64(0)
		rq.SendPID = true
		rq.CopyHandles = []nxtypes.Handle{svc.CurrentProcessHandle}

		rs = ipc.ResponseFmt{}
		rs.Objects = make([]ipc.Object, 1)

		err = ipc.Send(proxyServiceObject, &rq, &rs)
		if err != nil {
			ipcErr, ok := err.(nxerrors.IPCError)
			if !ok || ipcErr.Result != nxerrors.AMBusy {
				return err
			}

			time.Sleep(time.Second)
		} else {
			resCode = 0
		}
	}

	proxyObject = rs.Objects[0]
	proxyInit = true

	iSelfControllerObject, err = GetObject(proxyObject, 1)
	if err != nil {
		return err
	}
	iscInit = true

	iWindowController, err = GetObject(proxyObject, 2)
	return err
}

func forceFinalize() {
	if debug {
		println("am::ForceFinalize()")
	}

	_ = ipc.Close(&iWindowController)
	_ = ipc.Close(&iSelfControllerObject)
	_ = ipc.Close(&proxyObject)
	_ = ipc.Close(&proxyServiceObject)

	if amDomain != nil {
		_ = ipc.CloseSession(amDomain.Session)
		amDomain = nil
	}
	amInitializations = 0
}

func Finalize() {
	if debug {
		println("am::Finalize()")
	}
	amInitializations--
	if amInitializations <= 0 {
		forceFinalize()
	}
}
