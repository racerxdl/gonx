package ipc

import (
	"github.com/racerxdl/gonx/nx/nxtypes"
	"unsafe"
)

// IPCObject Represents either an object within an IPC domain or a standalone object
type Object struct {
	ObjectID   int32 // -1 if this represents a session, >= 0 if this represents a domain object
	Content    uint64
	IsBorrowed bool
}

func (o *Object) SetSession(session nxtypes.Handle) {
	o.Content = uint64(session)
	o.ObjectID = -1
}

func (o Object) GetSession() nxtypes.SessionHandle {
	if o.ObjectID == -1 {
		return nxtypes.SessionHandle(o.Content & 0xFFFFFFFF)
	}

	return nxtypes.SessionHandle(0)
}

func (o Object) SetDomain(domain *Domain) {
	o.Content = uint64(uintptr(unsafe.Pointer(domain)))
}

func (o Object) GetDomain() *Domain {
	if o.ObjectID >= 0 {
		return (*Domain)(unsafe.Pointer(uintptr(o.Content)))
	}

	return nil
}
