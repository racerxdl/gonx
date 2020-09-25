package am

import (
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/ipc"
	"github.com/racerxdl/gonx/services/sm"
)

var amInitializations = 0
var proxyServiceObject ipc.Object
var proxyObject ipc.Object
var iscObject ipc.Object
var iwcObject ipc.Object
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
				_ = ipc.Close(iscObject)
			}
			if proxyInit {
				_ = ipc.Close(proxyObject)
			}
			if proxyServiceInit {
				_ = ipc.Close(proxyServiceObject)
			}
			if domainInit {
				_ = ipc.CloseSession(amDomain.Session)
			}
		}
		sm.Finalize()
	}()

	err = sm.Init()
	if err != nil {
		return err
	}

	err = sm.GetService(&proxyServiceObject, "appletAE")
	if err != nil {
		return err
	}
	proxyServiceInit = true

	amDomain, err = ipc.ConvertToDomain(&proxyServiceObject)
	if err != nil {
		return err
	}
	domainInit = true

	// Open Application Proxy
	rq := ipc.MakeDefaultRequest(200)
	rq.SetRawDataFromUint64(0)
	rq.SendPID = true
	rq.CopyHandles = []nxtypes.Handle{0xFFFF8001}

	rs := ipc.ResponseFmt{}
	rs.Objects = make([]ipc.Object, 1)

	err = ipc.Send(proxyServiceObject, &rq, &rs)
	if err != nil {
		return err
	}

	proxyObject = rs.Objects[0]
	proxyInit = true

	iscObject, err = GetObject(proxyObject, 1)
	if err != nil {
		return err
	}
	iscInit = true

	iwcObject, err = GetObject(proxyObject, 2)
	return err
}

func forceFinalize() {
	_ = ipc.Close(iwcObject)
	_ = ipc.Close(iscObject)
	_ = ipc.Close(proxyObject)
	_ = ipc.Close(proxyServiceObject)
	_ = ipc.CloseSession(amDomain.Session)
	amInitializations = 0
}

func Finalize() {
	amInitializations--
	if amInitializations <= 0 {
		forceFinalize()
	}
}
