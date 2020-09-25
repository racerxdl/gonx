package ipc

import (
	"github.com/racerxdl/gonx/nx/nxtypes"
)

// PackedMessage represents a IPC Message Data to be packed
type PackedMessage struct {
	Type Type

	Buffers []*Buffer

	DataSection []byte

	CopyHandles []nxtypes.Handle
	MoveHandles []nxtypes.Handle
	SendPID     bool
}

// Message Describes an incoming IPC message. Used as an intermediate during unpacking.
type Message struct {
	MessageType        uint16
	RawDataSectionSize uint32 // in Words
	NumXDescriptors    uint32
	NumADescriptors    uint32
	NumBDescriptors    uint32
	NumWDescriptors    uint32
	CDescriptorFlags   uint32
	XDescriptors       []uint32
	ADescriptors       []uint32
	BDescriptors       []uint32
	WDescriptors       []uint32
	CDescriptors       []uint32
	NumCopyHandles     uint32
	NumMoveHandles     uint32
	CopyHandles        []uint32
	MoveHandles        []uint32
	HasPID             bool
	PID                uint64
	PrePadding         int
	PostPadding        int
	DataSection        []uint32
}
