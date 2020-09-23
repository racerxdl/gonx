package ipc

import "github.com/racerxdl/gonx/nx/nxtypes"

// Maximum allocated descriptors is always 16
const maxIPCDescriptors = 16
const debug = false
const debugDumpBeforeSend = false

const (
	TypeIPCInvalid            Type = 0
	TypeIPCLegacyRequest      Type = 1
	TypeIPCClose              Type = 2
	TypeIPCLegacyControl      Type = 3
	TypeIPCRequest            Type = 4
	TypeIPCControl            Type = 5
	TypeIPCRequestWithContext Type = 6
	TypeIPCControlWithContext Type = 7
)

// Request Represents an unmarshalled outgoing IPC request
// see http://switchbrew.org/index.php?title=IPC_Marshalling#IPC_Command_Structure
type Request struct {
	Type        Type
	Buffers     []*Buffer
	RequestID   uint32
	RawData     []byte
	SendPID     bool
	CopyHandles []nxtypes.Handle
	MoveHandles []nxtypes.Handle
	Objects     []Object
	CloseObject bool
}

//// SetRawDataFromPointer copies the specified data to internal RawData field
//func (i *Request) SetRawDataFromPointer(ptr unsafe.Pointer, dataLen uintptr) {
//    i.RawData = make([]byte, dataLen)
//    Memcpy(unsafe.Pointer(&i.RawData[0]), ptr, dataLen)
//}

func (i *Request) SetRawDataFromUint64(data uint64) {
	i.RawData = make([]byte, 8)
	i.RawData[0] = byte(data)
	i.RawData[1] = byte(data >> 8)
	i.RawData[2] = byte(data >> 16)
	i.RawData[3] = byte(data >> 24)
	i.RawData[4] = byte(data >> 32)
	i.RawData[5] = byte(data >> 40)
	i.RawData[6] = byte(data >> 48)
	i.RawData[7] = byte(data >> 56)
}

func (i *Request) SetRawDataFromUint32Slice(data []uint32) {
	i.RawData = make([]byte, len(data)*4)
	for n, v := range data {
		i.RawData[n*4+0] = byte(v)
		i.RawData[n*4+1] = byte(v >> 8)
		i.RawData[n*4+2] = byte(v >> 16)
		i.RawData[n*4+3] = byte(v >> 24)
	}
}

func MakeDefaultRequest(requestId uint32) Request {
	return Request{
		Type:      TypeIPCRequest,
		RequestID: requestId,
	}
}
