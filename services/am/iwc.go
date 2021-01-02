package am

import (
	"encoding/binary"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/nx/nxtypes"
	"github.com/racerxdl/gonx/services/ipc"
)

func IwcGetAppletResourceUserId() (nxtypes.ARUID, error) {
	if amInitializations <= 0 {
		return 0, nxerrors.AMNotInitialized
	}

	rq := ipc.MakeDefaultRequest(1)
	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 8) // one uint64

	err := ipc.Send(iwcObject, &rq, &rs)
	if err != nil {
		return 0, err
	}

	return nxtypes.ARUID(binary.LittleEndian.Uint64(rs.RawData)), nil
}

func IwcAcquireForegroundRights() error {
	if amInitializations <= 0 {
		return nxerrors.AMNotInitialized
	}

	rq := ipc.MakeDefaultRequest(10)
	rs := ipc.ResponseFmt{}

	return ipc.Send(iwcObject, &rq, &rs)
}
