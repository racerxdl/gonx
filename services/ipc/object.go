package ipc

import (
	"github.com/racerxdl/gonx/nx/nxtypes"
	"runtime"
	"unsafe"
)

// Used for keeping GC from collecting it
var domainHolders []*Domain

func cacheDomain(domain *Domain) {
	for _, v := range domainHolders {
		if domain == v {
			return // Already cached
		}
	}

	domainHolders = append(domainHolders, domain)
}

// IPCObject Represents either an object within an IPC domain or a standalone object
type Object struct {
	ObjectID   int32 // -1 if this represents a session, >= 0 if this represents a domain object
	Content    uint64
	IsBorrowed bool
}

func (o *Object) Recycle() {
	if o.Content == 0 {
		return
	}

	for i, v := range domainHolders {
		ptr := uint64(uintptr(unsafe.Pointer(v)))
		if ptr == o.Content {
			o.Content = 0
			domainHolders = append(domainHolders[1:], domainHolders[i+1:]...)
			return
		}
	}

	o.Content = 0
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
	cacheDomain(domain)
	o.Content = uint64(uintptr(unsafe.Pointer(domain)))
}

func (o Object) GetDomain() *Domain {
	if o.ObjectID >= 0 && o.Content != 0 && uintptr(o.Content) > runtime.GetHeapBase() {
		return (*Domain)(unsafe.Pointer(uintptr(o.Content)))
	}

	return nil
}
