package am

import (
	"encoding/binary"
	"github.com/racerxdl/gonx/nx/nxerrors"
	"github.com/racerxdl/gonx/services/ipc"
)

func IscCreateManagedDisplayLayer() (uint64, error) {
	if amInitializations <= 0 {
		return 0, nxerrors.AMNotInitialized
	}

	rq := ipc.MakeDefaultRequest(40)
	rs := ipc.ResponseFmt{}
	rs.RawData = make([]byte, 8) // one uint64

	err := ipc.Send(iscObject, &rq, &rs)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(rs.RawData), nil
}

func IscApproveToDisplay() error {
	if amInitializations <= 0 {
		return nxerrors.AMNotInitialized
	}

	rq := ipc.MakeDefaultRequest(51)
	rs := ipc.ResponseFmt{}

	return ipc.Send(iscObject, &rq, &rs)
}
