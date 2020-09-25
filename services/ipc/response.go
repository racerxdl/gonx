package ipc

import (
	"github.com/racerxdl/gonx/nx/nxtypes"
)

// ResponseFmt Describes format expectations for an incoming IPC response
//
// Represents the expectations for an IPC response and contains pointers to buffers for
// response data to be written to.
type ResponseFmt struct {
	CopyHandles []nxtypes.Handle
	MoveHandles []nxtypes.Handle
	Objects     []Object
	RawData     []byte
	HasPID      bool
	PID         *uint64
}
